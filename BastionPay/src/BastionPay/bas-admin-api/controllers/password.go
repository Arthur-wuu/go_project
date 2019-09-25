package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
)

type PasswordController struct {
	redis       *common.Redis
	db          *gorm.DB
	config      *config.Config
	userModel   *models.UserModel
	secretModel *models.SecretModel
}

func NewPasswordController(redis *common.Redis, db *gorm.DB, config *config.Config) *PasswordController {
	return &PasswordController{
		redis: redis,
		db:    db, config: config,
		userModel:   models.NewUserModel(db),
		secretModel: models.NewSecretModel(db),
	}
}

func (p *PasswordController) Modify(ctx context.Context) {
	var (
		secretId uint
		params   = struct {
			OldPassword string `json:"old_password"`
			Password    string
		}{}
		err error
	)
	userId := common.GetUserIdFromCtx(ctx)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	// 验证参数
	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.OldPassword == "", "old_password"},
		{params.Password == "", "password"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	user, err := p.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	verify, err := p.secretModel.Verify(user.SecretID, params.OldPassword)
	if err != nil || !verify {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("Verify err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "INCORRECT_OLD_PASSWORD", apibackend.BASERR_INCORRECT_PWD.Code()))
		return
	}

	// 创建密码
	secretId, err = p.secretModel.CreateSecret(params.Password)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateSecret err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_PASSWORD_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	// 修改密码
	if err = p.userModel.UpdateSecretId(userId, secretId); err != nil {
		ZapLog().With(zap.Error(err)).Error("UpdateSecretId err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "RESET_PASSWORD_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))
	ctx.Next()
}

func (p *PasswordController) Inquire(ctx context.Context) {
	var (
		params = struct {
			CaptchaToken string `json:"captcha_token"`
			Username     string
		}{}
		result = struct {
			Email       string `json:"email"`
			Phone       string `json:"phone"`
			CountryCode string `json:"country_code"`
			Ga          bool   `json:"ga"`
		}{}
		err error
	)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.CaptchaToken == "", "captcha_token"},
		{params.Username == "", "username"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 图形验证码
	captchaPass, err := common.NewVerification(p.redis, "forget_password", common.VerificationTypeCaptcha).
		Check(params.CaptchaToken, 0, "")
	if err != nil || !captchaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	// 用户是否存在
	user, err := p.userModel.GetUserByName(params.Username)
	if err != nil || user == nil {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserByName err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_USER_NOT_EXISTED", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	result.Email = user.Email
	result.Phone = user.Phone
	result.CountryCode = user.CountryCode

	if user.Ga != "" {
		result.Ga = true
	}

	ctx.JSON(common.NewSuccessResponse(ctx, result))
}

func (p *PasswordController) Reset(ctx context.Context) {
	var (
		secretId uint
		params   = struct {
			Username   string
			Password   string
			EmailToken string `json:"email_token"`
			SmsToken   string `json:"sms_token"`
			GaValue    string `json:"ga_value"`
		}{}
		err error
	)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	if params.Username == "" {
		ZapLog().Error("Username is nil err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", "username"))
		return
	}

	// 用户是否存在
	user, err := p.userModel.GetUserByName(params.Username)
	if err != nil || user == nil {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("GetUserByName err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_USER_NOT_EXISTED", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	// 验证token参数
	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Password == "", "password"},
		{user.Email != "" && params.EmailToken == "", "email_token"},
		{user.Phone != "" && params.SmsToken == "", "sms_token"},
		{user.Ga != "" && params.GaValue == "", "ga_value"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 验证邮箱
	if user.Email != "" {
		tokenPass, err := common.NewVerification(p.redis, "forget_password", common.VerificationTypeEmail).
			Check(params.EmailToken, 0, user.Email)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_EMAIL_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 验证手机
	if user.Phone != "" {
		tokenPass, err := common.NewVerification(p.redis, "forget_password", common.VerificationTypeSms).
			Check(params.SmsToken, 0, user.CountryCode+user.Phone)
		if err != nil || !tokenPass {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("NewVerification err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_SMS_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
			return
		}
	}

	// 验证GA
	if user.Ga != "" {
		ga := common.NewGA()
		bol, err := ga.Verify(user.Ga, params.GaValue)
		if err != nil || !bol {
			if err != nil {
				ZapLog().With(zap.Error(err)).Error("Verify err")
				//				glog.Error(err.Error())
			}
			ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_GA_VALUE_AUTHENTICATION", apibackend.BASERR_INCORRECT_GA_PWD.Code()))
			return
		}
	}

	// 创建密码
	secretId, err = p.secretModel.CreateSecret(params.Password)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateSecret err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_PASSWORD_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	if err = p.userModel.UpdateSecretId(user.ID, secretId); err != nil {
		ZapLog().With(zap.Error(err)).Error("model UpdateSecretId err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "RESET_PASSWORD_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	if _, err = common.NewTermBlocker(p.redis, p.config.TermBlockLimits, "login_pwd_incorrect", user.Uuid, ctx).Done(true); err != nil {
		ZapLog().With(zap.Error(err)).Error("TermBlocker Done err")
	}

	ctx.JSON(common.NewSuccessResponse(ctx, nil))

	ctx.Values().Set("uid", user.ID)
	ctx.Values().Set("safe", true)
	ctx.Next()
}
