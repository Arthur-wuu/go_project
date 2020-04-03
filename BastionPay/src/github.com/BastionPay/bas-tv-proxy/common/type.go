package common

import (
	"BastionPay/bas-tv-proxy/api"
	"BastionPay/bas-tv-proxy/base"
	"net/url"
	"strings"
)

type PushHandler interface {
	Push(reqUUid string, ver uint64, apiMsg interface{})
}

type WsServerInf interface {
	GetId() string
	Send(data []byte)
}

//请求定义
type Requester struct {
	Qid             string //同一连接 唯一，请求推送的会话id
	Uuid            string //程序唯一
	Url             *url.URL
	Values          url.Values
	ConId           string //程序唯一
	Sub             string
	counter         uint64 //会话递增，0表示首次，>=1表示推送
	WsServer        *base.WsServer
	clearHander     func(conId, qid string)
	packPushHandler base.PackPushHanderType
	packResHandler  base.PackResHanderType
}

func NewRequester(wsServer *base.WsServer, urlReq *url.URL, reqUUID string, clearHand base.ClearHanderType, packPushHand base.PackPushHanderType, packResHand base.PackResHanderType) base.Requester {
	//创建新请求
	reqNew := new(Requester)
	reqNew.ConId = wsServer.GetId()
	reqNew.Uuid = reqUUID
	reqNew.WsServer = wsServer
	reqNew.Url = urlReq
	reqNew.Values = urlReq.Query()
	reqNew.Qid = urlReq.Query().Get("qid")
	reqNew.clearHander = clearHand
	reqNew.packPushHandler = packPushHand
	reqNew.packResHandler = packResHand
	if strings.HasSuffix(urlReq.Path, "cancel") {
		sub := api.EVENT_nonesub
		reqNew.Values.Set("sub", sub)
		reqNew.Sub = sub
	} else {
		sub := reqNew.Values.Get("sub")
		if len(sub) == 0 || sub >= api.EVENT_SubMax {
			sub = api.EVENT_nonesub
		}
		reqNew.Values.Set("sub", sub)
		reqNew.Sub = sub
	}
	reqNew.counter = 0
	return reqNew
}

func (this *Requester) GetSub() string {
	return this.Sub
}

func (this *Requester) GetConId() string {
	return this.ConId
}

func (this *Requester) GetQid() string {
	return this.Qid
}

func (this *Requester) SetUuid(uuid string) {
	this.Uuid = uuid
}

func (this *Requester) GetUuid() string {
	return this.Uuid
}

func (this *Requester) GetCounter() uint64 {
	this.counter++
	return this.counter
}

func (this *Requester) GetFirstPath() string {
	paths := strings.Split(this.Url.Path, "/")
	if len(paths) >= 4 {
		return paths[3]
	}
	return ""
}

func (this *Requester) GetPath() string {
	return this.Url.Path
}

func (this *Requester) GetParamValue(param string) string {
	return this.Values.Get(param) //这个效率很低啊
}

func (this *Requester) IsSimpleReq() bool {
	return this.Sub == api.EVENT_nonesub
}

func (this *Requester) IsSub() bool {
	return this.Sub == api.EVENT_sub
}

func (this *Requester) IsUnSub() bool {
	return this.Sub == api.EVENT_unsub || strings.HasSuffix(this.Url.Path, "cancel")
}

func (this *Requester) CancelSub() {
	this.Sub = api.EVENT_unsub
	this.Values.Set("sub", api.EVENT_unsub)
	//this.Url.Query().Set("sub", api.EVENT_unsub)这部操作无效
}

func (this *Requester) OnResponseWithPack(errCode int32, apiMag interface{}) error {
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

func (this *Requester) OnResponse(errCode int32, data []byte) error {
	if errCode != api.ErrCode_Success {
		this.SubErrClear(errCode)
	}
	this.WsServer.Send(data)
	return nil
}

func (this *Requester) SubErrClear(errCode int32) {
	if this.IsSub() && errCode != api.ErrCode_Success && this.clearHander != nil {
		this.clearHander(this.ConId, this.Qid)
	}
}

func (this *Requester) OnPushResponseWithPack(data interface{}) error {
	content, err := this.packPushHandler(this, data)
	if err != nil {
		return err
	}
	this.WsServer.Send(content)
	return nil
}

func (this *Requester) OnPushResponse(data []byte) error {
	this.WsServer.Send(data)
	return nil
}

//func NewCtxRequester(ctx iris.Context) *CtxRequester {
//
//}
//
//type CtxRequester struct{
//	Ctx iris.Context
//}
//
//func ()
