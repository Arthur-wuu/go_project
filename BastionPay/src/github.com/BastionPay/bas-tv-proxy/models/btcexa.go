package models

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/api"
	"BastionPay/bas-tv-proxy/base"
	"BastionPay/bas-tv-proxy/common"
	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/type"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var GBtcExaModels BtcExaModels

type BtcExaModels struct {
	mWsCon     *base.WsCon
	mWsHandler *BtcExaModelsWsHandler
	mUrl       *url.URL
	mHeader    http.Header
	mPush      common.PushHandler
	sync.Mutex
}

func (this *BtcExaModels) Init(pushHander common.PushHandler) error {
	source := &config.GConfig.BtcExa
	url, err := url.Parse(source.WsUrl)
	if err != nil {
		return err
	}
	header := make(http.Header)
	header.Set("ApiKey", source.ApiKey)
	this.mUrl = url
	this.mHeader = header
	this.mPush = pushHander
	return nil
}

func (this *BtcExaModels) Start() error {
	return this.WsConnection()
}

func (this *BtcExaModels) Stop() {

}

func (this *BtcExaModels) WsConnection() error {
	hd := NewBtcExaModelsWsHandler()
	c := new(base.WsCon)
	err := c.Init(uuid.New(), this.mUrl, this.mHeader, hd)
	if err != nil {
		return err
	}
	hd.mWsCon = c
	c.Start()
	this.Lock()
	this.mWsHandler = hd
	this.mWsCon = c
	this.Unlock()
	//	hd.Test()
	fmt.Println("WsConnection succ")
	return nil
}

func (this *BtcExaModels) WsSend(conid, ReqUuid, sub string, btcexaReq *_type.ReqBtcExaSub) error {
	this.Lock()
	newConFlag := (this.mWsCon == nil)
	this.Unlock()
	if newConFlag {
		ZapLog().Info("BtcExaModels need WsConnection")
		if err := this.WsConnection(); err != nil {
			return err
		}
	}
	json1 := btcexaReq.GetReqStr()
	content := []byte(json1)
	fmt.Println(string(content))

	//content, err := json.Marshal(btcexaReq)
	//if err != nil {
	//	ZapLog().Error("marshal err", zap.Error(err))
	//	return err
	//}
	needSendFlag := true
	if sub == api.EVENT_sub {
		needSendFlag = this.mWsHandler.Sub(btcexaReq.GetTopic(), ReqUuid, content)
	} else if sub == api.EVENT_unsub {
		needSendFlag = this.mWsHandler.UnSub(btcexaReq.GetTopic(), ReqUuid)
	}
	if !needSendFlag {
		ZapLog().Info("no need send sub or unsub")
		return nil
	}
	ZapLog().Info(string(content))
	err := this.mWsCon.Send(websocket.BinaryMessage, content)
	if err != nil {
		if sub == api.EVENT_sub {
			this.mWsHandler.UnSub(btcexaReq.Topic, ReqUuid)
		}
	}
	return err
}

//http
func (this *BtcExaModels) HttpbtcKXian(obj, period, limit string) (*_type.ResBtcExaKLine, error) {
	value := make(url.Values)
	if len(obj) != 0 {
		value.Set("trading_pair", obj)
	}
	if len(period) != 0 {
		value.Set("period", period)
	}
	if len(limit) != 0 {
		value.Set("limit", limit)
	}

	param := value.Encode()

	//sign := this.genSign(param)
	//content, err := base.HttpSend(config.GConfig.BtcExa.HttpUrl+"/api/market/kline?"+param, nil, "GET", map[string]string{ "Content-Type":"application/x-www-form-urlencoded","ApiKey":config.GConfig.BtcExa.ApiKey, "Signature":sign })
	content, err := base.HttpSend(config.GConfig.BtcExa.HttpUrl+"/api/market/kline?"+param, nil, "GET", nil)

	if err != nil {
		return nil, err
	}
	fmt.Println("kxian shuju : ", string(content))
	res := new(_type.ResBtcExaKLine)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil, err
	}
	if res.Status.Code != 0 {
		return nil, fmt.Errorf("%d %s", res.Status.Code, res.Status.Msg)
	}
	return res, nil
}

//btcexa交易对
func (this *BtcExaModels) HttpObjList() (*_type.ResBtcExaObjList, error) {

	content, err := base.HttpSend(config.GConfig.BtcExa.HttpUrl+"/api/market/trading_pairs", nil, "GET", nil)

	if err != nil {
		return nil, err
	}

	res := new(_type.ResBtcExaObjList)

	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil, err
	}
	if res.Status.Code != 0 {
		return nil, fmt.Errorf("%d %s", res.Status.Code, res.Status.Msg)
	}

	return res, nil
}

//params应该是按照名称排序并且url编码的   生成btcexa的sign
func (this *BtcExaModels) genSign(params string) string {
	str := "GET\n" + config.GConfig.BtcExa.HttpUrl + "/kline\n" + params + "\n" + config.GConfig.BtcExa.Secret_key
	h := sha256.New()
	h.Write([]byte(str))
	fmt.Println(str)
	bs := h.Sum(nil)
	fmt.Println(len(bs))

	return string(bs)
}

//type SubKeyInfo struct{
//	Data   []byte     //发送消息，断线重发的时候使用
//	SubMap map[string] *SubItem  //uuid, version
//}
//
//func (this *SubKeyInfo) Add(sessionId string, ver uint64){
//	this.SubMap[sessionId] = &SubItem{ver}
//}
//
//func (this *SubKeyInfo) Del(sessionId string){
//	delete(this.SubMap, sessionId)
//}
//
//func (this *SubKeyInfo) Len() int {
//	return len(this.SubMap)
//}

//func NewSubKeyInfo() *SubKeyInfo{
//	s := new(SubKeyInfo)
//	s.SubMap = make(map[string] *SubItem)
//	return s
//}

func NewBtcExaModelsWsHandler() *BtcExaModelsWsHandler {
	c := new(BtcExaModelsWsHandler)
	c.mRecvResChan = make(chan error, 5)
	c.mSubKeyMap = make(map[string]*SubKeyInfo)
	return c
}

//handler
type BtcExaModelsWsHandler struct {
	mRecvResChan chan error
	mWsCon       *base.WsCon
	mSubKeyMap   map[string]*SubKeyInfo //topic
	ConId        string
	mSendLock    sync.Mutex
	sync.Mutex
}

func (this *BtcExaModelsWsHandler) IsRecvPingHander() bool {
	return true
}

func (this *BtcExaModelsWsHandler) IsRecvPongHander() bool {
	return true
}

func (this *BtcExaModelsWsHandler) RecvPingHander(message string) string {
	return ""
}

func (this *BtcExaModelsWsHandler) RecvPongHander(message string) string {
	return ""
}

func (this *BtcExaModelsWsHandler) SendPingHander() []byte {
	return []byte("{\"event\":\"ping\"}")
}

func (this *BtcExaModelsWsHandler) BeforeSendHandler() {
	this.mSendLock.Lock()
	ZapLog().Info("BeforeSendHandler lock")
	for {
		select {
		case <-this.mRecvResChan:
			break
		default:
			return
		}
	}
}

func (this *BtcExaModelsWsHandler) AfterSendHandler() error {
	defer this.mSendLock.Unlock()
	ZapLog().Info("AfterSendHandler Unlock")
	select {
	case err := <-this.mRecvResChan:
		return err
	case <-time.After(time.Second * 30):
		return errors.New("chan timeout 30s")
	}
	return nil
}

func (this *BtcExaModelsWsHandler) RecvHandler(conId string, message []byte) {
	slic := make([]interface{}, 0)
	err := json.Unmarshal(message, &slic)
	if err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("btcexa_sub_msg", string(message)))
		return
	}
	if len(slic) <= 0 {
		ZapLog().Error("unMarshal err , no data")
		return
	}

	topic, ok := slic[0].(string)
	ty, ok := slic[1].(string)
	status, ok := slic[2].(string)

	newTopic := strings.TrimLeft(topic, "sub.")
	apiMsg := new(api.MSG)

	if ty == "r" {
		if strings.ToUpper(status) == "OK" {
			this.mRecvResChan <- nil
		} else {
			this.mRecvResChan <- fmt.Errorf("%s", ty)
		}
		return
	}
	if ty == "i" {
		//return init的数据放开
		res := new(_type.ResBtcExaSubI)
		res.Topic = topic
		res.Type = ty
		res.Status = status
		data1, ok := slic[3].([]interface{})
		if !ok {
			ZapLog().Error("change i data wrong1")
			return
		}
		dataValue1 := make([][]string, 0)
		for i := 0; i < len(data1); i++ {
			data2, ok := data1[i].([]interface{})
			if !ok {
				ZapLog().Error("change i data wrong2 ", zap.Any("type", reflect.TypeOf(data1[i])))
				return
			}

			dataValue2 := make([]string, 0)
			for j := 0; j < len(data2); j++ {
				data3 := data2[j].(string)
				dataValue2 = append(dataValue2, data3)
			}
			dataValue1 = append(dataValue1, dataValue2)
		}
		res.Data = dataValue1

		if strings.Contains(newTopic, "kline") {
			apiKxian := BtcExaSubToApiKXianI(res)
			ZapLog().Info("recv=", zap.Any("res=", apiKxian))
			topicArr := strings.Split(newTopic, ".")
			obj := ""
			if len(topicArr) >= 3 {
				obj = topicArr[1]
			}
			apiMsg.AddQuoteKlineSingle(api.NewQuoteKlineSingle(obj, apiKxian))
		} else {
			ZapLog().Error("GBtcExaModels unknown sub msg")
			return
		}
	}

	if ty == "u" {
		res := new(_type.ResBtcExaSubU)
		res.Topic = topic
		res.Type = ty
		res.Status = status
		data, ok := slic[3].([]interface{})
		if !ok {
			ZapLog().Error("recv data wrong1")
			return
		}
		dataValue := make([]string, 0)
		if data != nil {
			for i := 0; i < len(data); i++ {
				data, ok := data[i].(string)
				if !ok {
					ZapLog().Error("recv data wrong2")
					return
				}
				dataValue = append(dataValue, data)
			}
			res.Data = dataValue
		} else {
			return
		}
		if strings.Contains(newTopic, "kline") {
			apiKxian := BtcExaSubToApiKXianU(res)
			ZapLog().Info("recv=", zap.Any("res=", apiKxian))
			topicArr := strings.Split(newTopic, ".")
			obj := ""
			if len(topicArr) >= 3 {
				obj = topicArr[1]
			}
			apiMsg.AddQuoteKlineSingle(api.NewQuoteKlineSingle(obj, apiKxian))
		} else {
			ZapLog().Error("GBtcExaModels unknown sub msg")
			return
		}
	}
	// else 各种消息类型
	//推送消息
	this.Lock()
	defer this.Unlock()
	subKeyInfoList, ok := this.mSubKeyMap[newTopic]
	if !ok {
		ZapLog().Error("unknown topic", zap.String("topic", newTopic))
		return
	}
	for k, v := range subKeyInfoList.SubMap {
		ZapLog().Info("push uuid =" + k)
		base.GWsServerMgr.Push(k, v.Ver, apiMsg)
	}
}

func (this *BtcExaModelsWsHandler) ReSendAllRequest() {
	errMap := make(map[string]bool)
	subMap := this.GetSubMap()
	ZapLog().Info("start ReSendAllRequest", zap.Int("topicnum", len(subMap)))
	for topic, info := range subMap {
		ZapLog().Info("start topic " + topic)
		if err := this.mWsCon.Send(websocket.BinaryMessage, info); err != nil {
			ZapLog().Sugar().Errorf("ReSendAllRequest fail %v", topic)
			errMap[topic] = true
		} else {
			ZapLog().Info("ReSendAllRequest success", zap.String("topic", topic))
		}
		ZapLog().Info("end topic " + topic)
	} //重试3次
}

func (this *BtcExaModelsWsHandler) GetSubMap() map[string][]byte {
	mm := make(map[string][]byte)
	this.Lock()
	defer this.Unlock()
	for topic, info := range this.mSubKeyMap {
		mm[topic] = info.Data
	}
	return mm
}

func (this *BtcExaModelsWsHandler) Sub(topic, sessionId string, data []byte) bool {
	this.Lock()
	defer this.Unlock()
	subKeyInfo, ok := this.mSubKeyMap[topic]
	if !ok {
		subKeyInfo = NewSubKeyInfo()
		this.mSubKeyMap[topic] = subKeyInfo
	}
	subKeyInfo.Add(sessionId, 0)
	subKeyInfo.Data = data
	return !ok
}

func (this *BtcExaModelsWsHandler) UnSub(topic, sessionId string) bool {
	this.Lock()
	defer this.Unlock()
	subKeyInfo, ok := this.mSubKeyMap[topic]
	if !ok {
		return false
	}
	subKeyInfo.Del(sessionId)
	if subKeyInfo.Len() == 0 {
		delete(this.mSubKeyMap, topic)
		return true
	} else {
		return false
	}

}

func BtcExaSubToApiKXianU(k1 *_type.ResBtcExaSubU) []*api.KXian {
	k2 := make([]*api.KXian, 0)
	tmp := new(api.KXian)

	if len(k1.Data) <= 0 {
		ZapLog().Error("data length <= 0 btcexa U data nil")
		return nil
	}
	string := k1.Data[0]
	in, err := strconv.ParseInt(string, 10, 64)
	if err != nil {
		ZapLog().Error("string to int err", zap.Error(err))
	}

	tmp.SetShiJian(in)
	tmp.SetKaiPanJia(k1.Data[1])
	tmp.SetZuiGaoJia(k1.Data[2])
	tmp.SetZuiDiJia(k1.Data[3])
	tmp.SetShouPanJia(k1.Data[4])
	tmp.SetChengJiaoLiang(k1.Data[5])

	k2 = append(k2, tmp)

	return k2
}

func BtcExaSubToApiKXianI(k1 *_type.ResBtcExaSubI) []*api.KXian {
	k2 := make([]*api.KXian, 0)
	for j := 0; j < len(k1.Data); j++ {
		tmp := new(api.KXian)

		if len(k1.Data) <= 0 {
			ZapLog().Error("data length <= 0 btcexa I data nil")
			return nil
		}

		for i := 0; i < len(k1.Data[j]); i++ {
			switch i {
			case 0:
				in, err := strconv.ParseInt(k1.Data[j][0], 10, 64)
				if err != nil {
					ZapLog().Error("string to int err", zap.Error(err))
				}
				tmp.SetShiJian(int64(in))
				break
			case 1:
				tmp.SetKaiPanJia(k1.Data[j][1])
				break
			case 2:
				tmp.SetZuiGaoJia(k1.Data[j][2])
				break
			case 3:
				tmp.SetZuiDiJia(k1.Data[j][3])
				break
			case 4:
				tmp.SetShouPanJia(k1.Data[j][4])
				break
			case 5:
				tmp.SetChengJiaoLiang(k1.Data[j][5])
				break
			}
		}
		k2 = append(k2, tmp)
	}

	return k2
}
