package comsumer

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-recharge/base"
	"BastionPay/baspay-recharge/config"
	"BastionPay/baspay-recharge/db"
	"BastionPay/baspay-recharge/models/table"
	"bytes"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math/rand"
	"sync"

	//"/BastionPay/baspay-recharge/models/table"
	//"BastionPay/baspay-recharge/table"
	"fmt"
	"strconv"
	"time"
)

const (
	Prefix_New_Id = "Last_Id"
	//Red_Request_Url = "https://test-activity.bastionpay.io"
)

var GTasker Tasker

type Tasker struct {
	Login
	Transfer
	Refresh
	GetUid

	RefreshToken *string
	Token        *string
	sync.Mutex
}

type RedResponse struct {
	Assets string  `json:"assets,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

func (this *Tasker) Start() {
	//先把要登录bastionPay的用户登录了， token拿到，才可以转账
	token, refreshToken, err := this.Login.LoginBastionPay(config.GConfig.Login.Phone, config.GConfig.Login.Pwd)
	if err != nil {
		ZapLog().Sugar().Errorf("get token err [%v]", err)
		return
	}
	//存一下登录返回的token和refreshToken
	this.Token = &token
	this.RefreshToken = &refreshToken
	fmt.Println("[** lgoin token & refresh token** ]", *this.Token, *this.RefreshToken)
	//this.StoreNewIDToRedis(0)

	//起查表的任务
	go func() {
		for {
			this.Lock()
			t := this.Token
			this.Unlock()
			this.run(*t)

		}
	}()

	//起刷新token的任务, 遇到token问题，重新登录
	go func() {
		for {
			//睡四分钟再去刷新token
			time.Sleep(time.Second * time.Duration(60*60*8))
			this.Lock()
			t, rt, errCode, err := this.Refresh.RefreshToken(*this.Token, *this.RefreshToken) //token

			this.Token = &t
			this.RefreshToken = &rt
			fmt.Println("[** token & refresh token** ]", *this.Token, *this.RefreshToken)

			if errCode == 1014 || errCode == 1001 {
				t2, rt2, err2 := this.Login.LoginBastionPay(config.GConfig.Login.Phone, config.GConfig.Login.Pwd)

				this.Token = &t2
				this.RefreshToken = &rt2
				fmt.Println("[**relogin token & refresh token** ]", *this.Token, *this.RefreshToken)
				this.Unlock()
				if err2 != nil {
					ZapLog().Sugar().Errorf("get token err [%v]", err)
					continue
				}
			}
			this.Unlock()

			if err != nil {
				ZapLog().Sugar().Errorf("refresh token err [%v]", err)
				continue
			}
		}
	}()
}

//每隔几秒钟查询一次

func (this *Tasker) run(token string) { // token string

	id, err := this.getNewIdFromRedis()
	fmt.Println("[** get new id from redis** ]:", id)
	if err != nil {
		ZapLog().Sugar().Errorf("get new time from redis err[%v]", err)
		time.Sleep(time.Second * time.Duration(120))
		return
	}

HERE: //查询符合新注册的用户集合
	redInfos, err := this.GetNewRegisterServer(id)
	if err != nil {
		ZapLog().Sugar().Errorf("get data err[%v]", err)
		time.Sleep(time.Second * time.Duration(10))
		goto HERE
	}

	// 从表中拿最新的id，存redis
	idNew, err := this.GetLastId()
	if idNew == nil {
		time.Sleep(time.Second * time.Duration(10))
		goto HERE
	}

	if err != nil {
		ZapLog().Sugar().Errorf("get table new id err[%v]", err)
	}
	this.StoreNewIDToRedis(*idNew.Id)

	//批量去请求 币种金额等信息
	if len(redInfos) == 0 || redInfos == nil {
		time.Sleep(time.Second * time.Duration(10))
		goto HERE
	}

	for _, vr := range redInfos {
		content := new(table.Content)
		json.Unmarshal([]byte(*vr.NotifyContent), content)
		if content.Phone == nil || content.CountryCode == nil || len(*content.Phone) < 3 || len(*content.CountryCode) < 1 {
			continue
		}
		robberSearch := Parse(*content.Phone, *content.CountryCode)
		robberSearchbytes, _ := json.Marshal(robberSearch)

		res, err := base.HttpSendSer(config.GConfig.RedPackUrl.RedPackUrl+"/v1/ft/fissionshare/robber/search", bytes.NewBuffer(robberSearchbytes), "POST", map[string]string{"Api-Key": "1234567890"})
		if err != nil {
			ZapLog().Sugar().Errorf("request red err[%v]", err)
		}

		fmt.Println("[** select money res** ]:", string(res))
		robberInfoBack := new(Robber)
		json.Unmarshal(res, robberInfoBack)

		//对比注册时间和活动结束时间
		//if *content.Phone + *content.CountryCode == robberInfoBack.Data[0].Phone + robberInfoBack.Data[0].CountryCode{
		//	if *content.RegistTime > robberInfoBack.Data[0].OffAt{   //活动时间结束
		//		SetFlag2WithTimeOut(robberInfoBack.Data[0].Id)
		//		continue
		//	}
		//}

		//一个个设置成trans_flag=1  再去bastion去请求 转账
		for _, v := range robberInfoBack.Data {
			v.CountryCode = "+" + v.CountryCode[1:]
			if *content.Phone+*content.CountryCode == v.Phone+v.CountryCode {
				if *content.RegistTime/1000 > v.OffAt || v.OffAt <= time.Now().Unix()-7*24*3600 { //活动时间结束
					SetFlag2WithTimeOut(v.Id)
					continue
				}
			}
			//根据通知的参数，先把trans_flag从0 设置为1，再充值
			sliceId := v.Id
			robberUpdate, _ := json.Marshal(map[string]interface{}{
				"id":            sliceId,
				"transfer_flag": 1,
			})

			res, err := base.HttpSendSer(config.GConfig.RedPackUrl.RedPackUrl+"/v1/ft/fissionshare/robber/set-transferflag", bytes.NewBuffer(robberUpdate), "POST", map[string]string{"Api-Key": "1234567890"})
			if err != nil {
				ZapLog().Error("set transflag err")
				continue
			}

			responseMsg := new(Response)
			json.Unmarshal(res, responseMsg)

			if responseMsg.Code == apibackend.BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED.Code() {
				ZapLog().Error("transfer maybe double")
				continue
			}
			if responseMsg.Code != 0 {
				ZapLog().Error("set flag err", zap.Error(err))
				continue
			}
			code := "+" + v.CountryCode[1:]
			uid, _, err := this.GetUidByPhone(v.Phone, code, *GTasker.Token)
			fmt.Println("**uid**", uid)
			//充值
			pwd := config.GConfig.Login.ZfPwd

			requestNo := time.Now().UnixNano()
			requestNoStr := strconv.FormatInt(requestNo, 10)
			rands := RandString(4)
			requestNoStr = requestNoStr + rands

			//userId := strconv.FormatInt(v.UserId, 10)
			coin := fmt.Sprintf("%s", v.Coin)

			status, err := this.Transfer.TransferCoin(v.Symbol, pwd, *GTasker.Token, uid, requestNoStr, coin)
			fmt.Println("[** trans status** ]:", status)
			if err != nil {
				ZapLog().Sugar().Errorf("transfer err[%v]", err)
			}
			if err != nil {
				ZapLog().Sugar().Errorf("transfer coin err[%v]", err)
			}
			if status == 2 {
				ZapLog().Sugar().Info("transfer succ")
			}
			if status == 3 {
				ZapLog().Sugar().Info("transfer fail")
			}
			if status == 1 {
				ZapLog().Sugar().Info("transfer running...")
			}
		}
	}
	for _, vr := range redInfos {
		content := new(table.Content)
		json.Unmarshal([]byte(*vr.NotifyContent), content)
		if content.Phone == nil || content.CountryCode == nil || len(*content.Phone) < 3 || len(*content.CountryCode) < 1 {
			continue
		}
		robberSearch := Parse(*content.Phone, *content.CountryCode)
		robberSearchbytes, _ := json.Marshal(robberSearch)

		res, err := base.HttpSendSer(config.GConfig.LuckDrawUrl.LuckDrawUrl+"/v1/ft/luckdraw/drawer/search", bytes.NewBuffer(robberSearchbytes), "POST", map[string]string{"Api-Key": "1234567890"})
		if err != nil {
			ZapLog().Sugar().Errorf("request red err[%v]", err)
		}

		fmt.Println("[** select money res draw** ]:", string(res))
		robberInfoBack := new(Robber)
		json.Unmarshal(res, robberInfoBack)

		//一个个设置成trans_flag=1  再去bastion去请求 转账
		for _, v := range robberInfoBack.Data {
			v.CountryCode = "+" + v.CountryCode[1:]
			if *content.Phone+*content.CountryCode == v.Phone+v.CountryCode {
				if *content.RegistTime/1000 > v.OffAt { //活动时间结束
					SetFlag2WithTimeOut(v.Id)
					continue
				}
			}
			//根据通知的参数，先把trans_flag从0 设置为1，再充值
			sliceId := v.Id

			robberUpdate, _ := json.Marshal(map[string]interface{}{
				"id":            sliceId,
				"transfer_flag": 1,
			})

			res, err := base.HttpSendSer(config.GConfig.LuckDrawUrl.LuckDrawUrl+"/v1/ft/luckdraw/drawer/set-transferflag", bytes.NewBuffer(robberUpdate), "POST", map[string]string{"Api-Key": "1234567890"})
			if err != nil {
				ZapLog().Error("set transflag err")
				continue
			}

			responseMsg := new(Response)
			json.Unmarshal(res, responseMsg)

			if responseMsg.Code == apibackend.BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED.Code() {
				ZapLog().Error("transfer maybe double")
				continue
			}
			if responseMsg.Code != 0 {
				ZapLog().Error("set flag err", zap.Error(err))
				continue
			}
			code := "+" + v.CountryCode[1:]
			uid, _, err := this.GetUidByPhone(v.Phone, code, *GTasker.Token)
			fmt.Println("**uid draw**", uid)
			//充值
			pwd := config.GConfig.Login.ZfPwd

			requestNo := time.Now().UnixNano()
			requestNoStr := strconv.FormatInt(requestNo, 10)
			rands := RandString(4)
			requestNoStr = requestNoStr + rands

			//userId := strconv.FormatInt(v.UserId, 10)
			coin := fmt.Sprintf("%s", v.Coin)

			status, err := this.Transfer.TransferCoin(v.Symbol, pwd, *GTasker.Token, uid, requestNoStr, coin)
			fmt.Println("[** trans status draw** ]:", status)
			if err != nil {
				ZapLog().Sugar().Errorf("transfer err[%v]", err)
			}
			if err != nil {
				ZapLog().Sugar().Errorf("transfer coin err[%v]", err)
			}
			if status == 2 {
				ZapLog().Sugar().Info("transfer succ")
			}
			if status == 3 {
				ZapLog().Sugar().Info("transfer fail")
			}
			if status == 1 {
				ZapLog().Sugar().Info("transfer running...")
			}
		}

	}

	time.Sleep(time.Second * time.Duration(10))

}

//传切片到裂变服务
type RobberSearch struct {
	ActivityUUId string `valid:"optional" json:"activity_uuid,omitempty"`
	CountryCode  string `valid:"optional" json:"country_code,omitempty"`
	Phone        string `valid:"optional" json:"phone,omitempty"`
	SrcUid       int64  `valid:"optional" json:"src_uid,omitempty"`
	TransferFlag int    `valid:"optional" json:"transfer_flag,omitempty"`
	Valid        int    `valid:"optional" json:"valid,omitempty"`
}

//返回这个
type Robber struct {
	Code    int    `valid:"optional" json:"code,omitempty"`
	Message string `valid:"optional" json:"message,omitempty"`
	Data    []RobberData
}

type RobberData struct {
	Id           int64           `json:"id,omitempty"`
	ActivityUUId string          `json:"activity_uuid,omitempty" `
	RedId        int             `json:"red_id,omitempty"  `
	AppId        string          `json:"app_id,omitempty" `
	UserId       int64           `json:"user_id,omitempty"`
	CountryCode  string          `json:"country_code,omitempty"`
	Phone        string          `json:"phone,omitempty" `
	Symbol       string          `json:"symbol,omitempty" `
	Coin         decimal.Decimal `json:"coin,omitempty"`
	SrcUrl       string          `json:"src_url,omitempty"`
	SrcUid       int64           `json:"src_uid,omitempty"`
	//ExpireAt   	*int64     `json:"expire_at,omitempty"  gorm:"column:expire_at;type:bigint(20)"`
	TransferFlag   int    `json:"transfer_flag,omitempty"`
	SponsorAccount string `json:"sponsor_account,omitempty" `
	OffAt          int64  `json:"off_at,omitempty"`
	//Table
}

type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

func Parse(phone, countryCode string) RobberSearch {

	flag := new(int)
	*flag = 0
	countryCode = "0" + countryCode[1:]

	return RobberSearch{
		CountryCode:  countryCode,
		Phone:        phone,
		TransferFlag: *flag,
	}
}

func (this *Tasker) GetNewRegister(id int64) ([]table.RederInfoT, error) {
	result := make([]table.RederInfoT, 0)

	err := db.GDbMgr.Get().Table("reder_info").Select("uid, regist_time, country, phone, type, channel").Where("type = ? and id > ?", "1", id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Tasker) GetNewRegisterServer(id int64) ([]table.MqNotify, error) {
	result := make([]table.MqNotify, 0)

	err := db.GDbMgr.Get().Table("MQ_NOTIFY").Select("NOTIFY_CONTENT").Where(" ID > ?", id).Find(&result).Error //NOTIFY_TYPE = ? and
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Tasker) GetLastId() (*table.MqNotify, error) {
	Info := new(table.MqNotify)

	if err := db.GDbMgr.Get().Table("MQ_NOTIFY").Select("ID").Last(Info).Error; err != nil {
		return nil, err
	}
	return Info, nil
}

func (this *Tasker) StoreNewIDToRedis(id int64) error {
	_, err := db.GRedis.Do("SET", Prefix_New_Id, id)
	if err != nil {
		return err
	}
	return nil
}

func (this *Tasker) getNewIdFromRedis() (int64, error) {
	times, err := db.GRedis.Do("GET", Prefix_New_Id)
	if err != nil {
		return 0, err
	}
	timeStr := string(times.([]byte))
	timeInt, _ := strconv.ParseInt(timeStr, 10, 64)
	return timeInt, err
}

func GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	strings := fmt.Sprintf("%s", ud)

	return strings
}

func SetFlag2WithTimeOut(sliceId int64) {
	//if len(sliceId) == 0 {
	//	return
	//}

	robberUpdate, _ := json.Marshal(map[string]interface{}{
		"id":            sliceId,
		"transfer_flag": 2,
	})

	res, err := base.HttpSend(config.GConfig.RedPackUrl.RedPackUrl+"/v1/fissionshare/robber/set-transferflag", bytes.NewBuffer(robberUpdate), "POST", nil)
	if err != nil {
		ZapLog().Sugar().Errorf("request red err[%v]", err)
	}

	responseMsg := new(Response)
	json.Unmarshal(res, responseMsg)

	if responseMsg.Code == 1000100 {
		ZapLog().Info("set flag param nil")
		return
	}

	if responseMsg.Code != 0 {
		ZapLog().Error("set flag err", zap.Error(err))
	}
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
