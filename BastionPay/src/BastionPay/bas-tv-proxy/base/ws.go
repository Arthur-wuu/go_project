package base

import (
	. "BastionPay/bas-base/log/zap"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type WsMsg struct {
	Type     int
	Data     []byte
	ExpireAt int64
}

type WsCon struct {
	mId       string
	mUrl      *url.URL
	mHeader   http.Header
	conn      *websocket.Conn
	exitCh    chan bool
	waitGroup sync.WaitGroup
	isAlive   bool
	dialer    *websocket.Dialer
	mHander   WsHandler
	mSendCh   chan *WsMsg
	mRecvCh   chan *WsMsg
	sync.Mutex
}

func (this *WsCon) Init(id string, url *url.URL, header http.Header, wshand WsHandler) error {
	this.mUrl = url
	this.mHeader = header
	this.mId = id
	this.mHander = wshand
	this.exitCh = make(chan bool)
	this.mSendCh = make(chan *WsMsg, 1000)
	this.mRecvCh = make(chan *WsMsg, 1000)
	this.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Second * 10,
	}

	if err := this.newCon(); err != nil {
		return err
	}
	return nil
}

func (this *WsCon) Start() {
	this.isAlive = true
	go this.runRead()
	go this.runWrite()
	return
}

func (this *WsCon) Send(messageType int, data []byte) error {
	this.mHander.BeforeSendHandler()
	this.mSendCh <- &WsMsg{messageType, data, time.Now().Add(time.Second * 25).Unix()}
	if err := this.mHander.AfterSendHandler(); err != nil {
		return err
	}
	return nil
}

func (this *WsCon) write(messageType int, data []byte) error {
	this.conn.SetWriteDeadline(time.Now().Add(time.Second * 15)) // 这个接口一次有效
	return this.conn.WriteMessage(messageType, data)
}

//一旦连接发送或者读取超时过了一次，往后发送均无效，也就是说连接已经费掉了
func (this *WsCon) runWrite() {
	this.waitGroup.Add(1)
	defer this.waitGroup.Done()
	defer this.conn.Close() //让runRead顺利终止
	pingTick := time.Tick(time.Second * 10)
	for { //超时错误怎么办
		select {
		case msg := <-this.mSendCh:
			nowTime := time.Now().Unix()
			if msg.ExpireAt < nowTime {
				//消息超时
				ZapLog().Warn("msg to much timetout", zap.Int64("time", nowTime-msg.ExpireAt), zap.String("conid", this.mId))
				break
			}
			if err := this.write(msg.Type, msg.Data); err != nil {
				ZapLog().Error("ws write err", zap.Error(err), zap.String("conid", this.mId))
				return
			}
		case <-this.exitCh:
			return
		case <-pingTick:
			ZapLog().Debug("time to send ping", zap.String("conid", this.mId))
			content := this.mHander.SendPingHander()
			if err := this.conn.WriteControl(websocket.PingMessage, content, time.Now().Add(time.Second*6)); err != nil {
				ZapLog().Error("ws WriteControl err", zap.Error(err), zap.String("conid", this.mId))
				return
			}
			ZapLog().Debug("send ping ok")

		}
	}
}

//read 比wirte判断网络异常的概率大多了
func (this *WsCon) runRead() {
	defer this.reConnect()
	this.waitGroup.Add(1)
	defer this.waitGroup.Done()
	/*一旦连接发送或者读取超时过了一次，往后发送或者读取均无效，也就是说连接已经费掉了。
	    连接断开后，sendping 不会报错
		ReadMessage永远阻塞，连接断开后ReadMessage要60+秒（30秒发哦5分钟不等）才能返回，而且报错 timeout，这种方式显然不合适。
	    ReadMessage带超时，要配合ping pong，在收到ping和pong的时候重新设置读超时。这种方式断开后最终返回是超时。
	    SetReadDeadline
	*/
	for {
		//阻塞模式，不会有超时或者信号打断错误
		msgType, data, err := this.conn.ReadMessage()
		if err != nil {
			ZapLog().Error("ws ReadMessage err", zap.Error(err), zap.String("conid", this.mId))
			return
		}
		this.conn.SetReadDeadline(time.Now().Add(time.Second * 22)) //=sendPing*2+2, 这个接口一次有效
		switch msgType {
		case websocket.CloseMessage:
			ZapLog().Error("ws ReadMessage CloseMessage", zap.String("conid", this.mId))
			return
		case websocket.BinaryMessage, websocket.TextMessage:
			this.mHander.RecvHandler(this.mId, data)
		}
	}
}

func (this *WsCon) newCon() error {
	ZapLog().Info("newCon", zap.String("conid", this.mId))
	c, _, err := this.dialer.Dial(this.mUrl.String(), this.mHeader)
	if err != nil {
		ZapLog().Error("newCon err", zap.Error(err), zap.String("conid", this.mId))
		return err
	}
	this.Lock()
	defer this.Unlock()
	this.conn = c
	this.setPingPongHandler()
	return nil
}

func (this *WsCon) setPingPongHandler() {
	if this.mHander.IsRecvPingHander() {
		this.conn.SetPingHandler(func(appData string) error {
			this.conn.SetReadDeadline(time.Now().Add(time.Second * 22))
			ZapLog().Debug("recv ping and send pong")
			newData := this.mHander.RecvPingHander(appData)
			err := this.conn.WriteControl(websocket.PongMessage, []byte(newData), time.Now().Add(time.Second*5))
			if err != nil {
				ZapLog().Error("WriteControl PongMessage err", zap.Error(err), zap.String("conid", this.mId))
			}
			return nil
		})
	}
	if this.mHander.IsRecvPongHander() {
		this.conn.SetPongHandler(func(appData string) error {
			this.conn.SetReadDeadline(time.Now().Add(time.Second * 22))
			ZapLog().Debug("recv  pong", zap.String("conid", this.mId))
			this.mHander.RecvPongHander(appData)
			return nil
		})
	}

}

func (this *WsCon) reConnect() {
	if !this.isAlive {
		return
	}
	ZapLog().Info("Find Disconnect and will reConnect", zap.String("conid", this.mId))
	close(this.exitCh)    //这里让 runwrite退出
	this.waitGroup.Wait() //必须要等到两个线程 退出
	this.exitCh = make(chan bool)

	for {
		err := this.newCon()
		if err != nil {
			ZapLog().Error("newCon err so sleep 120s", zap.String("conid", this.mId))
			//			logger.Error("PoccSdk:: connect to pocc failed, err(%v)", err)
			time.Sleep(120 * time.Second)
		} else {
			this.Start()
			//重新推
			ZapLog().Info("reConnect ok and will reSub all topic", zap.String("conid", this.mId))
			this.mHander.ReSendAllRequest()
			return
		}
	}
}

func (this *WsCon) stop() {
	this.Lock()
	defer this.Unlock()
	if this.conn != nil {
		this.isAlive = false
		this.conn.WriteMessage(websocket.CloseMessage, nil)
		//		logger.Info("PoccSdk::send close msg to pocc done.")
		close(this.exitCh)
		this.waitGroup.Wait()
		this.conn = nil
	}
	ZapLog().Info("wscon stop" + this.mId)
}

func (this *WsCon) GetId() string {
	return this.mId
}
