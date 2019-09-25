package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
)

var (
	TokenWhiteList = []string{
		"/v1/user/account/login",
		"/v1/user/account/login/ga",
		"/v1/user/account/register",
		"/v1/user/account/exists",
		"/v1/user/account/password/reset",
		"/v1/user/account/password/inquire",
		"/v1/user/account/verification",
		"/v1/user/account/verification/email",
		"/v1/user/account/verification/sms",
		"/v1/user/account/verification/captcha",
		"/v1/user/account/verification/ga",
		"/v1/user/debug/pprof",
		"/v1/user/notice/getlist",
		"/v1/user/notice/get",
		"/v1/user/notice/count",
		"/v1/user/bastionpay/asset_attribute",
	}
)

type InterceptorController struct {
	redis  *common.Redis
	config *config.Config
}

func NewInterceptorController(redis *common.Redis, config *config.Config) *InterceptorController {
	return &InterceptorController{redis: redis, config: config}
}

func (i *InterceptorController) Interceptor(ctx context.Context) {
	var (
		tokenString string
		inWhiteList bool
		appClaims   *common.AppClaims
		err         error
	)

	//ctx.Header("Access-Control-Allow-Origin", "*")
	//ctx.Header("Access-Control-Allow-Headers", "Authorization,X-Requested-With,X_Requested_With,Content-Type,Access-Token,Accept-Language")
	//ctx.Header("Access-Control-Expose-Headers", "*")
	//ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//ctx.Header("Access-Control-Allow-Credentials","true")
	if ctx.Method() == "OPTIONS" {
		ctx.WriteString("OPTIONS")
		return
	}

	pass, _ := common.DoLimiter(i.redis, i.config.PathLimits, ctx)
	if !pass {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "ACCESS_IS_TOO_FREQUENT", apibackend.BASERR_OPERATE_FREQUENT.Code()))
		return
	}

	inWhiteList = common.InArray(TokenWhiteList, ctx.Path())
	ZapLog().Info("path:>>" + ctx.Path())
	//	glog.Info("path:>>", ctx.Path())

	tokenString = ctx.GetHeader("Authorization")
	if tokenString != "" {
		appClaims, err = common.JwtParse(i.config.Token.Secret, i.config.Token.Expiration, tokenString)
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("JwtParse err")
			//			glog.Error(err.Error())
			appClaims = &common.AppClaims{}
		}

		ZapLog().Sugar().Infof("%+v", appClaims)
		//		glog.Info(appClaims)
		if appClaims != nil {
			ctx.Values().Set("app_claims", appClaims)
		}
	}

	if inWhiteList {
		ctx.Next()
	} else {
		ZapLog().Sugar().Infof("%+v", appClaims)
		//		glog.Info(appClaims)
		if appClaims == nil || appClaims.UserId == 0 || appClaims.Safe == false {
			ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_TOKEN_INVALID.Code()))
		} else {
			ctx.Next()
		}
	}
}
