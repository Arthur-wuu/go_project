package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/api"
	"BastionPay/merchant-api/baspay"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Assets struct {
		Controllers
	}
)

//查询可用的币种信息
func (this *Assets) AvaliAssets(ctx iris.Context) {

	param := new(api.AvAssets)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	res, err := new(baspay.Assets).Parse(param).Send()
	if err != nil {
		ZapLog().Error("select trade info err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, res)
}
