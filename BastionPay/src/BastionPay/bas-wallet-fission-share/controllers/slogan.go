package controllers

import (
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Slogan struct {
		Controllers
	}
)

func (this *Slogan) Add(ctx iris.Context) {
	param := new(api.SloganAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Slogan).ParseAdd(param).Add()
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, data)
}

func (this *Slogan) Update(ctx iris.Context) {
	param := new(api.Slogan)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Slogan).Parse(param).Update()
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, vip)
}

func (this *Slogan) ListForBack(ctx iris.Context) {
	param := new(api.SloganList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Slogan).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, vip)
}

func (this *Slogan) ListForFront(ctx iris.Context) {
	param := new(api.SloganList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Slogan).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, vip)
}

func (this *Slogan) GetAll(ctx iris.Context) {
	param := new(api.SloganGets)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Slogan).GetAll(*param.ActivityId)
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, vip)
}
