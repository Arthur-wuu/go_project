package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/config"
	"github.com/satori/go.uuid"
	"fmt"
	//"github.com/shopspring/decimal"
	"sync"
	"time"
)


var GLoginTasker LoginTasker

type  LoginTasker struct{
	Login
	//Transfer
	Refresh
	//GetUid

	RefreshToken  *string
	Token         *string
	sync.Mutex
}

type  RedResponse struct{
	Assets    string      `json:"assets,omitempty"`
	Amount    float64     `json:"amount,omitempty"`
}

func (this *LoginTasker) Start() {
	//登录bastionPay的用户登录
	token, refreshToken, err := this.Login.LoginBastionPay(config.GConfig.Login.Phone,config.GConfig.Login.Pwd)
	if  err != nil {
		ZapLog().Sugar().Errorf("get token err [%v]", err)
		return
	}
	//存一下登录返回的token和refreshToken
	this.Token = &token
	this.RefreshToken = &refreshToken
	fmt.Println("[** lgoin token & refresh token** ]", *this.Token, *this.RefreshToken)





	//起刷新token的任务, 遇到token问题，重新登录
	go func(){
		for {
			//睡再去刷新token
			time.Sleep(time.Second * time.Duration( 60*60*8 ))
			this.Lock()
			t, rt, errCode, err := this.Refresh.RefreshToken(*this.Token, *this.RefreshToken) //token

			this.Token = &t
			this.RefreshToken = &rt
			fmt.Println("[** token & refresh token** ]", *this.Token, *this.RefreshToken)


			if errCode == 1014 || errCode == 1001  {
				t2, rt2, err2 := this.Login.LoginBastionPay(config.GConfig.Login.Phone,config.GConfig.Login.Pwd)

				this.Token = &t2
				this.RefreshToken = &rt2
				fmt.Println("[**relogin token & refresh token** ]", *this.Token, *this.RefreshToken)
				//this.Unlock()
				if  err2 != nil {
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



//传切片到裂变服务
type RobberSearch struct{
	ActivityUUId  string    `valid:"optional" json:"activity_uuid,omitempty"`
	CountryCode   string    `valid:"optional" json:"country_code,omitempty"`
	Phone         string    `valid:"optional" json:"phone,omitempty"`
	SrcUid        int64     `valid:"optional" json:"src_uid,omitempty"`
	TransferFlag  int       `valid:"optional" json:"transfer_flag,omitempty"`
	Valid         int       `valid:"optional" json:"valid,omitempty"`
}
//返回这个
type Robber struct{
	Code      int    `valid:"optional" json:"code,omitempty"`
	Message   string `valid:"optional" json:"message,omitempty"`
	Data      []RobberData
}


type RobberData struct{
	Id         		int64       `json:"id,omitempty"`
	ActivityUUId 	string    `json:"activity_uuid,omitempty" `
	RedId      		int       `json:"red_id,omitempty"  `
	AppId      		string    `json:"app_id,omitempty" `
	UserId     		int64     `json:"user_id,omitempty"`
	CountryCode 	string    `json:"country_code,omitempty"`
	Phone      		string    `json:"phone,omitempty" `
	Symbol    	 	string    `json:"symbol,omitempty" `
	Coin     		string    `json:"coin,omitempty"`
	SrcUrl     		string    `json:"src_url,omitempty"`
	SrcUid     		int64     `json:"src_uid,omitempty"`
	//ExpireAt   	*int64     `json:"expire_at,omitempty"  gorm:"column:expire_at;type:bigint(20)"`
	TransferFlag 	int       `json:"transfer_flag,omitempty"`
	SponsorAccount  string    `json:"sponsor_account,omitempty" `
	OffAt           int64     `json:"off_at,omitempty"`
	//Table
}

type Response struct{
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`

}


func  Parse(phone, countryCode string) RobberSearch {

	flag := new(int)
	*flag = 0
	countryCode = "0"+countryCode[1:]

	return RobberSearch{
		CountryCode : countryCode,
		Phone: phone,
		TransferFlag: *flag,
	}
}




func  GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	strings :=fmt.Sprintf("%s", ud)

	return strings
}


