package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/models"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type Markets []string

type OptionalController struct {
	db            *gorm.DB
	optionalModel *models.OptionalModel
}

func NewOptionalController(db *gorm.DB) *OptionalController {
	return &OptionalController{
		db:            db,
		optionalModel: models.NewOptionalModel(db),
	}
}

func (o *OptionalController) Get(ctx iris.Context) {
	userId := common.GetUserIdFromCtx(ctx)

	markets, err := o.optionalModel.Get(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("model get err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_OPTIONAL_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, markets))
}

func (o *OptionalController) Update(ctx iris.Context) {
	var (
		params []string
		err    error
	)

	userId := common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if err = o.optionalModel.Delete(userId); err != nil {
		ZapLog().With(zap.Error(err)).Error("model delete err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SET_OPTIONAL_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	if _, err = o.optionalModel.Create(userId, params); err != nil {
		ZapLog().With(zap.Error(err)).Error("model create err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SET_OPTIONAL_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
}
