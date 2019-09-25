package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-userhelp/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	//"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
)

type AppVersion struct {
	Controllers
}

func (this *AppVersion) Add(ctx iris.Context) {
	param := new(models.AppVersionSet)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//num ,err := param.RowsAffectNumUpdate()
	//
	//fmt.Println("param.RowsAffectNumUpdate()", num)
	//if num == 0 {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  "info has exist in db ")
	//	return
	//}
	//if err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  "info has exist in db ")
	//	return
	//}
	label, err := new(models.AppWhiteLabel).GetBy(param.LabelId)
	if err != nil {
		ZapLog().Error("Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}
	if label == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), err.Error())
		return
	}
	param.Name = label.Name

	result, err := param.Set()
	if err != nil {
		ZapLog().Error("Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, result)
}

func (this *AppVersion) GetForFront(ctx iris.Context) {
	param := new(models.AppVersionGet)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//label, err := new(models.AppWhiteLabel).GetBy(param.Name)
	//if err != nil {
	//	ZapLog().Error( "Add param to mysql err", zap.Error(err))
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
	//	return
	//}
	//if label == nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(),  err.Error())
	//	return
	//}

	result, err := param.GetForFront(param.Name, param.SysType)
	if err != nil {
		ZapLog().Error("Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}
	if result == nil {
		this.Response(ctx, "")
		return
	}

	this.Response(ctx, result)
}

func (this *AppVersion) List(ctx iris.Context) {
	param := new(models.AppVersionList)

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

func (this *AppVersion) Update(ctx iris.Context) {
	param := new(models.AppVersionUpdate)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	_, err = param.Update()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		return
	}
	this.Response(ctx, "update succ")
}
