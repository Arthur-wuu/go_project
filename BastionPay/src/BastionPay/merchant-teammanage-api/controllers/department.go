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
	Department struct {
		Controllers
	}
)

func (this *Department) Add(ctx iris.Context) {
	param := new(api.DepartmentAdd)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Department).ParseAdd(param).Add()
	if err != nil {
		ZapLog().Error("Activity Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *Department) Get(ctx iris.Context) {
	param := new(api.Department)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).Parse(param).Get()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Department) Del(ctx iris.Context) {
	param := new(api.Department)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	existFlag, err := new(models.Employee).ExistByDepartmentId(*param.Id)
	if err != nil {
		ZapLog().Error("ExistByDepartmentId Db err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if existFlag {
		ZapLog().Error("Department exist employee and cannot del err")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "Department exist employee and cannot del err")
		return
	}

	err = new(models.Department).Parse(param).Del()
	if err != nil {
		ZapLog().Error("del Db err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, nil)
}

func (this *Department) GetFront(ctx iris.Context) {
	param := new(api.FtDepartmentGet)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).FtParseGet(param).Get()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Department) GetsFront(ctx iris.Context) {
	param := new(api.FtDepartmentGets)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).GetsByCompany(*param.CompanyId)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Department) Update(ctx iris.Context) {
	param := new(api.Department)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).Parse(param).Update()
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Department) List(ctx iris.Context) {
	param := new(api.DepartmentList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).ParseList(param).ListWithConds(param.Page, param.Size, nil)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *Department) ListFront(ctx iris.Context) {
	param := new(api.FtDepartmentList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.Department).FtParseList(param).ListWithConds(param.Page, param.Size, nil)
	if err != nil {
		ZapLog().Error("Activity GetForFront err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	this.Response(ctx, vip)
}
