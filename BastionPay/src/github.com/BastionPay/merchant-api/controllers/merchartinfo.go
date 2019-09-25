package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/merchant-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	//"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
)

//商户的后台配置
type MerchantConfig struct {
	Controllers
}

func (this *MerchantConfig) Add(ctx iris.Context) {
	param := new(models.MerchantAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	err = param.Add()
	if err != nil {
		ZapLog().Error( "Add param to mysql err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
		return
	}

	this.Response(ctx, "add succ")
}

func (this *MerchantConfig) Get(ctx iris.Context) {
	param := new(models.MerchantGet)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		//ZapLog().Error( "param err", zap.Error(err))
		return
	}

	data, err := param.Get()
	//fmt.Println("data ** len",data.List)
	if err != nil {
		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
		return
	}

	this.Response(ctx, data)
}

func (this *MerchantConfig) Update(ctx iris.Context) {
	param := new(models.MerchantUpdate)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	err = param.Update()
	if  err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		return
	}
	this.Response(ctx, "update succ")
}


func (this *MerchantConfig) Delete(ctx iris.Context) {
	param := new(models.MerchantDel)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		ZapLog().Error( "param err", zap.Error(err))
		return
	}

	err = param.Delete()
	if  err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
		return
	}
	this.Response(ctx, "del succ")
}


func (this *MerchantConfig) List(ctx iris.Context) {
	param := new(models.MerchantList)

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