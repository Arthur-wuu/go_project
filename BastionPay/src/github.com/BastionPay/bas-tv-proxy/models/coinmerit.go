package models

import (
	"BastionPay/bas-tv-proxy/base"
	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/type"
	"BastionPay/bas-tv-proxy/api"
	"net/url"
	"time"
	"fmt"
	"crypto/md5"
	"strings"
	"errors"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"sync"
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"BastionPay/bas-tv-proxy/common"
	"strconv"
)


var GCoinMeritModels CoinMeritModels
type CoinMeritModels struct{
	mWsCon *base.WsCon
	mWsHandler  *CoinMeritModelsWsHandler
	mUrl *url.URL
	mHeader http.Header
	mPush       common.PushHandler
	sync.Mutex
}

func (this *CoinMeritModels) Init(pushHander common.PushHandler) error {
	source := &config.GConfig.CoinMerit
	url,err := url.Parse(source.WsUrl)
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

func (this *CoinMeritModels) Start() error {
	return this.WsConnection()
}

func (this *CoinMeritModels) Stop() {

}

func (this *CoinMeritModels) WsSend(conid, ReqUuid, sub string, cmReq *_type.ReqCoinMeritSub) error {
	this.Lock()
	newConFlag :=  (this.mWsCon == nil )
	this.Unlock()
	if newConFlag {
		ZapLog().Info("CoinMeritModels need WsConnection")
		if err := this.WsConnection(); err != nil {
			ZapLog().Error("CoinMeritModels WsConnect err", zap.String("conid", "coinmerit"), zap.Error(err))
			return err
		}
	}
	content, err := json.Marshal(cmReq)
	if err != nil {
		ZapLog().Error("marshal err", zap.Error(err))
		return err
	}
	ZapLog().Info("will send="+ string(content), zap.String("sub", sub))
	needSendFlag := true
	if sub == api.EVENT_sub {
		needSendFlag = this.mWsHandler.Sub(cmReq.Topic, ReqUuid, content)
		ZapLog().Info("sub topic "+ cmReq.Topic+"  "+ReqUuid)
	}else if sub == api.EVENT_unsub  {
		needSendFlag = this.mWsHandler.UnSub(cmReq.Topic, ReqUuid)
		ZapLog().Info("unsub topic "+ cmReq.Topic+"  "+ReqUuid)
	}
	if !needSendFlag {
		ZapLog().Info("noSeed sub or unsub", zap.String("conid", this.mWsCon.GetId()))
		return nil
	}
	ZapLog().Info("sned====>>>>>"+ string(content))
	err = this.mWsCon.Send(websocket.BinaryMessage, content);
	if err != nil {
		ZapLog().Error("CoinMeritModels Send err", zap.String("conid", this.mWsCon.GetId()), zap.Error(err))
		if sub == api.EVENT_sub {
			this.mWsHandler.UnSub(cmReq.Topic, ReqUuid)
		}
	}
	return err
}

func (this *CoinMeritModels) WsConnection() error {
	hd := NewCoinMeritModelsWsHandler(this.mPush)
	c:= new(base.WsCon)
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
	return nil
}

func (this *CoinMeritModels) HttpExa() (*_type.ResCoinMeritExchanges, error){
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := this.genSign("exchanges", timeStamp, "")
	content, err := base.HttpSend(config.GConfig.CoinMerit.HttpUrl+"/exchanges", nil, "GET", map[string]string{"ApiKey":config.GConfig.CoinMerit.ApiKey, "Timestamp":timeStamp, "ApiSign":sign })
	if err != nil {
		return nil,err
	}
	res := new(_type.ResCoinMeritExchanges)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err))
		return nil,err
	}
	if res.Status_code != 200 {
		return nil, fmt.Errorf("%d %s", res.Status_code, res.Message)
	}
	return res,nil
}

func (this *CoinMeritModels) HttpObjList(exa string) (*_type.ResCoinMeritCurrencyPairs, error){
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := this.genSign("currency_pairs", timeStamp, "exchange="+exa)
	content, err := base.HttpSend(config.GConfig.CoinMerit.HttpUrl+"/currency_pairs?exchange="+exa, nil, "GET", map[string]string{"ApiKey":config.GConfig.CoinMerit.ApiKey, "Timestamp":timeStamp, "ApiSign":sign })
	if err != nil {
		return nil,err
	}
	res := new(_type.ResCoinMeritCurrencyPairs)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("content", string(content)))
		return nil,err
	}
	if res.Status_code != 200 {
		return nil, fmt.Errorf("%d %s", res.Status_code, res.Message)
	}
	return res,nil
}

func (this *CoinMeritModels) HttpKXian(exchange, symbol, tp, size, since string) (*_type.ResCoinMeritKLine, error){
	value := make(url.Values)
	if len(exchange) != 0 {
		value.Set("exchange", exchange)
	}
	if len(symbol) != 0 {
		value.Set("symbol", symbol)
	}
	if len(tp) != 0 {
		value.Set("type", tp)
	}
	if len(size) != 0 {
		count,_ := strconv.Atoi(size)
		if count == 0 {
			size = "13"
		}else if count > 2000 {
			size = "2000"
		}
		value.Set("size", size)
	}else{
		value.Set("size", "13")
	}
	if len(since) != 0 {
		value.Set("since", since)
	}

	param := value.Encode()
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	sign := this.genSign("kline", timeStamp, param)
//	ZapLog().Info("models coinmerit "+config.GConfig.CoinMerit.HttpUrl+"/kline?"+param)
	content, err := base.HttpSend(config.GConfig.CoinMerit.HttpUrl+"/kline?"+param, nil, "GET", map[string]string{"ApiKey":config.GConfig.CoinMerit.ApiKey, "Timestamp":timeStamp, "ApiSign":sign })
	if err != nil {
		return nil,err
	}
	res := new(_type.ResCoinMeritKLine)
	if err = json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil,err
	}
	ZapLog().Info("res=", zap.Int("num", len(res.Data)))
	if res.Status_code != 200 {
		return nil, fmt.Errorf("%d %s", res.Status_code, res.Message)
	}
	return res,nil
}

//params应该是按照名称排序并且url编码的
func (this *CoinMeritModels) genSign(path, timestamp string, params string) string {
	md := md5.New()
	str := config.GConfig.CoinMerit.Secret_key +path+timestamp+params+config.GConfig.CoinMerit.Secret_key
//	fmt.Println(str)
	//newStr := md.Sum([]byte(str))
	md.Write([]byte(str))
	newStr := md.Sum(nil) //这个只是追加
//	fmt.Println(len(newStr))
	return strings.ToUpper(string(fmt.Sprintf("%X", newStr)))
}

type SubKeyInfo struct{
	Data   []byte     //发送消息，断线重发的时候使用
	SubMap map[string] *SubItem  //uuid, version
}

func (this *SubKeyInfo) Add(sessionId string, ver uint64){
	this.SubMap[sessionId] = &SubItem{ver}
}

func (this *SubKeyInfo) Del(sessionId string){
	delete(this.SubMap, sessionId)
}

func (this *SubKeyInfo) Len() int {
	return len(this.SubMap)
}

func NewSubKeyInfo() *SubKeyInfo{
	s := new(SubKeyInfo)
	s.SubMap = make(map[string] *SubItem)
	return s
}

type SubItem struct{
	Ver       uint64 //涉及到版本就是用
}

func NewCoinMeritModelsWsHandler(pushH common.PushHandler) *CoinMeritModelsWsHandler {
	c := new(CoinMeritModelsWsHandler)
	c.mRecvResChan = make(chan error, 5)
	c.mSubKeyMap = make(map[string] *SubKeyInfo)
	c.mPushHand = pushH
	return c
}

//handler
type CoinMeritModelsWsHandler struct{
	mRecvResChan  chan error
	mWsCon        *base.WsCon
	mSubKeyMap    map[string] *SubKeyInfo  //topic,
	ConId         string
	mPushHand         common.PushHandler
	mSendLock     sync.Mutex
	sync.Mutex
}

func (this * CoinMeritModelsWsHandler) IsRecvPingHander() bool {
	return true
}

func (this * CoinMeritModelsWsHandler) IsRecvPongHander() bool {
	return true
}

func (this * CoinMeritModelsWsHandler) RecvPingHander(message string) string {
	return ""
}

func (this * CoinMeritModelsWsHandler) RecvPongHander(message string) string {
	return ""
}

func (this * CoinMeritModelsWsHandler) SendPingHander() []byte {
	return []byte("{\"event\":\"ping\"}")
}

func (this * CoinMeritModelsWsHandler) BeforeSendHandler()  {
	this.mSendLock.Lock()
	ZapLog().Info("BeforeSendHandler lock")
	for {
		select{
		case  <-this.mRecvResChan:
			break
		default:
			return
		}
	}
}

func (this * CoinMeritModelsWsHandler) AfterSendHandler() error {
	defer this.mSendLock.Unlock()
	ZapLog().Info("AfterSendHandler Unlock")
	select{
		case err := <-this.mRecvResChan:
			return err
		case <- time.After(time.Second * 30):
			return errors.New("chan timeout")
	}
	return nil
}

func (this * CoinMeritModelsWsHandler)Test(){
	go func (){
		for true {
			time.Sleep(time.Second*5)
			ZapLog().Info("time to push")
			msg := new(_type.ResCoinMeritSub)
			msg.Topic = "cm_huobi_BCH_BTC_kline_1min"
			msg.Data = make([]_type.CoinMeritSubKXian, 0)
			kxian := new(_type.CoinMeritSubKXian)
			kxian.Pair = "hd_test"
			kxian.Data = []float64{1234567890, 12.3, 34.2, 34, 22, 11.3}
			msg.Data = append(msg.Data, *kxian)
			content,_ := json.Marshal(msg)
			this.RecvHandler("test", content)
		}
	}()

}

func (this * CoinMeritModelsWsHandler) RecvHandler(conId string, message []byte){
	res := new(_type.ResCoinMeritSub)
	err := json.Unmarshal(message, res)
	if err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("cm_sub_msg", string(message)), zap.String("conid", this.ConId))
		return
	}
	if res.Error_code != nil {
		if *res.Error_code == 0 {
			this.mRecvResChan <- nil
		}else{
			this.mRecvResChan <- fmt.Errorf("%d", *res.Error_code)
		}
		return
	}

	apiMsg := new(api.MSG)
	if strings.Contains(res.Topic, "kline") {
		apiKxian := CoinMeritSubToApiKXian(res)
		ZapLog().Info("recv=", zap.Any("res=", apiKxian))
		topicArr := strings.Split(res.Topic, "_")
		obj:=""
		if len(topicArr) >= 4 {
			obj = topicArr[2] +"_"+topicArr[3]
		}
		apiMsg.AddQuoteKlineSingle(api.NewQuoteKlineSingle(obj,apiKxian))
	}else{
		ZapLog().Info("GCoinMeritModels unknown sub msg", zap.String("topic", res.Topic), zap.String("conid", this.ConId))
		return
	}// else 各种消息类型
	//推送消息
	this.Lock()
	defer this.Unlock()
	subKeyInfoList, ok := this.mSubKeyMap[res.Topic]
	if !ok {
		ZapLog().Error("unknown topic",zap.String("topic", res.Topic), zap.String("conid", this.ConId))
		return
	}
	for k,v := range subKeyInfoList.SubMap {
		ZapLog().Info("push uuid ="+ k)
		base.GWsServerMgr.Push(k, v.Ver, apiMsg)
	}
}

//这里以后再优化吧,注意死锁
func (this * CoinMeritModelsWsHandler) ReSendAllRequest() {
	errMap := make(map[string] bool)
	subMap := this.GetSubMap()
	ZapLog().Info("start ReSendAllRequest", zap.Int("topicnum", len(subMap)))
	for topic,info := range subMap {
		ZapLog().Info("start topic "+topic)
		if err := this.mWsCon.Send(websocket.BinaryMessage, info); err != nil {
			ZapLog().Sugar().Errorf("ReSendAllRequest fail %v", topic)
			errMap[topic] = true
		}else{
			ZapLog().Info("ReSendAllRequest success", zap.String("topic", topic))
		}
		ZapLog().Info("end topic "+topic)
	}//重试3次
}

func (this *CoinMeritModelsWsHandler) GetSubMap() map[string] []byte {
	mm := make(map[string] []byte)
	this.Lock()
	defer this.Unlock()
	for topic,info := range this.mSubKeyMap {
		mm[topic] = info.Data
	}
	return mm
}

func (this * CoinMeritModelsWsHandler) Sub(topic, sessionId string, data []byte) bool {
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

func (this * CoinMeritModelsWsHandler) UnSub(topic, sessionId string) bool {
	this.Lock()
	defer this.Unlock()
	subKeyInfo, ok := this.mSubKeyMap[topic]
	if !ok {
		return false
	}
	subKeyInfo.Del(sessionId)
	if subKeyInfo.Len() == 0{
		delete(this.mSubKeyMap, topic)
		return true
	}else{
		return false
	}

}


func CoinMeritSubToApiKXian(k1 *_type.ResCoinMeritSub) []*api.KXian {
	k2 := make([]*api.KXian, 0)
	for j:=0; j< len(k1.Data);j++ {
		tmp := new(api.KXian)
		for i := 0; i < len(k1.Data[j].Data); i++ {
			switch i {
			case 0:
				tmp.SetShiJian(int64(k1.Data[j].Data[0]))
				break
			case 1:
				tmp.SetKaiPanJia(fmt.Sprintf("%f", k1.Data[j].Data[1]))
				break
			case 2:
				tmp.SetZuiGaoJia(fmt.Sprintf("%f", k1.Data[j].Data[2]))
				break
			case 3:
				tmp.SetZuiDiJia(fmt.Sprintf("%f", k1.Data[j].Data[3]))
				break
			case 4:
				tmp.SetShouPanJia(fmt.Sprintf("%f", k1.Data[j].Data[4]))
				break
			case 5:
				tmp.SetChengJiaoLiang(fmt.Sprintf("%f", k1.Data[j].Data[5]))
				break
			}
		}
		k2 = append(k2, tmp)
	}
	return k2
}