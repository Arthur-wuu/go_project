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
	Activity struct {
		Controllers
	}
)

func (this *Activity) Add(ctx iris.Context) {
	param := new(api.ActivityAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Activity).ParseAdd(param).Add()
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, data)
}

func (this *Activity) Update(ctx iris.Context) {
	param := new(api.Activity)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Activity).Parse(param).Update()
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Activity) ListForBack(ctx iris.Context) {
	param := new(api.ActivityList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Activity).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Activity) ListForFront(ctx iris.Context) {
	param := new(api.ActivityList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Activity).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("Get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}
