package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type GaController struct {
	redis     *common.Redis
	db        *gorm.DB
	config    *config.Config
	userModel *models.UserModel
}

func NewGaController(redis *common.Redis, db *gorm.DB, config *config.Config) *GaController {
	return &GaController{
		redis: redis,
		db:    db, config: config,
		userModel: models.NewUserModel(db),
	}
}

func (g *GaController) Generate(ctx iris.Context) {
	var (
		userId   uint
		username string
		err      error
	)
	userId = common.GetUserIdFromCtx(ctx)

	// 获取用户信息
	user, err := g.userModel.GetUserById(userId)
	if err != nil || user.Ga != "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	if user.RegistrationType == "email" {
		username = user.Email
	} else {
		username = user.Phone
	}

	// 生成GA
	ga := common.NewGA()
	err = ga.Generate(username)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ga Generate err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_SERVICE_TEMPORARILY_UNAVAILABLE.Code()))
		return
	}

	// 放入redis,等待验证
	id, err := common.NewVerification(g.redis, "bind_ga", common.VerificationTypeGa).GenerateGA(userId, ga.Secret)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("NewVerification err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_GA_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Id     string `json:"id"`
		Secret string `json:"secret"`
		Image  string `json:"image"`
	}{id, ga.Secret, ga.Image}))
}

func (g *GaController) Bind(ctx iris.Context) {
	var (
		userId uint
		params = struct {
			Id    string `json:"id"`
			Value string `json:"value"`
		}{}
		err error
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Id == "", "id"},
		{params.Value == "", "value"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 获取用户信息
	user, err := g.userModel.GetUserById(userId)
	if err != nil || user.Ga != "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	verification := common.NewVerification(g.redis, "bind_ga", common.VerificationTypeGa)

	bol, err := verification.Verify(params.Id, userId, params.Value, "")
	if err != nil || !bol {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
			//			glog.Error(err.Error())
		}

		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_GA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code()))
		return
	}

	// 绑定
	err = g.userModel.BindGa(userId, verification.Value)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("BindGa err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}

func (g *GaController) UnBind(ctx iris.Context) {
	var (
		userId uint
		err    error
		params = struct {
			Value      string `json:"value"`
			EmailToken string `json:"email_token"`
			SmsToken   string `json:"sms_token"`
		}{}
	)
	userId = common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	// 获取用户信息
	user, err := g.userModel.GetUserById(userId)
	if err != nil || user.Ga == "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code()))
		return
	}

	// 验证token参数
	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Value == "", "value"},
		{user.Email != "" && params.EmailToken == "", "email_token"},
		{user.Phone != "" && params.SmsToken == "", "sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 验证GA
	ga := common.NewGA()
	bol, err := ga.Verify(user.Ga, params.Value)
	if err != nil || !bol {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_GA_TOKEN_AUTHENTICATION", apibackend.BASERR_INCORRECT_GA_PWD.Code()))
		return
	}

	// 验证邮箱
	if user.Email != "" {
		tokenPass, err := common.NewVerification(g.redis, "unbind_ga", common.VerificationTypeEmail).
			Check(params.EmailToken, userId, user.Email)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//			glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 验证手机
	if user.Phone != "" {
		tokenPass, err := common.NewVerification(g.redis, "unbind_ga", common.VerificationTypeSms).
			Check(params.SmsToken, userId, user.Phone)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 解绑
	err = g.userModel.UnBindGa(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("UnBindGa err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}
