package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-userhelp/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	//"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
)

type UserHelp struct {
	Controllers
}

func (this *UserHelp) Add(ctx iris.Context) {
	param := new(models.UserHelpAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	result, err := param.Add()
	if err != nil {
		ZapLog().Error("Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, result)
}

func (this *UserHelp) List(ctx iris.Context) {
	param := new(models.UserHelpList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		//ZapLog().Error( "param err", zap.Error(err))
		return
	}

	data, err := param.List()
	//fmt.Println("data ** len",data.List)
	if err != nil {
		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, data)
}

func (this *UserHelp) Update(ctx iris.Context) {
	param := new(models.UserHelpUpdate)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := param.Update()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		return
	}
	this.Response(ctx, data)
}
