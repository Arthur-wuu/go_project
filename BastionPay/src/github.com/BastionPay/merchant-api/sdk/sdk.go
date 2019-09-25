package sdk

import (
	"BastionPay/merchant-api/baspay"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"BastionPay/merchant-api/util"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/base"
	"github.com/satori/go.uuid"
)

const testCommonUrl    =  "https://test-openapi.bastionpay.io/"
const commonUrl        =  "https://open-api.bastionpay.com/"

const transUrl         =  commonUrl+"open-api/trade/transfer"
const availAssetsUrl   =  commonUrl+"open-api/trade/avail_assets"
const createQrUrl      =  commonUrl+"open-api/trade/create_qr_trade"
const createSdkUrl     =  commonUrl+"open-api/trade/create_trade"
const createWapUrl     =  commonUrl+"open-api/trade/create_wap_trade"
const tradeInfoUrl     =  commonUrl+"open-api/trade/info"
const posUrl           =  commonUrl+"open-api/trade/pos"
const posOrderUrl      =  commonUrl+"open-api/trade/pos_orders"



var h5PrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEArQ+bLJlYR4dn0l8yRGdN2ZSFns9NYLdL/DBLB70twKeRu5nL
7cAm0SzIMj6YtXgY9itPwOup53xFTCtvN322YY9Cffb1oBkcIY2PmtyX1C7RbDpD
P7tzwZWdz+eP67be19EMqfYXyYoPmMnoL2sPlMY34L/Jgop+2xw8F6kMRSUyZ8WU
+N0OswYF2W9S++fLWe1wZ+37WaUtAiOe721CtIVyQF2B7SDdGHoJ5AxRjcSP+rd5
qck99iGF1aAT+MYDFsBDBg5jZl8z6y1XccOQ50EVaNW7Gt6E+7951rZFA5RgAZqO
0lJtKFeKC3TSIuWsDhioGJ6m0aLqK7SDGz3FhwIDAQABAoIBAQCOofQkl+X4XhMl
gawuUG4LS6utLfH6KlgH682K/VI+HF2yHpnCw8G6WIxPTOQTfH4mNaAvwotv7C45
DvtE4ul0EtyccPQUFV3oEYIwAmtoR4X3CzXtyxMmk6dTeOhXP4r+mJ81XUxRoOYl
6RLiMfzPg2b+Z1MvsfHHqMemQH+KZDjdxXLGazk1DNAQXt8V/5e1SJhVcpv48QRC
ZeZHVRQ3f/6FL66edpGpQ5/oLBxw5m4tHiTSRg4T+a7sO4cCoZapCcsxaKx+s9eJ
DYjIOxPTNweqVVRDjrJ+fXxXKefbYx2AddUCecUIacyPMhDbWs5x//+FJkCFKvyH
eefKxWxxAoGBAM9Vxi0HaMCHdUkzUh/fX0uV8+f1kdl2k4W0CycXQA5zt1h2f2ve
fVcAEff5YwKYSTWlVgySxqyOVtz/ellzcEKcHiDZadIeCDYgoa23v266mnVpPYqT
n8jFwHhyVrPD/S+iLz8PrpF7NPYkQUNn40RIxNrEjMvUjyzatqgv8sTrAoGBANWu
Y/6t40mcBkCoWy0UJBuujPskiFrR9GLXQZ0GWhlSC+KKCOnO8/suWCicQqD1CD51
Nalpk+E3+wbXI/qDUpRq5WBob+l4CWKugtHhgaUdw8c+naWK6xn9okmD0Gadd1qr
jtAiHludI9fkNmkVIzaUkK26uBwOK41g9jMW/ErVAoGAG8BbWkOXnc2DwVyBLYr0
cmWL1AxmjTj13fuPUpgmFskeTVTvET0igbacsRhMTFid0/RhZCVxOj+DGOmJMtfk
usWysqrnIxyp9LTBb0Mc+HE5o2WGuzmvNWxiqryDJmShSvLmaAZtU0OufxOzOJZ7
MPSchLuyLMYys8pCkJh6YikCgYB/YUJC5C4GB3jCupn/uW39AoUQgaq3WUmyUlfO
36Z+SabEGT1PBAv1xJ7RNrWRdgDAGucuYr3BGLoQTdgo0ng7+a1bV2a/astNhHJ/
40qBv8ih0fXwZWvZRpWj9Wwaf+xSpMqx0GUAgCCJ5oV5BxzCwLWumwx9zQSxdwfN
VPp5MQKBgQDCwEBWR4xzjbt5MoE8pggF0an8zRGv3Z6yP0vFP+q/I6ykCrfNxByP
uchEixbp9Oe43zDzw8r/dIcCPPRw44A7voWfrAQKhyUjMVVki+ovGSQBZ5I6n9/7
AV5pE9Mk3YF7sLozbtKGyZev1YxnV3poC1DS+hTbcJdLFPCj3TQC1A==
-----END RSA PRIVATE KEY-----
`)


type Config struct {
		KeyPath string   //`yaml:"keyPath"`
	}

var (
	//私钥
	prikey []byte

	//公钥
	pubkey []byte
)
type (
	//转账参数
	TransferParam struct{
		Amount           *string   `valid:"required" json:"amount,omitempty"`
		Assets           *string   `valid:"required" json:"assets,omitempty"`
		MerchantId       *string   `valid:"required" json:"merchant_id,omitempty"`
		MerchantTransNo   string   `valid:"optional" json:"-"`
		NotifyUrl        *string   `valid:"optional" json:"notify_url,omitempty"`
		PayeeId          *string   `valid:"required" json:"payee_id,omitempty"`
		ProductName      *string   `valid:"required" json:"product_name,omitempty"`
		Timestamp        *string   `valid:"optional" json:"timestamp,omitempty"`
		SignType         *string   `valid:"optional" json:"sign_type,omitempty"`
		Signature        *string   `valid:"optional" json:"signature,omitempty"`
	}

	TransRes struct{
		Code              int64    `json:"code"`
		Message           string   `json:"message"`
		Data              DataRes  `json:"data"`
	}

	DataRes struct {
		Assets             string  `json:"assets,omitempty"`
		Amount             string  `json:"amount,omitempty"`
		MerchantTransferNo string  `json:"merchant_transfer_no,omitempty"`
		Status			   int     `json:"status,omitempty"`
		TransferNo		   string  `json:"transfer_no,omitempty"`
	}
)

var GPaySdk PaySdk

type PaySdk struct{

}

//pem文件的路径
func (this* PaySdk) Init(config *Config) error {
	//加载pem文件，检测有效性
	if err := loadRsaKeys(config); err != nil {
		fmt.Println("err",err)
		fmt.Errorf("BastionPay Init: %s", err.Error())
		os.Exit(1)
	}
	return nil
}

//转账接口       数量，币种，商户id，收款人id
func (this* PaySdk) Transfer(amount, assets, merchant_id, payee_id, produceName, notifyUrl string)( string, error) {
	param := new(TransferParam)
	times := time.Now().Local().Format("2006-01-02 15:04:05")

	param.Timestamp = &times
	param.MerchantTransNo = GenerateUuid()
	param.Amount = &amount
	param.Assets = &assets
	param.PayeeId = &payee_id
	param.ProductName = &produceName
	param.MerchantId = &merchant_id
	param.NotifyUrl = &notifyUrl


	transferStatus, err := param.Send()
	if err != nil {
		return "fail", err
	}

	if transferStatus == "1" {
		return "unpay", err
	}

	if transferStatus == "3" {
		return "succ", err
	}

	if transferStatus == "4" {
		return "fail", err
	}

	return "fail", nil
}




// 加载数据
func loadRsaKeys(config *Config) error {
	var err error

	private := fmt.Sprintf("%s/%s", strings.Trim(config.KeyPath, "/"), "private.pem")
	prikey, err = ioutil.ReadFile(private)
	if err != nil {
		return err
	}

	public := fmt.Sprintf("%s/%s", strings.Trim(config.KeyPath, "/"), "public.pem")
	pubkey, err = ioutil.ReadFile(public)
	if err != nil {
		return err
	}
	return nil
}



//转账请求
func (this *TransferParam) Send() (string, error) {
	signType := "RSA"
	this.SignType = &signType

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"amount": this.Amount,
		"assets": this.Assets,
		"payee_id": this.PayeeId,
		"product_name": this.ProductName,
		"merchant_transfer_no": this.MerchantTransNo,
		"merchant_id": this.MerchantId,
		"timestamp": this.Timestamp,
		"notify_url": this.NotifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(baspay.CommonPrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error( "sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"amount": this.Amount,
		"assets": this.Assets,
		"payee_id": this.PayeeId,
		"merchant_transfer_no": this.MerchantTransNo,
		"merchant_id": this.MerchantId,
		"product_name": this.ProductName,
		"timestamp": this.Timestamp,
		"notify_url": this.NotifyUrl,
		"sign_type": this.SignType,
		"signature": finalSign,
	})
	fmt.Println("**reqBody**",string(reqBody))

	result, err := base.HttpSend(transUrl, bytes.NewBuffer(reqBody),"POST", nil)
	if err != nil {
		fmt.Println("err",err)
		//ZapLog().Error( "request err")
		return "", err
	}

	fmt.Println("**result**",string(result))

	transRes := new(TransRes)
	err = json.Unmarshal(result, transRes)
	fmt.Println("err:", err)
	fmt.Println("status:",transRes.Data.Status)
	status := strconv.Itoa(transRes.Data.Status)

	return status, nil
}


func RequestBodyToSignStr (body []byte) (string){
	requestParams := make(map[string]string,0)

	err := json.Unmarshal(body, &requestParams)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("requestbody to requestParams err")
		return ""
	}
	//将param的key排序，
	keysSort := make([]string, 0)
	for k, _ := range requestParams{
		keysSort = append(keysSort, k)
	}
	sort.Strings(keysSort)
	//拼接签名字符串
	signH5Str := ""
	for i:=0; i<len(keysSort); i++ {
		signH5Str += keysSort[i]+"="+requestParams[keysSort[i]]+"&"
	}
	signH5Str = signH5Str[0:len(signH5Str)-1]
	return  signH5Str
}




func  GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}

//
//
////有效币种接口 avail_assets
//func (this* PaySdk) AvailAssets( arg ... interface{})( interface{}, error) {
//	param := new(AvailAssetsParam)
//	if len(arg) ==  1 {
//		argParam := arg[0].(string)
//		param.Assets = &argParam
//		availAssetsRes, err := param.Send()
//		if err != nil {
//			return nil, err
//		}
//		return availAssetsRes, nil
//	}
//
//	if len(arg) ==  0 {
//		availAssetsRes, err := param.Send()
//		if err != nil {
//			return nil, err
//		}
//		return availAssetsRes, nil
//	}
//
//	return nil, nil
//}

//
//
////二维码下单接口       create_qr_trade
//func (this* PaySdk) CreateQrtrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
//	param := new(QrTrade)
//	times := time.Now().Local().Format("2006-01-02 15:04:05")
//
//	param.Timestamp = &times
//	param.MerchantTradeNo = GenerateUuid()
//	param.Amount = &amount
//	param.Assets = &assets
//	param.PayeeId = &payee_id
//	param.MerchantId = &merchant_id
//	param.ProductName = &product_name
//	param.ProductDetail = &product_detail
//	param.Remark = &remark
//	param.ReturnUrl = &return_url
//	param.ShowUrl = &show_url
//	param.NotifyUrl = &notifyUrl
//	param.ExpireTime = expire_time
//
//
//	createQrRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return createQrRes, nil
//}
//
//
//
//
//
////Sdk下单接口       create sdk trade
//func (this* PaySdk) CreateSdkTrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
//	param := new(SdkTrade)
//	times := time.Now().Local().Format("2006-01-02 15:04:05")
//
//	param.Timestamp = &times
//	param.MerchantTradeNo = GenerateUuid()
//	param.Amount = &amount
//	param.Assets = &assets
//	param.PayeeId = &payee_id
//	param.MerchantId = &merchant_id
//	param.ProductName = &product_name
//	param.ProductDetail = &product_detail
//	param.Remark = &remark
//	param.ReturnUrl = &return_url
//	param.ShowUrl = &show_url
//	param.NotifyUrl = &notifyUrl
//	param.ExpireTime = expire_time
//
//
//	createSdkRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return createSdkRes, nil
//}
//
//
//
//
//
//
////wap下单接口       create wap trade
//func (this* PaySdk) CreateWapTrade(amount, assets, merchant_id, payee_id, product_detail,product_name,remark,return_url,show_url, notifyUrl string , expire_time int64)( interface{}, error) {
//	param := new(WapTrade)
//	times := time.Now().Local().Format("2006-01-02 15:04:05")
//
//	param.Timestamp = &times
//	param.MerchantTradeNo = GenerateUuid()
//	param.Amount = &amount
//	param.Assets = &assets
//	param.PayeeId = &payee_id
//	param.MerchantId = &merchant_id
//	param.ProductName = &product_name
//	param.ProductDetail = &product_detail
//	param.Remark = &remark
//	param.ReturnUrl = &return_url
//	param.ShowUrl = &show_url
//	param.NotifyUrl = &notifyUrl
//	param.ExpireTime = expire_time
//
//
//	createWapRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return createWapRes, nil
//}
//
//
//
//
////查询交易信息       5aa32a38-4214-4e52-a384-5ad8e3b54fb0
//func (this* PaySdk) SearchTradeInfo ( merchant_id, merchant_trade_no, notifyUrl string )( interface{}, error) {
//	param := new(TradeInfo)
//	param.MerchantId = &merchant_id
//	param.MerchantTradeNo = &merchant_trade_no
//	param.NotifyUrl = &notifyUrl
//
//	infoRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return infoRes, nil
//}
//
//
//
////pos支付接口       create_qr_trade
//func (this* PaySdk) CreatePos (amount, assets, merchant_id, payee_id,pay_voucher,pos_merchine_id, product_detail,product_name,remark, notifyUrl string )( interface{}, error) {
//	param := new(PosTrade)
//	times := time.Now().Local().Format("2006-01-02 15:04:05")
//
//	param.Timestamp = &times
//	param.MerchantPosNo = GenerateUuid()
//	param.PosMachineId = &pos_merchine_id
//	param.PayVoucher = &pay_voucher
//	param.Amount = &amount
//	param.Assets = &assets
//	param.PayeeId = &payee_id
//	param.MerchantId = &merchant_id
//	param.ProductName = &product_name
//	param.ProductDetail = &product_detail
//	param.Remark = &remark
//	param.NotifyUrl = &notifyUrl
//
//
//	createQrRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return createQrRes, nil
//}
//
//
//
//
////pos支付订单列表接口
//func (this* PaySdk) PosOrderList (begin_time,end_time, merchant_id,merchant_pos_no, pos_merchine_id, page,page_size,trade_status, notifyUrl string )( interface{}, error) {
//	param := new(PosOrderList)
//	times := time.Now().Local().Format("2006-01-02 15:04:05")
//
//	param.TimeStamp = &times
//	param.MerchantPosNo = &merchant_pos_no
//	param.PosMachineId = &pos_merchine_id
//	param.MerchantId = &merchant_id
//	param.NotifyUrl = &notifyUrl
//	param.BeginTime = &begin_time
//	param.EndTime = &end_time
//	param.Page = &page
//	param.PageSize = &page_size
//	param.TradeStatus = &trade_status
//
//
//	createQrRes, err := param.Send()
//	if err != nil {
//		return "fail", err
//	}
//
//	return createQrRes, nil
//}
//
//



















