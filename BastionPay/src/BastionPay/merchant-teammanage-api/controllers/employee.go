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
	Employee struct {
		Controllers
	}
)

func (this *Employee) Add(ctx iris.Context) {
	param := new(api.EmployeeAdd)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Employee).ParseAdd(param).Add()
	if err != nil {
		ZapLog().Error("Activity Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *Employee) Get(ctx iris.Context) {
	param := new(api.Employee)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).Parse(param).Get()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Employee) Del(ctx iris.Context) {
	param := new(api.Employee)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = new(models.Employee).Parse(param).Del()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, nil)
}

func (this *Employee) GetFront(ctx iris.Context) {
	param := new(api.FtEmployeeGet)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).FtParseGet(param).Get()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Employee) GetsFront(ctx iris.Context) {
	param := new(api.FtEmployeeGets)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).FtParseGets(param).Gets()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Employee) Update(ctx iris.Context) {
	param := new(api.Employee)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).Parse(param).Update()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Employee) List(ctx iris.Context) {
	param := new(api.EmployeeList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).ParseList(param).ListWithConds(param.Page, param.Size, nil)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Employee) ListFront(ctx iris.Context) {
	param := new(api.FtEmployeeList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Employee).FtParseList(param).ListWithConds(param.Page, param.Size, nil)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}
