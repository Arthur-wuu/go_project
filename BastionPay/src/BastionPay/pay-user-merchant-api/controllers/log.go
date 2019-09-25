package controllers

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/config"
	"BastionPay/pay-user-merchant-api/models"
	"BastionPay/pay-user-merchant-api/utils"
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	"BastionPay/pay-user-merchant-api/api"
)

type LogController struct {
	Controllers
}

func NewLogController() *LogController {
	return &LogController{
	}
}

func (l *LogController) RecodeLog(ctx iris.Context) {
	var (
		path   string
		userId int
		ip     string
		method string
	)

	userId = int(common.GetUserIdFromCtx(ctx))
	path = ctx.Path()
	ip = common.GetRealIp(ctx)
	method = ctx.Method()

	ZapLog().With(zap.Int("userid", userId), zap.String("path", path), zap.String("ip", ip), zap.String("method", method)).Info("")
	//	glog.Info(userId, path, ip, method)

	if method == "POST" {
		switch path {
		case "/v1/user/account/login", "/account/login/ga":
			device := ctx.Values().GetString("device")
			go l.RecodeLoginLog(userId, ip, device)
			break
		case "/v1/user/account/register":
			go l.RecodeOperationLog(userId, ip, "OPERATION_REGISTER")
			break
		case "/v1/user/account/info/email":
			go l.RecodeOperationLog(userId, ip, "OPERATION_BIND_EMAIL")
			break
		case "/v1/user/account/info/phone":
			go l.RecodeOperationLog(userId, ip, "OPERATION_BIND_PHONE")
			break
		case "/v1/user/account/info/phone/rebind":
			go l.RecodeOperationLog(userId, ip, "OPERATION_REBIND_PHONE")
			break
		case "/v1/user/account/password/reset":
			go l.RecodeOperationLog(userId, ip, "OPERATION_RESET_PASSWORD")
			break
		case "/v1/user/account/password/modify":
			go l.RecodeOperationLog(userId, ip, "OPERATION_MODIFY_PASSWORD")
			break
		case "/v1/user/account/ga/bind":
			go l.RecodeOperationLog(userId, ip, "OPERATION_BIND_GA")
			break
		case "/v1/user/account/ga/unbind":
			go l.RecodeOperationLog(userId, ip, "OPERATION_UNBIND_GA")
			break
		}
	}
}

func (l *LogController) RecodeLoginLog(userId int, ip string, device string) {
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(config.GConfig.IpFind.Auth, ip)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("IpLocation err")
		//		glog.Error(err.Error())
		return
	}

	if ipInfo.Country == "" {
		country = "Unknown"
	} else {
		country = ipInfo.Country
	}

	if ipInfo.City == "" {
		city = "Unknown"
	} else {
		city = ipInfo.City
	}

	if device == "" {
		device = "web"
	}

	 err = new(models.LogLogin).ParseAdd(userId, ip, country, city, device).Add()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateLoginLog err")
		return
	}
}

func (l *LogController) RecodeOperationLog(userId int, ip string, operation string) {
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(config.GConfig.IpFind.Auth, ip)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("IpLocation err")
		//		glog.Error(err.Error())
		return
	}

	if ipInfo.Country == "" {
		country = "Unknown"
	} else {
		country = ipInfo.Country
	}

	if ipInfo.City == "" {
		city = "Unknown"
	} else {
		city = ipInfo.City
	}

	err = new(models.LogOperation).ParseAdd(userId, operation, ip, country, city).Add()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateOperationLog err")
		return
	}
}

func (this *LogController) GetLoginLog(ctx iris.Context) {
	var (
		userId uint
		err    error

		params api.PageParams

	)

	common.GetParams(ctx, &params)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	userId = common.GetUserIdFromCtx(ctx)

	data, count, err := new(models.LogLogin).List(userId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetLoginLog err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "QUERY_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}
	respData := make([]*api.ResLoginLog, len(data))

	for k, v := range data {
		respData[k] = &api.ResLoginLog {
			Id: v.Id,
			Ip: v.Ip,
			Country: v.Country,
			City:    v.City,
			Device:  v.Device,
			CreatedAt: v.CreatedAt,
		}
	}

	this.ResponsePage(ctx, respData,count, params.Page)
}

func (this *LogController) GetOperationLog(ctx iris.Context) {
	var (
		userId uint
		err    error
		params api.PageParams
	)

	common.GetParams(ctx, &params)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	userId = common.GetUserIdFromCtx(ctx)

	data, count, err := new(models.LogOperation).List(userId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetOperationLog err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "QUERY_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	respData := make([]*api.ResOperationLog, len(data))

	for k, v := range data {
		operation := i18n.Translate(ctx, v.Operation)
		if operation == "" {
			operation = v.Operation
		}
		respData[k] = &api.ResOperationLog{
			Id:        v.Id,
			Operation: &v.Operation,
			Ip :       &v.Ip,
			Country:   &v.Country,
			City:      &v.City,
			CreatedAt: v.CreatedAt,
		}
	}

	this.ResponsePage(ctx, respData,count, params.Page)
}
