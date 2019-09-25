package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	AccountMap struct {
		Controllers
	}
)

func (this *AccountMap) AddForBack(ctx iris.Context) {
	param := new(api.BkAccountMapAdd)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//数据库添加
	data, err := new(models.AccountMap).BkParseAdd(param)

	if err != nil {
		ZapLog().Error("account map add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *AccountMap) UpdateForBack(ctx iris.Context) {
	param := new(api.BkAccountMapUpdate)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	atm, err := new(models.AccountMap).BkParseUpdate(param)

	if err != nil {
		ZapLog().Error("acocunt map record update err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, atm)
}

func (this *AccountMap) ListForBack(ctx iris.Context) {
	param := new(api.BkAccountMapList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.AccountMap).BkParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("acocunt map record list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if vip == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *AccountMap) GetAccount(ctx iris.Context) {
	param := new(api.BkAccountMap)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if param.Id == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "id is nil")
		ZapLog().Error("id is nil", zap.Error(err))
		return
	}

	act, err := new(models.AccountMap).GetById(*param.Id)

	if err != nil {
		ZapLog().Error("account map record get err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, act)
}
