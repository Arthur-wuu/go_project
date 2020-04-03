package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/pay-user-merchant-api/api"
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/config"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
	"time"
)

const CONST_CTX_APPUSERINFO = "app_userinfo"

type Interceptor struct {
	Controllers
}

func (this *Interceptor) VerifyAccess(ctx context.Context) {
	if ctx.Method() == "OPTIONS" {
		ctx.WriteString("OPTIONS")
		return
	}

	pass, _ := common.DoLimiter(&common.GRedis, config.GPreConfig.PathLimits, ctx)
	if !pass {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "ACCESS_IS_TOO_FREQUENT", apibackend.BASERR_OPERATE_FREQUENT.Code()))
		return
	}

	_, inWhiteList := config.GPreConfig.PathWhiteList[ctx.Path()]
	if inWhiteList {
		ctx.Next()
		return
	}

	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		this.ExceptionSerive(ctx, apibackend.BASERR_TOKEN_INVALID.Code(), apibackend.BASERR_TOKEN_INVALID.Desc())
		return
	}

	basErr, userInfo := this.GetUserInfoByToken(tokenString)
	if basErr.Code() != 0 {
		return
	}

	ctx.Values().Set(CONST_CTX_APPUSERINFO, userInfo)

	//可以放在单独队列里，增加速度
	common.GRedis.GetConn().Expire(CONNST_TOKEN_Prefix+tokenString, time.Duration(config.GConfig.Token.Expirat)*time.Minute)
	ctx.Next()
}

func (this *Interceptor) GetUserInfoByToken(tokenString string) (apibackend.EnumBasErr, *api.BkLoginUserInfo) {
	userInfoBytes, err := common.GRedis.GetConn().Get(CONNST_TOKEN_Prefix + tokenString).Bytes()
	if err == redis.Nil {
		return apibackend.BASERR_TOKEN_EXPIRED, nil
	}
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis get err")
		return apibackend.BASERR_DATABASE_ERROR, nil
	}

	if len(userInfoBytes) < 3 {
		return apibackend.BASERR_TOKEN_PENDING, nil
	}
	userInnfo := new(api.BkLoginUserInfo)
	if err := json.Unmarshal(userInfoBytes, userInnfo); err != nil {
		ZapLog().With(zap.Error(err)).Error("json.Unmarshal err")
		return apibackend.BASERR_DATA_UNPACK_ERROR, nil
	}

	if userInnfo.Status == nil || *userInnfo.Status != 1 {
		return apibackend.BASERR_TOKEN_PENDING, nil
	}

	return apibackend.BASERR_SUCCESS, userInnfo
}
