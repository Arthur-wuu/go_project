package base

import (
	"sync"
	"github.com/gorilla/websocket"
	"time"
	"net/http"
	"net/url"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"errors"
)

type WsMsg struct{
	Type int
	Data []byte
	ExpireAt int64
}

type WsCon struct {
	mUrl    *url.URL
	mHeader http.Header
	conn   *websocket.Conn
	exitCh    chan bool
	waitGroup sync.WaitGroup
	isAlive   bool
	dialer    *websocket.Dialer
	mSendCh   chan *WsMsg
	mRecvCh   chan *WsMsg
	mSendPingHander func() []byte
	mRecvPingHander func(string) []byte
	mRecvPongHander func(string) []byte
	sync.Mutex
}

func (this *WsCon) Init(urlStr string, sendPingHander func() []byte, recvPingHander, recvPongHander func(string) []byte) error {
	url,err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	this.mUrl = url
	//this.mHeader = header
	this.exitCh = make(chan bool)
	this.mSendCh =  make(chan *WsMsg, 1000)
	this.mRecvCh =  make(chan *WsMsg, 1000)
	this.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Second * 10,
	}

	this.mRecvPingHander = recvPingHander
	this.mRecvPongHander = recvPongHander
	this.mSendPingHander = sendPingHander

	if err := this.newCon();err != nil {
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
	for {
		select {
		case _,ok := <-this.mSendCh:
			if !ok {
				return errors.New("close sendCh err")
			}
			continue
		case _,ok := <-this.mRecvCh:
			if !ok {
				return errors.New("close recvCh err")
			}
			continue
		default:
			break
		}
		break
	}
	select{
		case this.mSendCh <- &WsMsg{messageType, data, time.Now().Add(time.Second*30).Unix()}:
		default:
			return errors.New("chan full err")
	}
	return nil
}

func (this *WsCon) Recv() (*WsMsg, error) {
	select{
	case  msg,ok := <- this.mRecvCh:
		if !ok {
			return nil, errors.New("recv closed chan err")
		}
		return msg,nil
	case <- time.After(time.Second * 30):
		return nil, errors.New("chan rcv time out ")
	}
	return nil, nil
}

func (this *WsCon) write(messageType int, data []byte) error {
	this.conn.SetWriteDeadline(time.Now().Add(time.Second *60))// 这个接口一次有效
	return this.conn.WriteMessage(messageType, data)
}

//一旦连接发送或者读取超时过了一次，往后发送均无效，也就是说连接已经费掉了
func (this *WsCon) runWrite() {
	this.waitGroup.Add(1)
	defer this.waitGroup.Done()
	defer this.conn.Close()  //让runRead顺利终止
	pingTick := time.NewTimer(time.Second * 60)

	for {//超时错误怎么办
		select {
		case msg := <-this.mSendCh:
			nowTime := time.Now().Unix()
			if msg.ExpireAt <  nowTime{
				//消息超时
				ZapLog().Warn("msg to much timetout", zap.Int64("time", nowTime-msg.ExpireAt))
				break
			}
			if err := this.write(msg.Type, msg.Data); err != nil {
				ZapLog().Error("ws write err", zap.Error(err))
				return
			}
			pingTick.Reset(time.Second * 60)
			time.Sleep(time.Second * 10)
		case <-this.exitCh:
			return
		case <-pingTick.C:
			ZapLog().Debug("time to send ping")
			content := this.mSendPingHander()
			if err := this.conn.WriteControl(websocket.PingMessage, content, time.Now().Add(time.Second * 10)); err != nil {
				ZapLog().Error("ws WriteControl err", zap.Error(err))
				return
			}
			time.Sleep(time.Second *10)
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
			ZapLog().Error("ws ReadMessage err", zap.Error(err))
			return
		}
		this.conn.SetReadDeadline(time.Now().Add(time.Second * 22)) //=sendPing*2+2, 这个接口一次有效
		switch msgType {
		case websocket.CloseMessage:
			ZapLog().Error("ws ReadMessage CloseMessage")
			return
		case websocket.BinaryMessage, websocket.TextMessage:
			this.mRecvCh <- &WsMsg{msgType, data, 0}
		}
	}
}

func (this *WsCon) newCon() error {
	//ZapLog().Info("newCon")
	c, _, err := this.dialer.Dial(this.mUrl.String(), this.mHeader)
	if err != nil {
		ZapLog().Error("newCon err", zap.Error(err))
		return err
	}
	this.Lock()
	defer this.Unlock()
	this.conn = c
	this.setPingPongHandler()
	return nil
}

func (this *WsCon) setPingPongHandler() {
		this.conn.SetPingHandler(func(appData string) error {
			this.conn.SetReadDeadline(time.Now().Add(time.Second * 22))
			ZapLog().Debug("recv ping and send pong")
			newData := this.mRecvPingHander(appData)
			if newData == nil || len(newData) == 0 {
				return nil
			}
			err := this.conn.WriteControl(websocket.PongMessage, []byte(newData), time.Now().Add(time.Second*5))
			if err != nil {
				ZapLog().Error("WriteControl PongMessage err", zap.Error(err))
			}
			time.Sleep(time.Second * 10)
			return nil
		})

		this.conn.SetPongHandler(func(appData string) error {
			this.conn.SetReadDeadline(time.Now().Add(time.Second * 22))
			ZapLog().Debug("recv  pong")
			this.mRecvPongHander(appData)
			return nil
		})


}

func (this *WsCon) reConnect() {
	if !this.isAlive {
		return
	}
	ZapLog().Info("Find Disconnect and will reConnect")
	close(this.exitCh) //这里让 runwrite退出
	this.waitGroup.Wait() //必须要等到两个线程 退出
	this.exitCh = make(chan bool)

	for {
		err := this.newCon()
		if err != nil {
			ZapLog().Error("newCon err so sleep 120s")
//			logger.Error("PoccSdk:: connect to pocc failed, err(%v)", err)
			time.Sleep(120 * time.Second)
		} else {
			this.Start()
			//重新推
			ZapLog().Info("reConnect ok and will reSub all topic")
			return
		}
	}
}

func (this *WsCon) Stop() {
	this.Lock()
	defer this.Unlock()
	if this.conn != nil {
		this.isAlive = false
		this.conn.WriteMessage(websocket.CloseMessage, nil)
//		logger.Info("PoccSdk::send close msg to pocc done.")
		close(this.exitCh)
		time.Sleep(time.Second*1)
		this.waitGroup.Wait()
		this.conn = nil
	}
	ZapLog().Info("wscon stop")
}


func  SendPingHander () []byte{

	return  []byte("ping")
}

func  RecvPingHander ( str string) []byte {

	return []byte(str)
}

func  RecvPongHander ( str string) []byte {

	return []byte(str)
}