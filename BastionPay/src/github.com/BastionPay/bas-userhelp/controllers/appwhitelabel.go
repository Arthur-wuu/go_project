package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-userhelp/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	//"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
)

type AppWhiteLabel struct {
	Controllers
}

func (this *AppWhiteLabel) Add(ctx iris.Context) {
	param := new(models.AppWhiteLabelAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	result, err := param.Add()
	if err != nil {
		ZapLog().Error( "Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
		return
	}

	this.Response(ctx, result)
}

func (this *AppWhiteLabel) List(ctx iris.Context) {
	param := new(models.AppWhiteLabelList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		//ZapLog().Error( "param err", zap.Error(err))
		return
	}

	data, err := param.List()
	if err != nil {
		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, data)
}

func (this *AppWhiteLabel) Update(ctx iris.Context) {
	param := new(models.AppWhiteLabel)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	data, err  := param.Update()
	if  err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		return
	}
	this.Response(ctx, data)
}