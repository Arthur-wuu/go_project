package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-im-assist/api"
	"BastionPay/bas-wallet-im-assist/base"
	"BastionPay/bas-wallet-im-assist/config"
	"BastionPay/bas-wallet-im-assist/db"
	"BastionPay/bas-wallet-im-assist/tencent"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
	"go.uber.org/zap"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	//向Tencent查询用户状态的url
	tencentUrl = "https://console.tim.qq.com/v4/openim/querystate?usersig=%s&identifier=%s&sdkappid=%s&random=%s&contenttype=json"

	//向bastion 发送红点通知的url
	bastionPayUrl = "/push_message?"
)

var GTasker Tasker
var Sign string

type Tasker struct {
	mSinglechatChan  chan *api.SingleChat
	mSendNotifyLimit *limiter.Limiter //一分钟一次 通知BastionPay 红点消息

}

type Message struct {
	Content     Content `json:"content"`
	CreateTime  int32   `json:"createTime"`
	HistoryFlag bool    `json:"historyFlag"`
	Id          int32   `json:"id"`
	MsgType     int32   `json:"msgType"`
	RedFlag     bool    `json:"redFlag"`
	SourceId    string  `json:"sourceId"`
	SubType     string  `json:"subType"`
	Title       Title   `json:"title"`
	UserId      int     `json:"userId"`
}

type Content struct {
	SenderId  string `json:"senderId"`
	ReceverId string `json:"receverId"`
	Content   string `json:"content"`
}

type Title struct {
	Zh_CN_title string `json:"zh_CN_title"`
	En_US_title string `json:"en_US_title"`
	Zh_TW_title string `json:"zh_TW_title"`
	Jpa_title   string `json:"jpa_title"`
	Spa_title   string `json:"spa_title"`
}

//type BastionPayRes struct {
//	Content       string   `json:"content"`
//}

func (task *Tasker) Init() {
	rate, err := limiter.NewRateFromFormatted("1-M")
	if err != nil {
		ZapLog().Error("tasker limit format err", zap.Error(err))
		return
	}
	store := memory.NewStore()
	task.mSendNotifyLimit = limiter.New(store, rate)
}

func (this *Tasker) Start() {
	this.mSinglechatChan = make(chan *api.SingleChat, 100000)
	for i := 0; i < 10; i++ {
		go this.sChatWorker()
	}

	go func() {
		this.genSign()
		time.Sleep(time.Second * 48 * 60 * 60)
	}()
}

func (this *Tasker) genSign() {
	appidInt, err := strconv.Atoi(config.GConfig.Tencent.Sdkappid)
	if err != nil {
		ZapLog().Info("appid string to int err", zap.Any("appid string to int err", err))
		return
	}

	sign, err := GenSig(appidInt, config.GConfig.Tencent.Key, config.GConfig.Tencent.Identifier, 48*60*60)
	if err != nil {
		ZapLog().Error("get sign err ", zap.Any("get sign err ", err))
		return
	}
	Sign = sign
}

func (this *Tasker) sChatWorker() {
	//查cache
	for {
		sigleChat := <-this.mSinglechatChan
		ZapLog().Info("get chan msg..", zap.Any("msg", sigleChat))
		ZapLog().Info("toaccount", zap.String("toaccount", sigleChat.ToAccount))
		toAccount, err := db.GCache.GetAccountCache(sigleChat.ToAccount)

		if err != nil {
			ZapLog().Info("get cache err", zap.Any("get cache err", err))
			//return
			//如果查的时候有错误
		}
		//fmt.Println("work******here", toAccount == nil, toAccount)
		if toAccount == nil {
			randNum := GenRandomString(32)
			tencentUrl := fmt.Sprintf(tencentUrl,
				Sign,
				config.GConfig.Tencent.Identifier,
				config.GConfig.Tencent.Sdkappid,
				randNum)

			account := make([]string, 0)
			account = append(account, sigleChat.ToAccount)

			reqBody, _ := json.Marshal(map[string]interface{}{
				"To_Account": account,
			})

			//腾讯查这个用户的状态
			res, err := base.HttpSend(tencentUrl, bytes.NewBuffer(reqBody), "POST", nil)
			ZapLog().Info("tencent res", zap.Any("res", string(res)))
			if err != nil {
				ZapLog().Error("get tencent statu err", zap.Error(err))
				continue
				//往腾讯拿用户状态出错了呢 return？
			}

			tencentRes := new(tencent.Response)
			if err = json.Unmarshal(res, tencentRes); err != nil {
				ZapLog().Error(" tencent response unmarshal err", zap.Error(err))
				continue
			}

			ZapLog().Info("tencentRes", zap.Any("res", tencentRes))

			//更新状态  在cache中存入用户的状态
			if tencentRes.ActionStatus != "OK" {
				ZapLog().Error(" tencent response unmarshal err", zap.Error(err))
				continue
			}

			if tencentRes.ActionStatus == "OK" {
				if len(tencentRes.QueryResult) <= 0 {
					ZapLog().Error(" queryResult  err", zap.Error(err))
				}
				db.GCache.SetAccountCache(tencentRes.QueryResult[0].ToAccount, tencentRes.QueryResult[0].State)
			}

			//判断 是否在线，在线则 给消息内容，给红点
			//if tencentRes.QueryResult[0].State == "Online" {
			//	flag, err := this.isLimit(tencentRes.QueryResult[0].ToAccount)
			//	ZapLog().Info("online is limit flag ",zap.Bool("onflag",flag))
			//	if err != nil {
			//		ZapLog().Error( "get account isLimit  err", zap.Error(err))
			//		return
			//	}
			//	fmt.Println("*****is limit flag ****** ",flag)
			//	if flag == false {
			//		_ ,err := this.SendTextNotify(tencentRes.QueryResult[0].ToAccount, sigleChat.MsgBody[0].MsgContent.Text, true)
			//		fmt.Println("*****online SendTextNotify  succ****** ")
			//		if err != nil {
			//			ZapLog().Error( "send message to pastionpay err", zap.Error(err))
			//			return
			//		}
			//
			//		ZapLog().Debug("send notify succ")
			//		return
			//	}else {
			//		//被限制，直接返回
			//		return
			//	}
			//	return
			//}

			//给消息内容，给红点   在线 不在线  都推
			if tencentRes.QueryResult[0].State == "Offline" || tencentRes.QueryResult[0].State == "Online" || tencentRes.QueryResult[0].State == "PushOnline" {

				//向bastion 发送个通知, 暂定一分钟推送一次红点消息
				//flag, err := this.isLimit(tencentRes.QueryResult[0].ToAccount)
				//ZapLog().Info("is limit flag ",zap.Bool("offflag",flag))
				//if err != nil {
				//	ZapLog().Error( "get account isLimit  err", zap.Error(err))
				//	return
				//}
				//if flag == false {
				ZapLog().Info("msgType1", zap.Any("msgtype:", sigleChat.MsgBody[0].MsgType))
				bastionRes, err := this.SendTextNotify(sigleChat.FromAccount, tencentRes.QueryResult[0].ToAccount, sigleChat.MsgBody[0].MsgContent.Text, true, sigleChat.MsgBody[0].MsgType)
				ZapLog().Info("Send succ", zap.Any("bastionRes:", string(bastionRes)))
				if err != nil {
					ZapLog().Error("send message to pastionpay err", zap.Error(err))
				}

				ZapLog().Debug("send notify succ")

				//}else {
				//	//被限制，不给推红点消息，直接返回
				//	return
				//}
			}
		} else {
			//flag, err := this.isLimit(sigleChat.ToAccount)
			//ZapLog().Info("online is limit flag ",zap.Bool("onflag",flag))
			//if err != nil {
			//	ZapLog().Error( "get account isLimit  err", zap.Error(err))
			//	return
			//}
			//fmt.Println("*****is limit flag ****** ",flag)
			//if flag == false {
			ZapLog().Info("msgType2", zap.Any("msgtype:", sigleChat.MsgBody[0].MsgType))
			bastionRes, err := this.SendTextNotify(sigleChat.FromAccount, sigleChat.ToAccount, sigleChat.MsgBody[0].MsgContent.Text, true, sigleChat.MsgBody[0].MsgType)
			ZapLog().Info("online SendTextNotify  succ", zap.Any("bastionRes:", string(bastionRes)))
			if err != nil {
				ZapLog().Error("send message to pastionpay err", zap.Error(err))
			}

			ZapLog().Debug("send notify succ")
			//}else {
			//被限制，直接返回
			//return
			//}
		}
	}
}

//向bastion 发送通知
func (this *Tasker) SendTextNotify(fromId, toId, text string, redflag bool, msgtype string) ([]byte, error) {
	if msgtype != "TIMFaceElem" && msgtype != "TIMSoundElem" && msgtype != "TIMImageElem" && msgtype != "TIMFileElem" && msgtype != "TIMTextElem" {
		ZapLog().Info("go here")
		return nil, nil
	}

	userId, err := strconv.Atoi(toId)
	if err != nil {
		ZapLog().Error("account userId string to int err", zap.Error(err))
		return nil, errors.New("account userId string to int err")
	}

	//msg := Message{}
	msg := new(Message)
	msg.Id = 1
	msg.UserId = userId
	msg.CreateTime = int32(time.Now().Unix())
	msg.HistoryFlag = false     //是否是历史消息
	msg.MsgType = 5             //消息类型
	msg.RedFlag = redflag       //红点消息
	msg.SourceId = "1"          //源消息id
	msg.SubType = "single-chat" //子消息类型

	//text = TransString(text)
	msg.Content.Content = text
	msg.Content.SenderId = fromId
	msg.Content.ReceverId = toId

	msg.Title.En_US_title = "$REMARK_OR_NICK_OR_ID$:" + TransString(text)
	msg.Title.Zh_CN_title = "$REMARK_OR_NICK_OR_ID$:" + TransString(text)
	msg.Title.Zh_TW_title = "$REMARK_OR_NICK_OR_ID$:" + TransString(text)
	msg.Title.Jpa_title = "$REMARK_OR_NICK_OR_ID$:" + TransString(text)
	msg.Title.Spa_title = "$REMARK_OR_NICK_OR_ID$:" + TransString(text) //TIMImageElem

	//if msgtype != "TIMFaceElem" && msgtype != "TIMSoundElem" && msgtype != "TIMImageElem" && msgtype != "TIMFileElem" {
	//	ZapLog().Info("go here")
	//	return nil, nil
	//}

	if msgtype == "TIMFaceElem" {
		msg.Content.Content = "[表情]"
		msg.Title.En_US_title = "$REMARK_OR_NICK_OR_ID$:sent a message to you"
		msg.Title.Zh_CN_title = "$REMARK_OR_NICK_OR_ID$:给你发送一条消息"
		msg.Title.Zh_TW_title = "$REMARK_OR_NICK_OR_ID$:給妳發送壹條消息"
		msg.Title.Jpa_title = "$REMARK_OR_NICK_OR_ID$:があなたにメッセージを送ります"
		msg.Title.Spa_title = "$REMARK_OR_NICK_OR_ID$:sent a message to you"
	}

	if msgtype == "TIMSoundElem" {
		msg.Content.Content = "[语音]"
		msg.Title.En_US_title = "$REMARK_OR_NICK_OR_ID$:sent you a voice"
		msg.Title.Zh_CN_title = "$REMARK_OR_NICK_OR_ID$:给你发送一条语音"
		msg.Title.Zh_TW_title = "$REMARK_OR_NICK_OR_ID$:給妳發送一條語音"
		msg.Title.Jpa_title = "$REMARK_OR_NICK_OR_ID$:があなたへ音声メッセージを送りました"
		msg.Title.Spa_title = "$REMARK_OR_NICK_OR_ID$:sent you a voice"
	}

	if msgtype == "TIMImageElem" {
		msg.Content.Content = "[图片]"
		msg.Title.En_US_title = "$REMARK_OR_NICK_OR_ID$:sent you a picture"
		msg.Title.Zh_CN_title = "$REMARK_OR_NICK_OR_ID$:给你发送一张图片"
		msg.Title.Zh_TW_title = "$REMARK_OR_NICK_OR_ID$:給你發送一張圖片"
		msg.Title.Jpa_title = "$REMARK_OR_NICK_OR_ID$:さんはあなたに写真を送ってあげます"
		msg.Title.Spa_title = "$REMARK_OR_NICK_OR_ID$:sent you a picture"
	}

	if msgtype == "TIMFileElem" {
		msg.Content.Content = "[文件]"
		msg.Title.En_US_title = "$REMARK_OR_NICK_OR_ID$ sent a file to you"
		msg.Title.Zh_CN_title = "$REMARK_OR_NICK_OR_ID$:给你发送一个文件"
		msg.Title.Zh_TW_title = "$REMARK_OR_NICK_OR_ID$:給妳發送壹個檔案"
		msg.Title.Jpa_title = "$REMARK_OR_NICK_OR_ID$:があなたにメッセージを送ります"
		msg.Title.Spa_title = "$REMARK_OR_NICK_OR_ID$:sent a file to you"
	}

	msgByte, err := json.Marshal(msg)

	v := make(url.Values)
	v.Set("message", string(msgByte))
	str := v.Encode()

	stringUrl := bastionPayUrl + str

	ZapLog().Info("send url", zap.Any("url", config.GConfig.BastionpayUrl.Bastionurl+stringUrl))
	bastionPayRes, err := base.HttpSend(config.GConfig.BastionpayUrl.Bastionurl+stringUrl, nil, "GET", nil)
	if err != nil {
		ZapLog().Error("send message to pastionpay err", zap.Error(err))
		return nil, err
	}
	return bastionPayRes, nil
}

func (this *Tasker) SendSChat(msg *api.SingleChat) {
	ZapLog().Info("msg to chan...", zap.Any("msg", msg))
	this.mSinglechatChan <- msg
}

//生成随机字符串 32位
func GenRandomString(l int) string {
	str := "123456789"
	bt := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bt[r.Intn(len(bt))])
	}
	return string(result)
}

func (task *Tasker) isLimit(key string) (bool, error) {
	ctx, err := task.mSendNotifyLimit.Get(nil, key)
	if err != nil {
		return true, err
	}
	if ctx.Reached {
		return true, nil
	}
	return false, nil
}

func TransString(s string) string {
	a := s
	if !strings.Contains(a, "[") || !strings.Contains(a, "]") { //不包含【 或者不包含 】 不处理
		return a
	}
	if !strings.Contains(a, "[") && !strings.Contains(a, "]") { //不包含【 并且 不包含 】 不处理
		return a
	}
	if strings.Contains(a, "[") && !strings.Contains(a, "]") { //包含【  并且不包含 】 不处理
		return a
	}
	if !strings.Contains(a, "[") && strings.Contains(a, "]") { //不包含【  并且 包含 】 不处理
		return a
	}
	for i := 0; i < 30; i++ {
		if strings.Contains(a, "[") && strings.Contains(a, "]") {
			index1 := strings.Index(a, "[")
			index2 := strings.Index(a, "]")
			a = a[:index1] + " " + a[index2+1:]
			//fmt.Println("a:",a)
		} else {
			break
		}
	}
	return a
}
