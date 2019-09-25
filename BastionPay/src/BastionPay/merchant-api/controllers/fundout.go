package controllers

import (
	"BastionPay/merchant-api/models"
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	FundOut struct {
		Controllers
	}
)

func (this *FundOut) Create(ctx iris.Context) {
	//创建订单
	//发送baspay
	//等回调通知吧
	param := new(api.FundOut)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	param.MerchantTradeNo = common.GenerateUuid()

	if err := new(models.FundOut).Parse(param).Add(); err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	res, err := new(baspay.FundOut).Parse(param).Send()
	if err != nil {
		ZapLog().Error("baspay Send err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	err = new(models.FundOut).UpdateByTradeNo(*res.MerchantFundoutNo, *res.FundoutNo, models.EUM_FUNDOUT_STATUS_APPLY)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}
	this.Response(ctx, &api.ResFundOut{MerchantFundoutNo: res.MerchantFundoutNo, FundoutNo: res.FundoutNo})
}

func (this *FundOut) List(ctx iris.Context) {
	//创建订单
	//发送baspay
	//等回调通知吧
	param := new(api.FundOutList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	res, err := new(models.FundOut).List(param.Page, param.Size, param.MerchantTradeNo, param.PayeeId, param.TradeNo)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, res)
}
