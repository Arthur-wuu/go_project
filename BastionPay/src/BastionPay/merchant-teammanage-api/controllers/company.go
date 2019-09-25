package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-teammanage-api/api"
	"BastionPay/merchant-teammanage-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Company struct {
		Controllers
	}
)

func (this *Company) Add(ctx iris.Context) {
	param := new(api.CompanyAdd)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Company).ParseAdd(param).Add()
	if err != nil {
		ZapLog().Error("Activity Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *Company) Get(ctx iris.Context) {
	param := new(api.Company)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Company).Parse(param).Get()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Company) Del(ctx iris.Context) {
	param := new(api.Company)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	existFlag, err := new(models.Department).ExistByCompanyId(*param.Id)
	if err != nil {
		ZapLog().Error("ExistByCompanyId Db err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if existFlag {
		ZapLog().Error("Company exist Department and cannot del err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "Company exist Department and cannot del err")
		return
	}

	err = new(models.Company).Parse(param).Del()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, nil)
}

func (this *Company) Update(ctx iris.Context) {
	param := new(api.Company)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Company).Parse(param).Update()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Company) List(ctx iris.Context) {
	param := new(api.CompanyList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Company).ParseList(param).ListWithConds(param.Page, param.Size, nil)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}
