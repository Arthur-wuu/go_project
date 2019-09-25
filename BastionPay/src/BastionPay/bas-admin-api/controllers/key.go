package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

type KeyController struct {
	config *config.Config
}

func NewKeyController(config *config.Config) *KeyController {
	return &KeyController{
		config: config,
	}
}

func (k *KeyController) Report(ctx iris.Context) {
	params := struct {
		UserKey     string `json:"user_key"`
		PublicKey   string `json:"public_key"`
		SourceIp    string `json:"source_ip"`
		CallbackUrl string `json:"callback_url"`
	}{}

	err := ctx.ReadJSON(&params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("ReadJSON err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "READ JSON PARAMS ERRORS", apibackend.BASERR_INVALID_PARAMETER.Code()))
		return
	}

	user, err := common.JwtParse(k.config.Token.Secret, k.config.Token.Expiration, ctx.GetHeader("Authorization"))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JwtParse err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "JWT PARSE ERRORS", apibackend.BASERR_TOKEN_INVALID.Code()))
		return
	}

	params.UserKey = user.Uuid
	params.PublicKey = strings.Trim(params.PublicKey, "\r")

	body, err := json.Marshal(params)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("JsonMarshal err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "MARSHAL JSON ERRORS", apibackend.BASERR_DATA_PACK_ERROR.Code()))
		return
	}

	_, res, err := bastionpay.CallApi(user.Uuid, string(body), "/v1/account/updateprofile")
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("bastionpay.CallApi err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "POST WALLET SERER ERRORS", apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code()))
		return
	}

	ctx.JSON(common.NewSuccessResponse(ctx, res))
}
