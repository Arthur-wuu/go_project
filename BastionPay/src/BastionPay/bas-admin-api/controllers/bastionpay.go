package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

type BastionPayController struct {
	redis             *common.Redis
	config            *config.Config
	supportedFunction map[string]string
}

func NewBastionPayController(redis *common.Redis, config *config.Config) *BastionPayController {
	bp := &BastionPayController{
		redis:             redis,
		config:            config,
		supportedFunction: make(map[string]string),
	}

	for _, path := range config.WalletPaths {
		index := strings.LastIndex(path, "/")

		relativePath := path[0:index]
		function := path[index+1:]

		bp.supportedFunction[function] = relativePath
	}

	return bp
}

func (bp *BastionPayController) User(ctx iris.Context) {
	params := struct {
		Function string      `json:"function"`
		Message  interface{} `json:"message"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_GET_PARAMS_ERRORS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	user, err := common.JwtParse(bp.config.Token.Secret, bp.config.Token.Expiration, ctx.GetHeader("Authorization"))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtParse err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "AUTHENTICATION_FAILED", apibackend.BASERR_TOKEN_INVALID.Code()))
		return
	}

	params.Function = strings.Trim(params.Function, "\r")
	relativePath, ok := bp.supportedFunction[params.Function]
	if !ok {
		ZapLog().With(zap.String("Function", params.Function)).Error("NOT SUPPORT FUNCTION")
		//		glog.Errorf("NOT SUPPORT FUNCTION:(%s)", params.Function)
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_NOT_SUPPORT_FUNCTION", apibackend.BASERR_UNSUPPORTED_METHOD.Code()))
		return
	}

	reqMsg, err := json.Marshal(params.Message)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_GET_MESSAGE_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}

	if err := bp.preHook(ctx, params.Function, user, reqMsg); err != nil {
		return
	}

	_, res, err := bastionpay.CallApi(user.Uuid, string(reqMsg), relativePath+"/"+params.Function)
	if err != nil {
		errMsg := err.Error()
		if strings.HasSuffix(err.Error(), "EOF") {
			errMsg = "timeout"
		}
		ZapLog().With(zap.Error(err), zap.String("path", relativePath+"/"+params.Function)).Error("bastionpay.CallApi err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_"+errMsg, apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, string(res)))
}

func (bp *BastionPayController) HandlerV1(ctx iris.Context) {
	params := struct {
		Message interface{} `json:"message"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_GET_PARAMS_ERRORS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}
	var user *common.AppClaims
	appclaims := ctx.Values().Get("app_claims")
	if appclaims != nil {
		ok := false
		user, ok = appclaims.(*common.AppClaims)
		if !ok {
			ZapLog().Error("app_claims type err")
			//		glog.Error(err.Error())
			ctx.JSON(common.NewErrorResponse(ctx, nil, "type err", apibackend.BASERR_UNKNOWN_BUG.Code()))
			return
		}
	} else {
		user = new(common.AppClaims)
	}

	path := strings.Trim(ctx.Path(), "\r")
	index := strings.LastIndex(path, "/")
	if index == -1 {
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_NOT_SUPPORT_FUNCTION", apibackend.BASERR_UNSUPPORTED_METHOD.Code()))
		return
	}

	function := path[index+1:]
	function = strings.Trim(function, "\r")

	relativePath, ok := bp.supportedFunction[function]
	if !ok {
		ZapLog().With(zap.String("Function", function)).Error("NOT SUPPORT FUNCTION")
		//		glog.Errorf("NOT SUPPORT FUNCTION:(%s)", params.Function)
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_NOT_SUPPORT_FUNCTION", apibackend.BASERR_UNSUPPORTED_METHOD.Code()))
		return
	}

	reqMsg, err := json.Marshal(params.Message)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Marshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_GET_MESSAGE_ERROR", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}

	if err := bp.preHook(ctx, function, user, reqMsg); err != nil {
		return
	}

	_, res, err := bastionpay.CallApi(user.Uuid, string(reqMsg), relativePath+"/"+function)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.String("path", relativePath+"/"+function)).Error("bastionpay.CallApi err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "BASTIONPAY_REQUEST_ERRORS", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, string(res)))
}

func (bp *BastionPayController) preHook(ctx iris.Context, function string, user *common.AppClaims, reqMsg []byte) error {

	if function == "deposit_order" || function == "withdrawal_order" || function == "transaction_bill" {
		reqParam := &struct {
			Async      int `json:"is_asyn"`
			Page_index int `json:"page_index"`
		}{}
		if err := json.Unmarshal([]byte(reqMsg), reqParam); err != nil {
			ZapLog().With(zap.Error(err)).Error("json.Unmarshal err")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_INVALID_PARAMETER.Code()))
			return err
		}
		if reqParam.Async != 1 || reqParam.Page_index != 1 {
			return nil
		}
		flag, remainCount, limitTime, err := common.DoLevelLimiter(bp.redis, bp.config.LevelPathLimits, ctx, "download_bill_turnover_file", user.Uuid, int(user.VipLevel))
		if err != nil {
			ZapLog().With(zap.Error(err)).Error("DoLevelLimiter err")
			ctx.JSON(common.NewErrorResponse(ctx, nil, "", apibackend.BASERR_DATABASE_ERROR.Code()))
			return err
		}
		if !flag {
			ctx.JSON(common.NewErrorResponse(ctx, nil, fmt.Sprintf("%d-%d", remainCount, limitTime), apibackend.BASERR_OPERATE_FREQUENT.Code()))
			ZapLog().With(zap.Int("remainCount", remainCount), zap.Int("limitTime", limitTime), zap.String("user.Uuid", user.Uuid)).Error("DoLevelLimiter over limit")
			return errors.New("limited")
		}
	}

	return nil
}
