package bastionpay

import (
	"crypto"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common/bastionpay/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var (
	gateway              = "http://35.173.156.149:8082/api"
	urlGetMessage        = "/v1/bastionpay/history_transaction_message"
	urlGetSupports       = "/v1/bastionpay/support_assets"
	urlCurrencrysInfo    = "/v1/bastionpay/asset_attribute"
	urlGetHeight         = "/v1/bastionpay/last_block_height"
	urlGetAddress        = "/v1/bastionpay/new_address"
	urlGetBalance        = "/v1/bastionpay/get_balance"
	urlWithdrawal        = "/v1/bastionpay/withdrawal"
	urlOrderHistory      = "/v1/bastionpay/history_transaction_order"
	msgChannelBufferSize = 50
)

type Session struct {
	*Config
	privateKeys          []byte
	publicKeys           []byte
	bastionPayPublicKeys []byte

	MsgChannel     chan *CallbackMsg
	lastMsgID      int64
	lastMsgIDMutex sync.Mutex
}

func New(config *Config) (*Session, error) {
	s := &Session{
		Config:     config,
		MsgChannel: make(chan *CallbackMsg, msgChannelBufferSize),
	}
	err := s.loadCerts()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Session) loadCerts() error {
	var err error
	s.privateKeys, err = ioutil.ReadFile(s.Config.PrivateKey)
	if err != nil {
		return err
	}

	s.publicKeys, err = ioutil.ReadFile(s.Config.PublicKey)
	if err != nil {
		return err
	}

	s.bastionPayPublicKeys, err = ioutil.ReadFile(s.Config.BastionPayPublicKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Send(url string, data interface{}, result interface{}) error {
	var err error

	source, err := s.encryption(data)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}

	res, err := utils.Post(url, source)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}

	//var resSource Source
	var response Response
	json.Unmarshal(res, &response)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}
	if response.Err != 0 {
		return errors.New(response.ErrMsg)
	}

	des, err := s.decryption(response.Value)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}

	err = json.Unmarshal(des, result)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}

	return nil
}

func (s *Session) encryption(data interface{}) (Source, error) {
	var (
		source     Source
		message    []byte
		encMessage []byte
	)

	source.ApiKey = s.ApiKey

	message, err := json.Marshal(data)
	if err != nil {
		return source, err
	}

	// 加密
	encMessage, err = utils.RsaEncrypt(message, s.bastionPayPublicKeys, utils.RsaEncodeLimit2048)
	if err != nil {
		return source, err
	}
	source.Message = base64.StdEncoding.EncodeToString(encMessage)

	// 签名
	hs := sha512.New()
	hs.Write(encMessage)
	hashData := hs.Sum(nil)

	signature, err := utils.RsaSign(crypto.SHA512, hashData, s.privateKeys)
	if err != nil {
		return source, err
	}
	source.Signature = base64.StdEncoding.EncodeToString(signature)

	return source, nil
}

func (s *Session) decryption(source Source) ([]byte, error) {
	var (
		cipherText []byte
		signature  []byte
		err        error
	)

	cipherText, err = base64.StdEncoding.DecodeString(source.Message)
	if err != nil {
		return nil, err
	}

	signature, err = base64.StdEncoding.DecodeString(source.Signature)
	if err != nil {
		return nil, err
	}

	// 验签
	hs := sha512.New()
	hs.Write([]byte(cipherText))
	hs.Sum(nil)

	err = utils.RsaVerify(crypto.SHA512, hs.Sum(nil), signature, s.bastionPayPublicKeys)
	if err != nil {
		return nil, err
	}

	// 解密
	return utils.RsaDecrypt(cipherText, s.privateKeys, utils.RsaDecodeLimit2048)
}

func (s *Session) SetLastMsgID(lastMsgID int64) {
	s.lastMsgID = lastMsgID
}

func (s *Session) handlePush(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		w.Write([]byte("404 page not found"))
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		l4g.Error(err.Error())
		return
	}

	res := Response{}
	err = json.Unmarshal(b, &res)
	if err != nil {
		l4g.Error(err.Error())
		return
	}

	decode, err := s.decryption(res.Value)
	if err != nil {
		l4g.Error(err.Error())
		return
	}

	callbackMsg := CallbackMsg{}
	err = json.Unmarshal(decode, &callbackMsg)
	if err != nil {
		l4g.Error(err.Error())
		return
	}

	if callbackMsg.ID > (s.lastMsgID + 1) {
		s.GetMsgs(s.lastMsgID+1, callbackMsg.ID)
	} else {
		if callbackMsg.ID == (s.lastMsgID + 1) {
			s.MsgChannel <- &callbackMsg
			s.lastMsgIDMutex.Lock()
			s.lastMsgID = callbackMsg.ID
			s.lastMsgIDMutex.Unlock()
		}
	}

	w.Write([]byte("Complate"))
}

// 启用充提回调监听
func (s *Session) EnableMonitoring(port string) error {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Service health"))
	}))
	http.Handle("/callback", http.HandlerFunc(s.handlePush))

	fmt.Printf("Now listening on :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		l4g.Error(err.Error())
		return err
	}
	return nil
}

// 获取回调数据
func (s *Session) GetMsgs(from int64, to int64) ([]*CallbackMsg, error) {
	res := []*CallbackMsg{}
	err := s.Send(url(urlGetMessage), RequestHistoryMsg{from, to}, &res)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		if v.ID == (s.lastMsgID + 1) {
			s.MsgChannel <- v
			s.lastMsgIDMutex.Lock()
			s.lastMsgID = v.ID
			s.lastMsgIDMutex.Unlock()
		}
	}

	return res, nil
}

// 获取支持币种
func (s *Session) GetSupports() ([]string, error) {
	symbols := []string{}
	err := s.Send(url(urlGetSupports), nil, &symbols)
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

// 获取支持币种信息
func (s *Session) GetCurrencysInfo(symbols []string) ([]*ResponseCurrency, error) {
	for k, _ := range symbols {
		symbols[k] = strings.ToLower(symbols[k])
	}

	currencys := []*ResponseCurrency{}
	err := s.Send(url(urlCurrencrysInfo), symbols, &currencys)
	if err != nil {
		return nil, err
	}
	return currencys, nil
}

// 获取块高
func (s *Session) GetHeight(symbols []string) ([]*ResponseHeight, error) {
	for k, _ := range symbols {
		symbols[k] = strings.ToLower(symbols[k])
	}

	heights := []*ResponseHeight{}
	err := s.Send(url(urlGetHeight), symbols, &heights)
	if err != nil {
		return nil, err
	}
	return heights, nil
}

// 获取地址
func (s *Session) GetAddress(symbol string, count int) ([]string, error) {
	symbol = strings.ToLower(symbol)

	res := &ResponseAddress{}
	err := s.Send(url(urlGetAddress), RequestAddress{symbol, count}, &res)
	if err != nil {
		return nil, err
	}
	return res.Addresses, nil
}

// 获取账户余额
func (s *Session) GetBalance(symbols []string) ([]*ResponseBalance, error) {
	for k, _ := range symbols {
		symbols[k] = strings.ToLower(symbols[k])
	}
	res := []*ResponseBalance{}
	err := s.Send(url(urlGetBalance), symbols, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 提币
func (s *Session) Withdrawal(id string, symbol string, address string, amount float64) (string, error) {
	symbol = strings.ToLower(symbol)

	res := &ResponseWithdrawal{}
	err := s.Send(url(urlWithdrawal), RequestWithdrawal{
		UserOrderID: id,
		Symbol:      symbol,
		Address:     address,
		Amount:      amount,
	}, &res)
	if err != nil {
		return "", err
	}
	return res.OrderID, nil
}

// 获取订单历史
func (s *Session) OrderHistory(symbol string, transType int, status int, maxUpdateTime int64, minUpdateTime int64) ([]*ResponseOrder, error) {
	symbol = strings.ToLower(symbol)

	orders := []*ResponseOrder{}
	err := s.Send(url(urlOrderHistory), RequestOrder{
		Symbol:        symbol,
		TransType:     transType,
		Status:        status,
		MaxUpdateTime: maxUpdateTime,
		MinUpdateTime: minUpdateTime,
	}, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func url(address string) string {
	return gateway + address
}
