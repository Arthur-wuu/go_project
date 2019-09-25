package base

import (
	"net/url"
	"net/http"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)


func NewWsMgr() *WsMgr {
	w := new(WsMgr)
	w.Init()
	return w
}

type WsMgr struct {
	mWsCons map[string] *WsCon
	sync.Mutex
}

func (this *WsMgr) Init() {
	this.mWsCons = make(map[string] *WsCon)
}

func (this *WsMgr) AddCon(conId string, wshand WsHandler, ul *url.URL, head http.Header) error {
	this.Lock()
	c,ok := this.mWsCons[conId];
	this.Unlock()
	if ok {
		return errors.New("wscon_id exist")
	}

	c = new(WsCon)
//	c.SetHandlers(this.dealResponse, this.mRecvPingHandler, this.mRecvPongHandler, this.mSendPingHandler)
	if err := c.Init(conId, ul, head, wshand);err != nil {
		return err
	}
	c.Start()
	this.Lock()
	this.mWsCons[conId] = c
	this.Unlock()
	return nil
}

func (this *WsMgr) RemoveCon(conId string) error {
	this.Lock()
	con,ok := this.mWsCons[conId];
	this.Unlock()
	if !ok {
		return nil
	}
	con.stop()
	this.Lock()
	delete(this.mWsCons, conId)
	this.Unlock()
	return nil
}


func (this *WsMgr) Send(conId string, data []byte) error {
	this.Lock()
	con, ok := this.mWsCons[conId];
	this.Unlock()
	if !ok {
		return errors.New("wscon_id nofind")
	}
	if err := con.Send(websocket.BinaryMessage, data); err != nil {
		return err
	}
	return nil
}

func (this *WsMgr) Stop() {
	this.Lock()
	defer this.Unlock()
	for _,v :=range this.mWsCons {
		v.stop()
	}
}
