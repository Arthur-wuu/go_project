package base

import (
	sws "github.com/kataras/iris/websocket"
	"BastionPay/bas-tv-proxy/config"
	cws "github.com/gorilla/websocket"
	"net/url"

	"BastionPay/bas-tv-proxy/api"

	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"github.com/pborman/uuid"
	"sync"
	"errors"
	"encoding/json"
	"strings"
)

var GWsServerMgr WsServerMgr

type WsServerMgr struct {
	mConf *config.Config
	mInConGroup map[string] *WsServer   //conId
	reqUuidMap  map[string] Requester //requuid==conId+qid
	mWsHandlers  map[string] func(request Requester)
	mPackPushHandler PackPushHanderType
	mPackResHandler PackResHanderType
	mNewRequesterHandler NewRequesterHandler
	sync.Mutex
}

func (this *WsServerMgr) Init(c *config.Config) {
	this.mConf = c
	this.mInConGroup = make(map[string] *WsServer)
	this.mWsHandlers = make(map[string] func(request Requester))
	this.reqUuidMap = make(map[string] Requester)
	if this.mPackPushHandler == nil {
		this.mPackPushHandler = this.DefaultPackPushMsg
	}
	if this.mPackResHandler == nil {
		this.mPackResHandler = this.DefaultPackResMsg
	}
	if this.mNewRequesterHandler == nil {
		this.mNewRequesterHandler = this.DefaultNewRequester
	}
}

func (this *WsServerMgr) RegPackHander(pushH PackPushHanderType, resH PackResHanderType) {
	this.mPackPushHandler = pushH
	this.mPackResHandler = resH
}

func (this *WsServerMgr) RegNewRequesterHandler(h NewRequesterHandler) {
	this.mNewRequesterHandler = h
}

func (this *WsServerMgr) RegHandler(path string,f func(request Requester)) {
	this.mWsHandlers[path] = f
}

func (this *WsServerMgr) AddReq(req Requester) error {
	this.Lock()
	defer this.Unlock()
	ws ,ok := this.mInConGroup[req.GetConId()]
	if !ok {
		return errors.New("nofind coinId")
	}
	if _,ok := this.reqUuidMap[req.GetUuid()]; ok {
		return errors.New("exist uuid")
	}
	if err :=ws.AddReq(req); err != nil {
		return err
	}
	this.reqUuidMap[req.GetUuid()] = req
	ZapLog().Info("AddReq ", zap.Int("reqCount", len(this.reqUuidMap)), zap.String("coinid", req.GetConId()))
	return nil
}

func (this *WsServerMgr) RemoveReq(conId, qid string) {
	this.Lock()
	defer this.Unlock()
	ws ,ok := this.mInConGroup[conId]
	if !ok {
		return
	}
	req := ws.GetRequester(qid)
	if req == nil {
		return
	}
	ws.RemoveReq(req.GetQid())
	delete(this.reqUuidMap, req.GetUuid())
	ZapLog().Info("RemoveReq ", zap.Int("reqCount", len(this.reqUuidMap)), zap.String("coinid", req.GetConId()))
}

func (this *WsServerMgr) AddConnection( ws *WsServer) (int, bool) {
	this.Lock()
	defer this.Unlock()
	_,ok := this.mInConGroup[ws.GetId()]
	this.mInConGroup[ws.GetId()] = ws
	return len(this.mInConGroup), ok
}

func (this *WsServerMgr) DelConnection(ws *WsServer) (int,bool) {
	this.Lock()
	defer this.Unlock()
	_,ok := this.mInConGroup[ws.GetId()]
	delete(this.mInConGroup, ws.GetId())
	return len(this.mInConGroup), ok
}



/**********ws*************/
func (this *WsServerMgr) HandleWsConnection(con sws.Connection){
	ZapLog().Info("Ws Connect Start "+con.ID())
	ws := NewWsServer(con)
	ws.Start()
	allcount,addExistFlag := this.AddConnection(ws)
	if addExistFlag {
		ZapLog().Error("Ws AddConnection but Id exist "+con.ID())
	}
	ZapLog().Info("Ws Connect Success "+con.ID(), zap.Int("allcount", allcount))
	con.OnMessage(func(data []byte){
		ZapLog().Info("con.OnMessage ", zap.String("req", string(data)), zap.String("conid", con.ID()))
		requrl, err := url.Parse(string(data))
		if err != nil {
			ZapLog().Error("url Parse err", zap.String("req", string(data)), zap.Error(err), zap.String("conid", con.ID()))
			this.Json(con, this.mNewRequesterHandler(ws, nil, "",this.RemoveReq, this.mPackPushHandler, this.mPackResHandler), api.ErrCode_UrlPath, nil)
			return
		}
		requester := this.mNewRequesterHandler(ws, requrl, uuid.New(), this.RemoveReq, this.mPackPushHandler, this.mPackResHandler)

		if requester.IsSub() { //先设置，万一失败了再删除，这样做可以让防止推送数据丢失
			if len(requester.GetQid()) == 0 {
				ZapLog().Error("param no Qid in sub", zap.String("req", string(data)), zap.String("conid", con.ID()))
				this.Json(con, requester, api.ErrCode_UrlPath, nil)
				return
			}
			if err = this.AddReq(requester); err != nil {
				ZapLog().Error("AddReq err", zap.String("req", string(data)), zap.String("uuid", requester.GetUuid()), zap.Error(err), zap.String("conid", con.ID()))
				this.Json(con, requester, api.ErrCode_InerServer, nil)
				return
			}
			ZapLog().Info("sub Msg ", zap.String("req", string(data)))
		}
		if requester.IsUnSub() { //反正你要取消订阅的，不管能不能成功，先取消再说
			if len(requester.GetQid()) == 0 {
				ZapLog().Error("param no Qid in unsub", zap.String("req", string(data)), zap.String("conid", con.ID()))
				this.Json(con, requester, api.ErrCode_UrlPath, nil)
				return
			}
			oldReqster := ws.GetRequester(requester.GetQid())
			if oldReqster == nil {
				ZapLog().Error("unsub not find old requester", zap.String("req", string(data)), zap.String("conid", con.ID()))
				this.Json(con, requester, api.ErrCode_UrlPath, nil)
				return
			}
			requester = oldReqster //应为取消订阅可能只传个qid过来，参数都没传，所以还得用老的会话信息
			requester.CancelSub()
			this.RemoveReq(con.ID(), requester.GetQid())
			ZapLog().Info("unsub Msg ", zap.String("req", string(data)))
		}

		handler, ok := this.mWsHandlers[requester.GetPath()]
		if !ok {
			ZapLog().Error("nofind path", zap.String("req", string(data)), zap.String("conid", con.ID()))
			this.Json(con, requester, api.ErrCode_UrlPath, nil)
			return
		}

		go handler(requester) //后续操作 走网络的，所以异步处理
	})
	con.OnDisconnect(func(){//同一个连接多次收到断开消息，坑爹
		ZapLog().Info("Ws DisConnect Start "+con.ID())
		allcount,existFlag := this.DelConnection(ws)
		if !existFlag {
			ZapLog().Error("Ws DisConnect not exist "+con.ID())
			return
		}
		ws.Stop()
		mm := ws.GetAllRequester()
		for _,v := range mm {
			if !v.IsSub() {
				continue
			}
			this.Lock()
			delete(this.reqUuidMap, v.GetUuid())
			this.Unlock()
			v.CancelSub()
			hand, ok := this.mWsHandlers[v.GetPath()]
			if !ok {
				continue
			}
			go hand(v)
		}
		ZapLog().Info("Ws DisConnect Success "+con.ID(), zap.Int("allcount", allcount))
	})
}

//持续推送给客户端, reqUUid 比起conid通用性更高
func (this *WsServerMgr) Push(reqUUid string, ver uint64, apiMsg interface{}) {
	this.Lock()
	requester, ok := this.reqUuidMap[reqUUid]
	this.Unlock()
	if !ok {
		ZapLog().Warn("nofind reqUUid", zap.String("reqUUid", reqUUid))
		return
	}

	if err := requester.OnPushResponseWithPack(apiMsg); err != nil {
		ZapLog().Error("OnPushResponseWithPack err", zap.String("reqUUid", reqUUid), zap.Error(err))
		return
	}
}

func (this *WsServerMgr) Json(con sws.Connection,reqster Requester, errCode int32, apiMsg interface{}) {
	content, err := this.mPackResHandler(reqster, errCode, apiMsg)
	if err != nil {
		ZapLog().Error("response Marshal err", zap.Error(err))
		return
	}
	if err =con.Write(cws.TextMessage, content); err != nil {//con自带锁
		ZapLog().Error("Ws Write err", zap.Error(err))
		return
	}
	return
}

func (this *WsServerMgr) DefaultPackPushMsg(reqster Requester, data interface{}) ([]byte, error) {
	content, ok := data.([]byte)
	if ok {
		return content,nil
	}
	return json.Marshal(data)
}

//错误和正常返回 都包含了
func (this *WsServerMgr) DefaultPackResMsg(reqster Requester, errCode int32, data interface{}) ([]byte, error) {
	content, ok := data.([]byte)
	if ok {
		return content,nil
	}
	return json.Marshal(data)
}

func (this *WsServerMgr) DefaultNewRequester(wsServer *WsServer, urlReq *url.URL, reqUUID string, clearHand ClearHanderType, packPushHand PackPushHanderType ,packResHand PackResHanderType) Requester {
	reqNew := new(DefaultRequester)
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
	}else{
		sub := reqNew.Values.Get("sub")
		if len(sub) == 0 || sub >= api.EVENT_SubMax{
			sub = api.EVENT_nonesub
		}
		reqNew.Values.Set("sub", sub)
		reqNew.Sub = sub
	}
	reqNew.counter = 0
	return reqNew
}