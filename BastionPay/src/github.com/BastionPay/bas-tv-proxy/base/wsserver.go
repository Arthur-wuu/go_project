package base
import(
	sws "github.com/kataras/iris/websocket"
	"sync"
	"github.com/gorilla/websocket"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"errors"
	"net"
)

func NewWsServer(c sws.Connection) *WsServer {
	w := new(WsServer)
	w.mCon = c
	w.reqIdMap = make(map[string] Requester)
	w.mExitCh = make(chan bool)
	w.chResp = make(chan []byte, 1000)
	return w
}
type WsServer struct{
	mCon           sws.Connection
	reqIdMap       map[string] Requester //qid, request
	mRunFlag       bool
	chResp         chan []byte //单个线程 轮训写入多个ws时，采用chan+ws线程速度更快
	mExitCh        chan bool
	mWaitGroup     sync.WaitGroup
	sync.Mutex
}

func (this *WsServer) AddReq(req Requester) error {
	this.Lock()
	defer this.Unlock()
	_,ok := this.reqIdMap[req.GetQid()]
	if ok {
		return errors.New("exist qid")
	}
	this.reqIdMap[req.GetQid()] =req
	return nil
}

func (this *WsServer) RemoveReq(qid string){
	this.Lock()
	defer this.Unlock()
	delete(this.reqIdMap, qid)
}

func (this *WsServer) GetRequester(qid string) Requester {
	this.Lock()
	defer this.Unlock()
	req,_ :=  this.reqIdMap[qid]
	return req
}

func (this *WsServer) GetAllRequester() map[string]Requester {
	this.Lock()
	defer this.Unlock()
	return  this.reqIdMap
}

//适合推送
func (this *WsServer) Send(data []byte) {
	if this.mRunFlag == false {
		return
	}
	select{
	case this.chResp <- data:
		break
	case <-this.mExitCh://防止ws线程退出，chResp没关，然后导致工作线程阻塞
		break
	}

}

func (this *WsServer) Start() {
	this.mRunFlag = true
	go this.run()
}

func (this *WsServer) run() error { //只需要写，读由框架做

	this.mWaitGroup.Add(1)
	defer this.mWaitGroup.Done()
	for {
		select {
		case resp, ok := <-this.chResp:
			if !ok {
				ZapLog().Error("chResp have closed", zap.String("conid", this.GetId()))
				if err := this.mCon.Write(websocket.CloseMessage, []byte{}); err != nil {
					if err == websocket.ErrCloseSent {
						return nil
					} else if e, ok := err.(net.Error); ok && e.Temporary() {
						return nil
					}
					ZapLog().Error("write CloseMessage err", zap.Error(err), zap.String("conid", this.GetId()))
				}
				return nil
			}
			if err := this.mCon.Write(websocket.TextMessage, resp); err != nil {
				if err != websocket.ErrCloseSent {//或者用this.mRunFlag, wsserverMgr 会判断断开的
					ZapLog().Error("write err", zap.Error(err), zap.String("conid", this.GetId()))
				}
				return err
			}
			//是否取消
		case <-this.mExitCh:
			if err := this.mCon.Write(websocket.CloseMessage, []byte{}); err != nil {
				if err == websocket.ErrCloseSent {
					return nil
				} else if e, ok := err.(net.Error); ok && e.Temporary() {
					return nil
				}
				ZapLog().Error("Write CloseMessage err" , zap.Error(err)  , zap.String("conid", this.GetId()))
			}
			return nil
		}
	}
	return nil
}

func (this *WsServer) Stop() error {
	if !this.mRunFlag {
		ZapLog().Error("WsServer stop more than one", zap.String("conid", this.GetId()))
		return nil
	}
	this.mRunFlag = false
	close(this.mExitCh)
	this.mWaitGroup.Wait()
	close(this.chResp)
	ZapLog().Info("WsServer stop", zap.String("conid", this.GetId()) )
	return nil
}

func (this *WsServer) GetId() string {
	return this.mCon.ID()
}