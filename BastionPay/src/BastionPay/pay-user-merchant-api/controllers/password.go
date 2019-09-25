package controllers

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/models"
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	"BastionPay/pay-user-merchant-api/api"
)

type PasswordController struct {
	Controllers
}

func NewPasswordController() *PasswordController {
	return &PasswordController{
	}
}

//func (this *PasswordController) Modify(ctx context.Context) {
//
//	params  := new(api.PasswordModify)
//	userId := common.GetUserIdFromCtx(ctx)
//
//	err := ctx.ReadJSON(&params)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
//		return
//	}
//
//	// 验证参数
//	pv := common.CheckParams(&[]common.ParamsVerification{
//		{params.OldPassword == "", "old_password"},
//		{params.Password == "", "password"}})
//	if pv != nil {
//		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
//		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
//		return
//	}
//
//	user, err := new(models.User).GetById(int64(userId))
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("GetUserById err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "GET_USER_ERROR")
//		return
//	}
//
//	verify, err := models.GSecretModel.Verify(user.SecretID, params.OldPassword)
//	if err != nil || !verify {
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("Verify err")
//		}
//		this.ExceptionSerive(ctx, apibackend.BASERR_INCORRECT_PWD.Code(), "INCORRECT_OLD_PASSWORD")
//		return
//	}
//
//	// 创建密码
//	secretId, err := models.GSecretModel.CreateSecret(params.Password)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("CreateSecret err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "CREATE_PASSWORD_ERROR")
//		return
//	}
//
//	// 修改密码
//	if err = new(models.User).UpdateSecretId(userId, secretId); err != nil {
//		ZapLog().With(zap.Error(err)).Error("UpdateSecretId err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "RESET_PASSWORD_ERROR")
//		return
//	}
//
//	this.Response(ctx, nil)
//	ctx.Next()
//}

func (this *PasswordController) Inquire(ctx context.Context) {

	params := new(api.PasswordInquire)
	err := ctx.ReadJSON(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.CaptchaToken == "", "captcha_token"},
		{params.Username == "", "username"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
		return
	}

	// 图形验证码
	captchaPass, err := common.NewVerification( "forget_password", common.VerificationTypeCaptcha).
		Check(params.CaptchaToken, 0, "")
	if err != nil || !captchaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION")
		return
	}

	// 用户是否存在
	user, err := new(models.User).GetByName(params.Username)
	if err != nil || user == nil {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserByName err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "THE_USER_NOT_EXISTED")
		return
	}

	result := api.ResPasswordInquire{
		Email : user.Email,
		Phone: user.Phone,
		CountryCode: user.PhonneDistrict,
	}
	if user.Ga != "" {
		result.Ga = true
	}

	ctx.JSON(common.NewSuccessResponse(ctx, result))
}

//func (this *PasswordController) Reset(ctx context.Context) {
//	params   := new(api.PasswordReset)
//
//	err := ctx.ReadJSON(params)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
//		return
//	}
//
//	if params.Username == "" {
//		ZapLog().Error("Username is nil err")
//		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", "username")
//		return
//	}
//
//	// 用户是否存在
//	user, err := new(models.User).GetByName(params.Username)
//	if err != nil || user == nil {
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("GetUserByName err")
//		}
//		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "THE_USER_NOT_EXISTED")
//		return
//	}
//
//	// 验证token参数
//	pv := common.CheckParams(&[]common.ParamsVerification{
//		{params.Password == "", "password"},
//		{user.Email != "" && params.EmailToken == "", "email_token"},
//		{user.Phone != "" && params.SmsToken == "", "sms_token"},
//		{user.Ga != "" && params.GaValue == "", "ga_value"}})
//	if pv != nil {
//		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
//		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS", pv.ErrMsg)
//		return
//	}
//
//	// 验证邮箱
//	if user.Email != "" {
//		tokenPass, err := common.NewVerification("forget_password", common.VerificationTypeEmail).
//			Check(params.EmailToken, 0, user.Email)
//		if err != nil || !tokenPass {
//			if err != nil {
//				ZapLog().With(zap.Error(err)).Error("NewVerification err")
//			}
//			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION")
//			return
//		}
//	}
//
//	// 验证手机
//	if user.Phone != "" {
//		tokenPass, err := common.NewVerification("forget_password", common.VerificationTypeSms).
//			Check(params.SmsToken, 0, user.PhonneDistrict+user.Phone)
//		if err != nil || !tokenPass {
//			if err != nil {
//				ZapLog().With(zap.Error(err)).Error("NewVerification err")
//			}
//			this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_SMS_TOKEN_AUTHENTICATION")
//			return
//		}
//	}
//
//	// 验证GA
//	if user.Ga != "" {
//		ga := common.NewGA()
//		bol, err := ga.Verify(user.Ga, params.GaValue)
//		if err != nil || !bol {
//			if err != nil {
//				ZapLog().With(zap.Error(err)).Error("Verify err")
//			}
//			this.ExceptionSerive(ctx, apibackend.BASERR_INCORRECT_GA_PWD.Code(), "FAILURE_OF_GA_VALUE_AUTHENTICATION")
//			return
//		}
//	}
//
//	// 创建密码
//	secretId, err := models.GSecretModel.CreateSecret(params.Password)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("CreateSecret err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "CREATE_PASSWORD_ERROR")
//		return
//	}
//
//	if err = new(models.User).UpdateSecretId(uint(*user.Id), secretId); err != nil {
//		ZapLog().With(zap.Error(err)).Error("model UpdateSecretId err")
//		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "RESET_PASSWORD_ERROR")
//		return
//	}
//
//	if _,err = common.NewTermBlocker(&common.GRedis, config.GPreConfig.TermBlockLimits, "login_pwd_incorrect", user.Uuid, ctx).Done(true); err != nil {
//		ZapLog().With(zap.Error(err)).Error("TermBlocker Done err")
//	}
//
//	ctx.JSON(common.NewSuccessResponse(ctx, nil))
//
//	ctx.Values().Set("uid", *user.Id)
//	ctx.Values().Set("safe", true)
//	appClaims := &common.AppClaims{
//		UserId:uint(*user.Id),
//		Uuid:user.Uuid,
//		Safe:true,
//	}
//	ctx.Values().Set("app_claims", appClaims)
//	ctx.Next()
//}