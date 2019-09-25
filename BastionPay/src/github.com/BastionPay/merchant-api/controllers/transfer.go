package controllers

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/comsumer"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/sdk"
	"fmt"
	"strconv"
	"time"

	"github.com/kataras/iris"
	"go.uber.org/zap"
)


type (
	Transfer struct {
		Controllers
	}
)

const(
	//Red_Request_Url = "https://test-activity.bastionpay.io"
)

type NotifyRequest struct{
	//Id      		*int       `json:"id,omitempty"  `
	//ActivityUuid    *string    `json:"activity_uuid,omitempty"  `
	//RedId      		*string       `json:"red_uuid,omitempty"  `
	//AppId     		*string    `json:"app_id,omitempty"`
	UserId     		*string    `json:"user_id,omitempty"`
	//CountryCode 	*string    `json:"country_code,omitempty"`
	//Phone      		*string    `json:"phone,omitempty" `
	Symbol    	 	*string    `json:"symbol,omitempty" `
	Coin     		*string   `json:"coin,omitempty"`
	MerchantId       *string   `json:"merchant_id,omitempty"`
	//SponsorAccount  *string    `json:"sponsor_account,omitempty" `
	//ApiKey          *string    `json:"api_key,omitempty" `
	//OffAt           *int64     `json:"off_at,omitempty" `
	//Lang            *string     `json:"language,omitempty" `
	//TransferFlag 	*int       `json:"transfer_flag,omitempty"`
}


type ResponseBack struct{
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

//70123b81-9754-4fd4-b2cd-ecac5d3947cd

func (this *Transfer) TransferATM(ctx iris.Context) {
	params := new(NotifyRequest)
	err :=ctx.ReadJSON(params)
	if err != nil {
		this.ExceptionSerive(ctx, 100901, "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}
	ZapLog().Info("params",zap.Any("params:" ,*params.Coin ))

	coin2cash_fee, err := strconv.ParseFloat(config.GConfig.Fee.Coin2cash, 64)
	if err != nil {
		ZapLog().Error( "string to float err", zap.Error(err))
		return
	}

	c := *params.Coin
	coin, err := strconv.ParseFloat( c, 64)
	if err != nil {
		ZapLog().Error( "string to float err", zap.Error(err))
		return
	}

	amounts := coin * (1 - coin2cash_fee)
	amountDec := Decimal(amounts)
	stringAmount := fmt.Sprintf("%v",amountDec)

	//err = ctx.ReadJSON(params)
	//if err != nil {
	//	ZapLog().Error("param err", zap.Error(err))
	//	return
	//}
	fmt.Println("[** notify params coming : **] :", *params.UserId, *params.Symbol, *params.Coin)

	//充值
	pwd := config.GConfig.Login.ZfPwd
	requestNo := time.Now().Unix()
	requestNoStr := strconv.FormatInt(requestNo, 10)
	//userId := strconv.FormatInt(uid, 10)
	//coin := fmt.Sprintf("%s", params.Coin)
	fmt.Println("[** coin num : **] :", *params.Coin)

	status, err := comsumer.GTransfer.TransferCoin(*params.Symbol, pwd, *comsumer.GLoginTasker.Token, *params.UserId, requestNoStr, stringAmount)
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
	return

}

	//sliceId := make([]int,0)
	//sliceId = append(sliceId, params.Id)


// atm通过open-api做转账    数据库里需要配个商户1 为atm的商户，
func (this *Transfer) TransferAtmOpenApi(ctx iris.Context) {
	params := new(NotifyRequest)
	err :=ctx.ReadJSON(params)
	if err != nil {
		this.ExceptionSerive(ctx, 100901, "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}
	ZapLog().Info("params",zap.Any("params:" ,*params.Coin ))

	coin2cash_fee, err := strconv.ParseFloat(config.GConfig.Fee.Coin2cash, 64)
	if err != nil {
		ZapLog().Error( "string to float err", zap.Error(err))
		return
	}

	c := *params.Coin
	coin, err := strconv.ParseFloat( c, 64)
	if err != nil {
		ZapLog().Error( "string to float err", zap.Error(err))
		return
	}

	amounts := coin * (1 - coin2cash_fee)
	amountDec := Decimal(amounts)
	stringAmount := fmt.Sprintf("%v",amountDec)

	//err = ctx.ReadJSON(params)
	//if err != nil {
	//	ZapLog().Error("param err", zap.Error(err))
	//	return
	//}
	fmt.Println("[** notify params coming : **] :", *params.UserId, *params.Symbol, *params.Coin)

	//充值
	//pwd := config.GConfig.Login.ZfPwd
	//requestNo := time.Now().Unix()
	//requestNoStr := strconv.FormatInt(requestNo, 10)
	////userId := strconv.FormatInt(uid, 10)
	////coin := fmt.Sprintf("%s", params.Coin)
	//fmt.Println("[** coin num : **] :", *params.Coin)

	//转账
	strStatus, err := sdk.GPaySdk.Transfer(stringAmount, *params.Symbol, *params.MerchantId, *params.UserId, "nil","http://nil.com")

	//status, err := comsumer.GTransfer.TransferCoin(*params.Symbol, pwd, *comsumer.GLoginTasker.Token, *params.UserId, requestNoStr, stringAmount)
	if err != nil {
		ZapLog().Sugar().Errorf("transfer coin err[%v]", err)
		this.ExceptionSerive(ctx, 1009, err.Error())
		return
	}
	if strStatus == "succ" {
		ZapLog().Sugar().Info("transfer succ")
		this.Response(ctx, " ")
		return
	}
	if strStatus == "fail" {
		ZapLog().Sugar().Info("transfer fail")
		this.ExceptionSerive(ctx, 1008, "transfer fail")
		return
	}
	if strStatus == "unpay" {
		ZapLog().Sugar().Info("transfer running...")
		this.ExceptionSerive(ctx, 1007, "transfer running...")
		return
	}
	return

}


