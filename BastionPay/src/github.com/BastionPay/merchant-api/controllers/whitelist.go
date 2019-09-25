package controllers

import (
	"BastionPay/bas-api/apibackend"
	"BastionPay/merchant-api/models"
	"github.com/kataras/iris"
)

type WhiteList struct {
	Controllers
}

//func (this *WhiteList) Add(ctx iris.Context) {
//	param := new(models.BkConfigAdd)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	err = param.Add()
//	if err != nil {
//		ZapLog().Error( "Add param to mysql err", zap.Error(err))
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(),  err.Error())
//		return
//	}
//
//	this.Response(ctx, "add succ")
//}

func (this *WhiteList) Get(ctx iris.Context) {
	param := new(models.WhiteListGet)

	err := ctx.ReadJSON( param)
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
	arr := make([]string,0)
	arr = append(arr, *data)

	this.Response(ctx, arr)
}
//
//
//func (this *BkConfig) GetCoinList (ctx iris.Context) {
//	param := new(models.BkConfigGet)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		//ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	data, err := param.GetCoinList()
//	//fmt.Println("data ** len",data.List)
//	if err != nil {
//		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
//		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
//		return
//	}
//
//	this.Response(ctx, data)
//}
//
//
//
//func (this *BkConfig) Update(ctx iris.Context) {
//	param := new(models.BkConfigUpdate)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	err = param.Update()
//	if  err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		return
//	}
//	this.Response(ctx, "update succ")
//}
//
//
//func (this *BkConfig) Delete(ctx iris.Context) {
//	param := new(models.BkConfigDel)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	err = param.Delete()
//	if  err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		return
//	}
//	this.Response(ctx, "del succ")
//}
//
//
//func (this *BkConfig) List(ctx iris.Context) {
//	param := new(models.BkConfigList)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err:"+err.Error())
//		//ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	data, err := param.List()
//	if err != nil {
//		//l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
//		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
//		return
//	}
//
//	this.Response(ctx, data)
//}