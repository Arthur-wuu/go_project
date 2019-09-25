package controllers

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-recharge/comsumer"
	"BastionPay/baspay-recharge/config"
	"fmt"
	//"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"time"

	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	RubbishNotify struct {
		Controllers
	}
)

//var r *rand.Rand
//const(
//	Red_Request_Url = "https://test-activity.bastionpay.io"
//)
//
//type NotifyRequest struct{
//	Id      		*int       `json:"id,omitempty"  `
//	RedId      		*string       `json:"red_uuid,omitempty"  `
//	UserId     		*int64     `json:"user_id,omitempty"`
//	CountryCode 	*string    `json:"country_code,omitempty"`
//	Phone      		*string    `json:"phone,omitempty" `
//	Symbol    	 	*string    `json:"symbol,omitempty" `
//	Coin     		*decimal.Decimal   `json:"coin,omitempty"`
//	SponsorAccount  *string    `json:"sponsor_account,omitempty" `
//	ApiKey          *string    `json:"api_key,omitempty" `
//	OffAt           *int64     `json:"off_at,omitempty" `
//	//TransferFlag 	*int       `json:"transfer_flag,omitempty"`
//}
//
//type NotifyResponse struct{
//	Code    int    `json:"code,omitempty"`
//	Message string `json:"message,omitempty"`
//	Data    string `json:"data,omitempty"`
//
//}
//
//type ResponseBack struct{
//	Code    int    `json:"code,omitempty"`
//	Message string `json:"message,omitempty"`
//	Data    string `json:"data,omitempty"`
//}

//70123b81-9754-4fd4-b2cd-ecac5d3947cd

func (this *RubbishNotify) TransferCallBack(ctx iris.Context) {
	params := new(NotifyRequest)

	err := ctx.ReadJSON(params)
	if err != nil {
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	//fmt.Println("[** notify params coming : **] :",params)

	if params.ApiKey == nil || *params.ApiKey != "70123b81-9754-4fd4-b2cd-ecac5d3947cd" {
		ZapLog().Error("api key nil or err")
		this.ExceptionSerive(ctx, 1000, "api key err")
		return
	}

	if params.UserId == nil {
		fmt.Println("[** start go no uid *****")
		//这里用phone去查询 uid，如果有注册，转账
		//consumer.Task
		code := *params.CountryCode
		if len(code) <= 2 {
			ZapLog().Error("country code err", zap.Error(err))
			return
		}
		code = "+" + code[1:]
		fmt.Println("[** get uid params* phone code]:", *params.Phone, code, *comsumer.GTasker.Token)
		fmt.Println("[** get uid params* Token ]:", *comsumer.GTasker.Token)
		uid, registime, err := comsumer.GTasker.GetUidByPhone(*params.Phone, code, *comsumer.GTasker.Token)
		if err != nil {
			ZapLog().Error("get uid by phone nil")
			this.ExceptionSerive(ctx, 100010, err.Error())
			return
		}
		if err != nil || uid == "" || registime == "" {
			ZapLog().Error("get uid by phone nil", zap.Error(err))
			this.Response(ctx, nil)
			return
		}
		fmt.Println("[** get uid by phone*]:", uid)

		registTime, err := strconv.ParseInt(registime, 10, 64)
		if registTime/1000 > *params.OffAt {
			ZapLog().Error("activite time out, no transfer...", zap.Error(err))
			this.Response(ctx, nil)
			return
		}

		////根据通知的参数，先把trans_flag从0 设置为1，再充值
		//robberUpdate, _ := json.Marshal(map[string]interface{}{
		//	"id": *params.Id ,
		//	"transfer_flag": 1,
		//})
		//
		//res, err :=base.HttpSend(Red_Request_Url+"/v1/fissionshare/robber/set-transferflag",bytes.NewBuffer(robberUpdate),"POST",nil)
		//if err != nil {
		//	ZapLog().Sugar().Errorf("request red err[%v]", err)
		//}
		//fmt.Println("[** set flag res***]",string(res))
		//
		//responseMsg := new(NotifyResponse)
		//json.Unmarshal(res, responseMsg)
		//
		//if responseMsg.Code == apibackend.BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED.Code() {
		//	this.Response(ctx,nil)
		//	return
		//}
		//
		//if responseMsg.Code == 1000100 {
		//	ZapLog().Info( "set flag param nil")
		//}
		//if responseMsg.Code != 0 {
		//	ZapLog().Error( "set flag err")
		//}

		//充值
		pwd := config.GConfig.Login.ZfPwd
		requestNo := time.Now().UnixNano()
		requestNoStr := strconv.FormatInt(requestNo, 10)
		rands := RandString(4)
		requestNoStr = requestNoStr + rands
		//userId := strconv.FormatInt(uid, 10)
		coin := fmt.Sprintf("%s", params.Coin)

		//fmt.Println("**symbol,pwd,token,userid,requestNo,coin**",params.Symbol, pwd, *comsumer.GTasker.Token, uid, requestNo, coin)
		status, err := comsumer.GTasker.Transfer.TransferCoin(*params.Symbol, pwd, *comsumer.GTasker.Token, uid, requestNoStr, coin)
		if err != nil {
			ZapLog().Sugar().Errorf("transfer coin err[%v]", err)
			this.ExceptionSerive(ctx, 1009, err.Error())
			return
		}
		if status == 2 {
			ZapLog().Sugar().Info("transfer succ")
			this.Response(ctx, " ")
			return
		}
		if status == 3 {
			ZapLog().Sugar().Info("transfer fail")
			this.ExceptionSerive(ctx, 1008, "transfer fail")
			return
		}
		if status == 1 {
			ZapLog().Sugar().Info("transfer running...")
			this.ExceptionSerive(ctx, 1007, "transfer running...")
			return
		}
		this.ExceptionSerive(ctx, 1010, err.Error())
		return
	}

	fmt.Println("[** start with uid *****")
	//根据通知的参数，先把trans_flag从0 设置为1，再充值
	//robberUpdate, _ := json.Marshal(map[string]interface{}{
	//	"id": *params.Id ,
	//	"transfer_flag": 1,
	//})
	//
	//res, err :=base.HttpSend(Red_Request_Url+"/v1/fissionshare/robber/set-transferflag",bytes.NewBuffer(robberUpdate),"POST",nil)
	//if err != nil {
	//	ZapLog().Sugar().Errorf("request red err[%v]", err)
	//}
	//fmt.Println("[** set flag with uid res***]", string(res))
	//
	//responseMsg := new(NotifyResponse)
	//json.Unmarshal(res, responseMsg)
	//
	//if responseMsg.Code == apibackend.BASERR_ACTIVITY_FISSIONSHARE_ROBBER_TRANSFERFLAG_NOT_AFFECTED.Code() {
	//	this.Response(ctx,nil)
	//	return
	//}
	//
	//if responseMsg.Code == 1000100 {
	//	ZapLog().Info( "set flag param nil")
	//	this.ExceptionSerive(ctx, 1008, "set-transferflag fail")
	//	return
	//}
	//if responseMsg.Code != 0 {
	//	ZapLog().Error( "set flag err", zap.Error(err))
	//}

	//充值
	pwd := config.GConfig.Login.ZfPwd
	requestNo := time.Now().Unix()
	requestNoStr := strconv.FormatInt(requestNo, 10)
	//userId := strconv.FormatInt(*params.UserId, 10)
	coin := fmt.Sprintf("%s", params.Coin)

	//fmt.Println("**symbol,pwd,token,userid,requestNo,coin**",params.Symbol, pwd, *comsumer.GTasker.Token, userId, requestNo, coin)
	status, err := comsumer.GTasker.Transfer.TransferCoin(*params.Symbol, pwd, *comsumer.GTasker.Token, *params.UserId, requestNoStr, coin)
	if err != nil {
		ZapLog().Sugar().Errorf("transfer coin err[%v]", err)
		this.ExceptionSerive(ctx, 1009, err.Error())
		return
	}
	if status == 2 {
		this.Response(ctx, " ")
		return
	}
	if status == 3 {
		ZapLog().Sugar().Info("transfer fail")
		this.ExceptionSerive(ctx, 1008, "transfer fail")
		return
	}
	if status == 1 {
		ZapLog().Sugar().Info("transfer running...")
		this.ExceptionSerive(ctx, 1007, "transfer running...")
		return
	}
	this.ExceptionSerive(ctx, 1010, err.Error())
	return
}

//sliceId := make([]int,0)
//sliceId = append(sliceId, params.Id)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
//func RandString(len int) string {
//	bytes := make([]byte, len)
//	for i := 0; i < len; i++ {
//		b := r.Intn(26) + 65
//		bytes[i] = byte(b)
//	}
//	return string(bytes)
//}
