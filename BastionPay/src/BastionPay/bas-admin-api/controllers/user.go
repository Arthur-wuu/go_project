package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/BastionPay/bas-tools/sdk.notify.mail"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

type UserController struct {
	redis       *common.Redis
	db          *gorm.DB
	config      *config.Config
	userModel   *models.UserModel
	secretModel *models.SecretModel
}

func NewUserController(redis *common.Redis, db *gorm.DB, config *config.Config) *UserController {
	models.GlobalBasModel.Init(config)
	return &UserController{
		redis:       redis,
		db:          db,
		config:      config,
		userModel:   models.NewUserModel(db),
		secretModel: models.NewSecretModel(db),
	}
}

func (u *UserController) Login(ctx iris.Context) {
	var (
		safe   bool
		params = struct {
			Username     string
			Password     string
			CaptchaToken string `json:"captcha_token"`
		}{}
		checkItems = struct {
			Email bool
			Phone bool
			Ga    bool
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
		{params.Username == "", "username"},
		{params.Password == "", "password"},
		{params.CaptchaToken == "", "captcha_token"}})

	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 检测图形验证码
	captchaPass, err := common.NewVerification(u.redis, "login", common.VerificationTypeCaptcha).
		Check(params.CaptchaToken, 0, "")
	if err != nil || !captchaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code()))
		return
	}

	user, err := u.userModel.GetUserByName(params.Username)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserByName err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	if user.Blocked == true {
		ZapLog().Error("USER_IS_LOCKED err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "USER_IS_LOCKED", apibackend.BASERR_BLOCK_ACCOUNT.Code()))
		return
	}

	termBlock := common.NewTermBlocker(u.redis, u.config.TermBlockLimits, "login_pwd_incorrect", user.Uuid, ctx)
	isBlock, err := termBlock.IsBlock()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("TermBlocker IsBlock err")
	} else if isBlock {
		ZapLog().Error("USER_IS_LOCKED err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "USER_IS_LOCKED", apibackend.BASERR_BLOCK_ACCOUNT.Code()))
		return
	}

	verify, err := u.secretModel.Verify(user.SecretID, params.Password)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Verify err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "INCORRECT_USERNAME_OR_PASSWORD", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	tbResper, err := termBlock.Done(verify)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("TermBlocker Done err")
	}

	if !verify {
		ZapLog().With(zap.Any("use", *user)).Error("Verify err")
		if !tbResper.OpenFlag {
			ctx.JSON(common.NewErrorResponse(ctx, nil, "INCORRECT_USERNAME_OR_PASSWORD", apibackend.BASERR_INCORRECT_PWD.Code()))
			return
		}
		if tbResper.OnBlock {
			ctx.JSON(common.NewErrorResponse(ctx, nil, "USER_IS_LOCKED", apibackend.BASERR_BLOCK_ACCOUNT.Code()))
			return
		}
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INCORRECT_PWD.Code()).SetMsgWithParams("INCORRECT_USERNAME_OR_PASSWORD_BLOCKER", tbResper.Remain_count, tbResper.Lock_time))
		return
	}

	if user.Ga == "" {
		safe = true
	} else {
		safe = false
	}

	if user.Email != "" {
		checkItems.Email = true
	}
	if user.Phone != "" {
		checkItems.Phone = true
	}
	if user.Ga != "" {
		checkItems.Ga = true
	}

	// 签发token
	token, expiration, err := common.JwtSign(u.config.Token.Secret, u.config.Token.Expiration,
		user.ID, safe, checkItems.Email, checkItems.Phone, checkItems.Ga, user.Uuid, user.VipLevel)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code()))
		return
	}

	userAllStatus, err := models.GlobalBasModel.GetUserAllAccountStatus(user.Uuid)
	if err != nil { //为了兼容，暂时不处理
		ZapLog().With(zap.Error(err), zap.String("userkey", user.Uuid)).Error("getUserAuditeStatus err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", common.ResponseError))
		//		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Token          string `json:"token"`
		Expiration     int64  `json:"expiration"`
		Safe           bool   `json:"safe"`
		AuditeStatus   uint   `json:"audite_status"`
		TransferStatus uint   `json:"transfer_status"`
	}{token, expiration, safe, userAllStatus.AuditeStatus, userAllStatus.TransferStatus}))

	ctx.Values().Set("device", "web")
	//ctx.Values().Set("uid", user.ID)
	//ctx.Values().Set("safe", safe)

	ctx.Values().Set("app_claims", &common.AppClaims{
		UserId: user.ID,
		Safe:   safe,
		Email:  checkItems.Email,
		Phone:  checkItems.Phone,
		Ga:     checkItems.Ga,
	})
	ctx.Next()
}

func (u *UserController) LoginWithGa(ctx iris.Context) {
	var (
		params = struct {
			GaToken string `json:"ga_token"`
		}{}
		//checkItems = struct {
		//	Email bool
		//	Phone bool
		//	Ga    bool
		//}{}
		err error
	)

	user, err := common.GetUserIdFromCtxUnsafe(ctx)
	if user.UserId <= 0 || err != nil {
		ZapLog().Sugar().Errorf("userid[%d] err[%v] AUTHENTICATION_FAILED failed", user.UserId, err)
		ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_UNKNOWN_BUG.Code()))
		return
	}

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.GaToken == "", "ga_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	// 检测GA
	gaPass, err := common.NewVerification(u.redis, "login", common.VerificationTypeGa).
		Check(params.GaToken, user.UserId, "")
	if err != nil || !gaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_GA_TOKEN_AUTHENTICATION", apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()))
		return
	}

	appClaims, err := common.GetAppClaims(ctx)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetAppClaims err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_UNKNOWN_BUG.Code()))
		return
	}
	appClaims.Safe = true

	// 签发token
	token, expiration, err := common.JwtSign(u.config.Token.Secret, u.config.Token.Expiration,
		user.UserId, appClaims.Safe, appClaims.Email, appClaims.Phone, appClaims.Ga, user.Uuid, 0)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code()))
		return
	}

	userAllStatus, err := models.GlobalBasModel.GetUserAllAccountStatus(user.Uuid)
	if err != nil { //为了兼容，暂时不处理
		ZapLog().With(zap.Error(err), zap.String("userkey", user.Uuid)).Error("getUserAuditeStatus err")
		//		ctx.JSON(common.NewErrorResponse(ctx, nil, "SERVER_INTERNAL_ERROR", common.ResponseError))
		//		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Token          string `json:"token"`
		Expiration     int64  `json:"expiration"`
		AuditeStatus   uint   `json:"audite_status"`
		TransferStatus uint   `json:"transfer_status"`
	}{token, expiration, userAllStatus.AuditeStatus, userAllStatus.TransferStatus}))

	ctx.Values().Set("device", "web")
	ctx.Values().Set("app_claims", appClaims)
	ctx.Next()
}

func (u *UserController) Register(ctx iris.Context) {
	var (
		secretId         uint
		username         string
		recipient        string
		verificationType string
		registrationType string
		tokenErrCode     int

		params = struct {
			CompanyName string `json:"company_name"`
			Email       string
			Phone       string
			CountryCode string `json:"country_code"`
			Password    string
			Citizenship string
			Language    string
			Timezone    string
			Token       string
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
		{params.Email == "" && params.Phone == "", "email or phone"},
		{params.Phone != "" && params.CountryCode == "", "country_code"},
		{params.Password == "", "password"},
		{params.Token == "", "token"},
		{params.CompanyName == "", "company_name"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	if params.Phone != "" {
		username = params.Phone
		recipient = params.CountryCode + params.Phone
		tokenErrCode = apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()
		verificationType = common.VerificationTypeSms
		registrationType = "phone"
	} else {
		username = params.Email
		recipient = params.Email
		tokenErrCode = apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code()
		verificationType = common.VerificationTypeEmail
		registrationType = "email"
	}

	// 检测手机或邮箱验证码
	tokenPass, err := common.NewVerification(u.redis, "register", verificationType).Check(params.Token, 0, recipient)
	if err != nil || !tokenPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILURE_OF_TOKEN_AUTHENTICATION", tokenErrCode))
		return
	}

	// 用户是否存在
	exists, err := u.userModel.UserExisted(username)
	if err != nil || exists {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("UserExisted err")
			//			glog.Error(err.Error())
		}
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_USER_HAS_ALREADY_EXISTED", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}
	compayNameNoSpace := strings.TrimSpace(params.CompanyName)
	if len(compayNameNoSpace) == 0 {
		ZapLog().With(zap.Error(err), zap.String("company", params.CompanyName)).Error("CompanyName blank err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "COMPANY_NAME_BLANK_ERR", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	exists, err = u.userModel.CompanyExisted(params.CompanyName)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("db err")
		ctx.JSON(common.NewErrorResponse(ctx, nil, "INNER_SERVER_ERR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}
	if exists {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "THE_COMPANY_HAS_ALREADY_EXISTED", apibackend.BASERR_OBJECT_EXISTS.Code()))
		return
	}

	// 创建密码
	secretId, err = u.secretModel.CreateSecret(params.Password)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateSecret err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_PASSWORD_ERROR", apibackend.BASERR_INCORRECT_FORMAT.Code()))
		return
	}

	//请求服务返回uuid
	b, _ := json.Marshal(map[string]interface{}{
		"user_class":   0,
		"level":        0,
		"is_frozen":    1,
		"user_name":    params.CompanyName,
		"user_mobile":  params.Phone,
		"user_email":   params.Email,
		"country_code": params.CountryCode,
		"language":     params.Language,
	})

	_, result, err := bastionpay.CallApi("", string(b), "/v1/account/register")
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_UUID_ERROR", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	res := struct {
		UserKey string `json:"user_key"`
	}{}

	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Unmarshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "PARSE_USER_KEY_ERROR", apibackend.BASERR_DATA_UNPACK_ERROR.Code()))
		return
	}

	// 创建账户
	userId, err := u.userModel.CreateUser(params.CompanyName, params.Email, params.Phone, params.CountryCode, secretId,
		params.Citizenship, params.Language, params.Timezone, registrationType, res.UserKey)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateUser err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_USER_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	// 签发token
	token, expiration, err := common.JwtSign(u.config.Token.Secret, u.config.Token.Expiration,
		userId, true, registrationType == "email", registrationType == "phone", false, res.UserKey, 0)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Token      string `json:"token"`
		Expiration int64  `json:"expiration"`
	}{token, expiration}))

	ctx.Values().Set("uid", userId)
	ctx.Values().Set("safe", true)
	ctx.Values().Set("app_claims", &common.AppClaims{
		UserId: userId,
		Safe:   true,
		Email:  registrationType == "email",
		Phone:  registrationType == "phone",
		Ga:     false,
	})
	ctx.Next()

	go func() {
		notifyConf := &u.config.Bas_notify
		lang := params.Language
		if lang != "zh-CN" {
			lang = "en-US"
		}

		if len(params.Phone) != 0 {
			err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(notifyConf.RegisterSuccessSmsTmp, lang, []string{params.CountryCode + params.Phone}, nil)
		} else {
			err = sdk_notify_mail.GNotifySdk.SendMailByGroupName(notifyConf.RegisterSuccessMailTmp, lang, []string{params.Email}, nil)
		}
		if err != nil {
			ZapLog().With(zap.Error(err), zap.Any("params", params)).Error("sdk_notify_mail err")
		} else {
			ZapLog().With(zap.Any("params", params)).Info("register notify ok")
		}
	}()
}

func (u *UserController) Exists(ctx iris.Context) {
	var (
		params = struct {
			Username     string
			CaptchaToken string `json:"captcha_token"`
		}{}
		err error
	)

	err = ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//	glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "FAILED_TO_GET_PARAMETERS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.Username == "", "username"},
		{params.CaptchaToken == "", "captcha_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		ctx.JSON(common.NewResponse(ctx).Error(apibackend.BASERR_INVALID_PARAMETER.Code()).
			SetMsgWithParams("MISSING_PARAMETERS", pv.ErrMsg))
		return
	}

	captchaPass, err := common.NewVerification(u.redis, "register", common.VerificationTypeCaptcha).
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
	exists, err := u.userModel.UserExisted(params.Username)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("UserExisted err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "SYSTEM_ERROR", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	//	glog.Info(exists)
	ZapLog().Sugar().Info(exists)
	ctx.JSON(common.NewSuccessResponse(ctx, exists))
}

func (u *UserController) RefreshToken(ctx iris.Context) {
	userId := common.GetUserIdFromCtx(ctx)

	appClaims, err := common.GetAppClaims(ctx)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetAppClaims err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_UNKNOWN_BUG.Code()))
		return
	}

	user, err := u.userModel.GetUserById(userId)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GET_USER_ERROR", apibackend.BASERR_ACCOUNT_NOT_FOUND.Code()))
		return
	}

	if user.Blocked == true {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "USER_IS_LOCKED", apibackend.BASERR_BLOCK_ACCOUNT.Code()))
		return
	}

	token, expiration, err := common.JwtSign(u.config.Token.Secret, u.config.Token.Expiration,
		userId, appClaims.Safe, appClaims.Email, appClaims.Phone, appClaims.Ga, appClaims.Uuid, appClaims.VipLevel)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "CREATE_TOKEN_ERROR", apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, &struct {
		Token      string `json:"token"`
		Expiration int64  `json:"expiration"`
	}{token, expiration}))
}
