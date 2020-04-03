package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type Markets []string

type OptionalController struct {
	Controllers
}

func NewOptionalController() *OptionalController {
	return &OptionalController{}
}

func (this *OptionalController) Get(ctx iris.Context) {
	userId := common.GetUserIdFromCtx(ctx)

	markets, err := models.GOptionalModel.Get(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("model get err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "GET_OPTIONAL_FAILED")
		return
	}

	this.Response(ctx, markets)
}

func (this *OptionalController) Update(ctx iris.Context) {
	var (
		params []string
	)

	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	if err = models.GOptionalModel.Delete(userId); err != nil {
		ZapLog().With(zap.Error(err)).Error("model delete err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SET_OPTIONAL_FAILED")
		return
	}

	if _, err = models.GOptionalModel.Create(userId, params); err != nil {
		ZapLog().With(zap.Error(err)).Error("model create err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SET_OPTIONAL_FAILED")
		return
	}

	this.Response(ctx, nil)
}
