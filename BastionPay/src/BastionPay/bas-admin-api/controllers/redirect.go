package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

func NewRedirectController(config *config.Config) *RedirectController {
	bp := &RedirectController{
		config: config,
	}
	return bp
}

type RedirectController struct {
	config *config.Config
}

func (bp *RedirectController) HandlerV1Gateway(ctx iris.Context) { //未测试
	user, err := common.JwtParse(bp.config.Token.Secret, bp.config.Token.Expiration, ctx.GetHeader("Authorization"))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtParse err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_TOKEN_INVALID.Code()))
		return
	}

	newPath := strings.Trim(ctx.Path(), "\r")
	reqMsg, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("body read err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "READ_BODY_ERROR", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	_, res, err := bastionpay.CallApi(user.Uuid, string(reqMsg), newPath)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "GATEWAY_REQUEST_ERRORS", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, string(res)))
}
