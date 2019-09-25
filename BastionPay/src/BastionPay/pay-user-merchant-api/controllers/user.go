package controllers

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/config"
	"BastionPay/pay-user-merchant-api/models"
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	"BastionPay/pay-user-merchant-api/api"
	"fmt"
	"BastionPay/pay-user-merchant-api/bas-user-api"

	"time"
	"github.com/satori/go.uuid"
	"strings"
	"encoding/json"
	"github.com/go-redis/redis"
)

const(
	CONNST_QRCODE_Prefix = "QRCODE_"
	CONNST_TOKEN_Prefix = "TOKEN_"
	CONNST_TOKEN_UserInfo = "{}"
)

type UserController struct {
	Controllers
}

func NewUserController() *UserController {
	return &UserController{
	}
}

func (this *UserController) LoginQr(ctx iris.Context) {
	qrCode,err := new(bas_user_api.ScanLoginQrcode).Send()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("userapi ScanLoginQrcode err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), "userapi ScanLoginQrcode err")
		return
	}

	uuidSu,err := uuid.NewV4()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("uuid.NewV4 err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "token gen err")
		return
	}

	replaceStr := fmt.Sprintf("%d", time.Now().UnixNano()%100000000)
	uuidStr := strings.Replace(uuidSu.String(), "-", replaceStr, -1)

	//存
	batcher := common.GRedis.GetConn().Pipeline()
	err = batcher.Set(CONNST_QRCODE_Prefix+qrCode, uuidStr, 10*time.Minute).Err()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis Pipeline SetXX err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis Pipeline SetXX err")
		return
	}
	err = batcher.Set(CONNST_TOKEN_Prefix+uuidStr,[]byte(CONNST_TOKEN_UserInfo), 10*time.Minute).Err()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis Pipeline SetXX err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis Pipeline SetXX err")
		return
	}
	_,err = batcher.Exec()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis Pipeline Exec err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis Pipeline Exec err")
		return
	}

	res := &api.ResLoginQr{
		QrCode :qrCode,
		Token: uuidStr,
		//Expiration: config.GConfig.Token.Expiration,
	}
	this.Response(ctx, res)
}

func (this *UserController) BkLoginQrCallBack(ctx iris.Context) {
	param := new(api.BkLoginUserInfo)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	token, err := common.GRedis.GetConn().Get(CONNST_QRCODE_Prefix+*param.LoginQrcode).Result()
	if err == redis.Nil {
		ZapLog().Error("token is nil, qrcode store in redis must have expire")
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), "nofind qrcode err")
		return
	}
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis Get err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis Get err")
		return
	}
	if len(token) <= 1 {
		ZapLog().Error("token len is 0, qrcode store bug")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "bug")
		return
	}

	param.Status = new(int)
	*param.Status = 1

	userInfo,err := json.Marshal(param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("json Marshal err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATA_PACK_ERROR.Code(), "json Marshal err")
		return
	}

	common.GRedis.GetConn().SetXX(CONNST_TOKEN_Prefix+token, userInfo, time.Duration(config.GConfig.Token.Expirat)*time.Minute)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis  SetXX err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis SetXX err")
		return
	}

	this.Response(ctx, nil)
}

func (this *UserController) LoginQrCheck(ctx iris.Context) {
	param := new(api.LoginQrCheck)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	userInfoBytes,err := common.GRedis.GetConn().Get(CONNST_TOKEN_Prefix+param.Token).Bytes()
	if err == redis.Nil {
		ZapLog().With(zap.Error(err)).Error("redis Get err")
		this.ExceptionSerive(ctx, apibackend.BASERR_TOKEN_INVALID.Code(), "redis get err")
		return
	}
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis get err")
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "redis get err")
		return
	}

	if len(userInfoBytes) <= 3 {
		//ZapLog().Error("redis err")
		this.ExceptionSerive(ctx, apibackend.BASERR_TOKEN_PENDING.Code(), "pending")
		return
	}
	userInfo := new(api.BkLoginUserInfo)
	if err := json.Unmarshal(userInfoBytes, userInfo); err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Unmarshal err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "json.Unmarshal err")
		return
	}
	if userInfo.Status ==nil || *userInfo.Status != 1 {
		//ZapLog().Error("pending err")
		this.ExceptionSerive(ctx, apibackend.BASERR_TOKEN_PENDING.Code(), "pending")
		return
	}

	this.Response(ctx, userInfo)
}

func (this *UserController) Login(ctx iris.Context) {
	var (
		safe   bool
		params = new(api.Login)
		checkItems = struct {
			Email bool
			Phone bool
			Ga    bool
		}{}
	)

	err := Tools.ShouldBindJSON(ctx, params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	// 检测图形验证码
	captchaPass, err := common.NewVerification( "login", common.VerificationTypeCaptcha).
		Check(params.CaptchaToken, 0, "")
	if err != nil || !captchaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INCORRECT_VERIFYCODE.Code(), "FAILURE_OF_CAPTCHA_TOKEN_AUTHENTICATION")
		return
	}


	user, err := new(models.User).GetByPhone(params.CountryCode, params.Phone)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserByName err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "GET_USER_ERROR")
		return
	}

	if user == nil {
		ZapLog().Error("USER nofind err", zap.Any("countrycode", params.CountryCode), zap.Any("phone", params.Phone))
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	if user.Status != 1 {
		ZapLog().Error("USER_IS_LOCKED err")
		this.ExceptionSerive(ctx, apibackend.BASERR_BLOCK_ACCOUNT.Code(), "USER_IS_LOCKED")
		return
	}

	termBlock := common.NewTermBlocker(&common.GRedis, config.GPreConfig.TermBlockLimits, "login_pwd_incorrect", fmt.Sprintf("%d", user.Id), ctx)
	isBlock, err := termBlock.IsBlock()
	if err !=nil {
		ZapLog().With(zap.Error(err)).Error("TermBlocker IsBlock err")
	}else if isBlock{
		ZapLog().Error("USER_IS_LOCKED err")
		this.ExceptionSerive(ctx, apibackend.BASERR_BLOCK_ACCOUNT.Code(), "USER_IS_LOCKED")
		return
	}


	//verify, err := models.GSecretModel.Verify(user.SecretID, params.Password)
	//if err != nil {
	//	ZapLog().With(zap.Error(err)).Error("Verify err")
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "INCORRECT_USERNAME_OR_PASSWORD")
	//	return
	//}
	verify := true

	tbResper, err := termBlock.Done(verify)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("TermBlocker Done err")
	}



	if !verify {
		ZapLog().With(zap.Any("use", *user)).Error("Verify err")
		if !tbResper.OpenFlag {
			this.ExceptionSerive(ctx, apibackend.BASERR_INCORRECT_PWD.Code(), "INCORRECT_USERNAME_OR_PASSWORD")
			return
		}
		if tbResper.OnBlock {
			this.ExceptionSerive(ctx,  apibackend.BASERR_BLOCK_ACCOUNT.Code(), "USER_IS_LOCKED")
			return
		}
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INCORRECT_PWD.Code(), "INCORRECT_USERNAME_OR_PASSWORD_BLOCKER", tbResper.Remain_count, tbResper.Lock_time)
		return
	}

	if user.Ga == "" {
		safe = true
	} else {
		safe = false
	}

	//if user.Email != "" {
		checkItems.Email = false
	//}
	if user.Phone != "" {
		checkItems.Phone = true
	}
	if user.Ga != "" {
		checkItems.Ga = true
	}

	// 签发token
	token, expiration, err := common.JwtSign(config.GConfig.Token.Secret, config.GConfig.Token.Expiration,
		uint(*user.Id), safe, checkItems.Email, checkItems.Phone, checkItems.Ga, params.Phone, user.VipLevel)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "CREATE_TOKEN_ERROR")
		return
	}

	this.Response(ctx, &api.ResLogin{
		Token: token,
		Expiration: expiration,
		Safe: safe,
	})


	ctx.Values().Set("device", "web")
	//ctx.Values().Set("uid", user.ID)
	//ctx.Values().Set("safe", safe)

	ctx.Values().Set("app_claims", &common.AppClaims{
		UserId: uint(*user.Id),
		Safe:   safe,
		Email:  checkItems.Email,
		Phone:  checkItems.Phone,
		Ga:     checkItems.Ga,
	})
	//ctx.Next()
}

func (this *UserController) LoginWithGa(ctx iris.Context) {

	params := new(api.LoginGa)

	user, err := common.GetUserIdFromCtxUnsafe(ctx)
	if user.UserId <= 0 || err != nil {
		ZapLog().Sugar().Errorf("userid[%d] err[%v] AUTHENTICATION_FAILED failed", user.UserId, err)
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "AUTHENTICATION_FAILED")
		return
	}

	err = ctx.ReadJSON(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "FAILED_TO_GET_PARAMETERS")
		return
	}

	pv := common.CheckParams(&[]common.ParamsVerification{
		{params.GaToken == "", "ga_token"}})
	if pv != nil {
		ZapLog().With(zap.String("error", pv.ErrMsg)).Error("CheckParams err")
		this.ExceptionSeriveWithParams(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "MISSING_PARAMETERS",pv.ErrMsg)
		return
	}

	// 检测GA
	gaPass, err := common.NewVerification("login", common.VerificationTypeGa).
		Check(params.GaToken, user.UserId, "")
	if err != nil || !gaPass {
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("NewVerification err")
		}
		this.ExceptionSerive(ctx, apibackend.BASERR_ADMIN_INVALID_VERIFY_STATUS.Code(), "FAILURE_OF_GA_TOKEN_AUTHENTICATION")
		return
	}

	appClaims, err := common.GetAppClaims(ctx)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetAppClaims err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "CREATE_TOKEN_ERROR")
		return
	}
	appClaims.Safe = true

	// 签发token
	token, expiration, err := common.JwtSign(config.GConfig.Token.Secret, config.GConfig.Token.Expiration,
		user.UserId, appClaims.Safe, appClaims.Email, appClaims.Phone, appClaims.Ga, user.Uuid, 0)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "CREATE_TOKEN_ERROR")
		return
	}

	this.Response(ctx, &api.ResLoginGa{
		Token: token,
		Expiration: expiration,
	})

	ctx.Values().Set("device", "web")
	ctx.Values().Set("app_claims", appClaims)
	//ctx.Next()
}

func (this *UserController) RefreshToken(ctx iris.Context) {
	userId := common.GetUserIdFromCtx(ctx)

	appClaims, err := common.GetAppClaims(ctx)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetAppClaims err")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "CREATE_TOKEN_ERROR")
		return
	}

	user, err := new(models.User).GetById(int64(userId))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetUserById err")
		this.ExceptionSerive(ctx,  apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), "GET_USER_ERROR")
		return
	}

	if user.Status != 1 {
		this.ExceptionSerive(ctx, apibackend.BASERR_BLOCK_ACCOUNT.Code(), "USER_IS_LOCKED")
		return
	}

	token, expiration, err := common.JwtSign(config.GConfig.Token.Secret, config.GConfig.Token.Expiration,
		userId, appClaims.Safe, appClaims.Email, appClaims.Phone, appClaims.Ga, appClaims.Uuid, appClaims.VipLevel)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtSign err")
		this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "CREATE_TOKEN_ERROR")
		return
	}

	this.Response(ctx, &api.RefreshToken{token, expiration})
}
