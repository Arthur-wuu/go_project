package controllers

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/models"
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	"BastionPay/pay-user-merchant-api/api"
)

type GaController struct {
	Controllers
}

func NewGaController() *GaController {
	return &GaController{

	}
}

func (this *GaController) Generate(ctx iris.Context) {
	var (
		userId   uint
		username string
		err      error
	)
	userId = common.GetUserIdFromCtx(ctx)

	// 获取用户信息
	user, err := new(models.User).GetById(int64(userId))
	if err != nil || user.Ga != "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "SYSTEM_ERROR")
		return
	}

	//if user.RegistrationType == "email" {
	//	username = user.Email
	//} else {
		username = user.Phone
	//}

	// 生成GA
	ga := common.NewGA()
	err = ga.Generate(username)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ga Generate err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SERVICE_TEMPORARILY_UNAVAILABLE.Code(), "SYSTEM_ERROR")
		return
	}

	// 放入redis,等待验证
	id, err := common.NewVerification("bind_ga", common.VerificationTypeGa).GenerateGA(userId, ga.Secret)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("NewVerification err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "GET_GA_FAILED")
		return
	}

	this.Response(ctx, &api.ResGaGenerate{id, ga.Secret, ga.Image})
}

func (this *GaController) Bind(ctx iris.Context) {

	params := new(api.GaBind)

	userId := common.GetUserIdFromCtx(ctx)
	err := ctx.ReadJSON(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Id == "", "id"},
		{params.Value == "", "value"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 获取用户信息
	user, err := new(models.User).GetById(int64(userId))
	if err != nil || user.Ga != "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "SYSTEM_ERROR")
		return
	}

	verification := common.NewVerification( "bind_ga", common.VerificationTypeGa)

	bol, err := verification.Verify(params.Id, userId, params.Value, "")
	if err != nil || !bol {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code(), "FAILURE_OF_GA_TOKEN_AUTHENTICATION")
		return
	}

	// 绑定
	err = new(models.User).BindGa(userId, verification.Value)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("BindGa err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "SYSTEM_ERROR")
		return
	}

	this.Response(ctx, nil)
	//ctx.Next()
}

func (this *GaController) UnBind(ctx iris.Context) {
	params := new(api.GaUnBind)

	userId := common.GetUserIdFromCtx(ctx)

	err := ctx.ReadJSON(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	// 获取用户信息
	user, err := new(models.User).GetById(int64(userId))
	if err != nil || user.Ga == "" {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserById err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), "SYSTEM_ERROR")
		return
	}

	// 验证token参数
	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Value == "", "value"},
		{user.Email != "" && params.EmailToken == "", "email_token"},
		{user.Phone != "" && params.SmsToken == "", "sms_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 验证GA
	ga := common.NewGA()
	bol, err := ga.Verify(user.Ga, params.Value)
	if err != nil || !bol {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_INCORRECT_GA_PWD.Code(), "FAILURE_OF_GA_TOKEN_AUTHENTICATION")
		return
	}

	// 验证邮箱
	if user.Email != "" {
		tokenPass, err := common.NewVerification("unbind_ga", common.VerificationTypeEmail).
			Check(params.EmailToken, userId, user.Email)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION")
			return
		}
	}

	// 验证手机
	if user.Phone != "" {
		tokenPass, err := common.NewVerification("unbind_ga", common.VerificationTypeSms).
			Check(params.SmsToken, userId, user.Phone)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
			}
			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_SMS_TOKEN_AUTHENTICATION")
			return
		}
	}

	// 解绑
	err = new(models.User).UnBindGa(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("UnBindGa err")
		this.ExceptionSerive(ctx,  apibackend.BASERR_DATABASE_ERROR.Code(), "SYSTEM_ERROR")
		return
	}

	this.Response(ctx, nil)
	//ctx.Next()
}
