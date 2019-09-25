package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Atm struct {
		Controllers
	}
)
//atm 有手续费的 订单
func (this *Atm) CreateAtmQr(ctx iris.Context) {
	param := new(api.QrTrade)

	err :=ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}
	param.MerchantTradeNo = common.GenerateUuid()
	param.ExpireTime = 900

	res, err := new(baspay.QrTrade).QrParse(param).SendDiscountQr()
	if err != nil {
		ZapLog().Error( "baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}
	if res.Code != 0 {
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: "create qr order err"})
		return
	}

	this.Response(ctx, res.Data)
}




