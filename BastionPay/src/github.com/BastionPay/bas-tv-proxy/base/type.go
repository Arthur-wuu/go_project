package base

import (
	"net/url"
	"strings"
	"BastionPay/bas-tv-proxy/api"
)

type WsHandler interface{
	RecvPingHander(message string) string
	RecvPongHander(message string) string
	SendPingHander() []byte
	AfterSendHandler() error
	BeforeSendHandler()
	RecvHandler(conId string, message []byte)
	IsRecvPingHander() bool
	IsRecvPongHander() bool
	ReSendAllRequest()
}

type ClearHanderType func(conId, qid string)
type PackPushHanderType func(reqster Requester, data interface{}) ([]byte, error)
type PackResHanderType func(reqster Requester, errCode int32, data interface{}) ([]byte,error)

type Requester interface {
	SetUuid(uuid string)
	GetUuid() string
	GetCounter() uint64
	GetFirstPath()string
	GetPath() string
	GetParamValue(param string) (string)
	IsSimpleReq() bool
	IsSub() bool
	IsUnSub() bool
	CancelSub()
	OnResponseWithPack(errCode int32, apiMag interface{}) error
	OnResponse(errCode int32, data []byte ) error
	SubErrClear(errCode int32)
	OnPushResponseWithPack(data interface{}) error
	OnPushResponse(data []byte) error
	GetQid() string
	GetConId() string
	GetSub() string
}

type NewRequesterHandler func(wsServer *WsServer, urlReq *url.URL, reqUUID string, clearHand ClearHanderType, packPushHand PackPushHanderType ,packResHand PackResHanderType) Requester

type DefaultRequester struct {
	Qid       string    //同一连接 唯一，请求推送的会话id
	Uuid      string    //程序唯一
	Url       *url.URL
	Values    url.Values
	ConId     string     //程序唯一
	Sub       string
	counter   uint64       //会话递增，0表示首次，>=1表示推送
	WsServer  *WsServer
	clearHander func(conId, qid string)
	packPushHandler  PackPushHanderType
	packResHandler  PackResHanderType
}

func (this *DefaultRequester) GetSub() string {
	return this.Sub
}

func (this *DefaultRequester) GetConId() string {
	return this.ConId
}

func (this *DefaultRequester) GetQid() string {
	return this.Qid
}

func (this *DefaultRequester) SetUuid(uuid string){
	this.Uuid = uuid
}

func (this *DefaultRequester) GetUuid() string {
	return this.Uuid
}

func (this *DefaultRequester) GetCounter() uint64 {
	this.counter++
	return this.counter
}

func (this *DefaultRequester) GetFirstPath()string {
	paths := strings.Split(this.Url.Path, "/")
	if len(paths) >= 4 {
		return paths[3]
	}
	return ""
}

func (this *DefaultRequester) GetPath() string {
	return this.Url.Path
}

func (this *DefaultRequester) GetParamValue(param string) (string) {
	return this.Values.Get(param) //这个效率很低啊
}

func (this *DefaultRequester) IsSimpleReq() bool {
	return this.Sub == api.EVENT_nonesub
}

func (this *DefaultRequester) IsSub() bool {
	return this.Sub == api.EVENT_sub
}

func (this *DefaultRequester) IsUnSub() bool {
	return this.Sub == api.EVENT_unsub || strings.HasSuffix(this.Url.Path, "cancel")
}

func (this *DefaultRequester) CancelSub() {
	this.Sub = api.EVENT_unsub
	this.Values.Set("sub", api.EVENT_unsub)
	//this.Url.Query().Set("sub", api.EVENT_unsub)这部操作无效
}

func (this *DefaultRequester) OnResponseWithPack(errCode int32, apiMag interface{}) error {
	if errCode != api.ErrCode_Success {
		this.SubErrClear(errCode)
	}
	content, err := this.packResHandler(this, errCode, apiMag)
	if err != nil {
		return err
	}
	this.WsServer.Send(content)
	return nil
}

func (this *DefaultRequester) OnResponse(errCode int32, data []byte ) error {
	if errCode != api.ErrCode_Success {
		this.SubErrClear(errCode)
	}
	this.WsServer.Send(data)
	return nil
}

func (this *DefaultRequester) SubErrClear(errCode int32) {
	if this.IsSub() && errCode != api.ErrCode_Success && this.clearHander != nil{
		this.clearHander(this.ConId, this.Qid)
	}
}

func (this *DefaultRequester) OnPushResponseWithPack(data interface{}) error {
	content, err := this.packPushHandler(this, data)
	if err != nil {
		return err
	}
	this.WsServer.Send(content)
	return nil
}

func (this *DefaultRequester) OnPushResponse(data []byte) error {
	this.WsServer.Send(data)
	return nil
}
