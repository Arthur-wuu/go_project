package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-admin-api/models"
	"github.com/BastionPay/bas-admin-api/utils"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"go.uber.org/zap"
)

type LogController struct {
	db       *gorm.DB
	logModel *models.LogModel
	config   *config.Config
}

func NewLogController(db *gorm.DB, config *config.Config) *LogController {
	return &LogController{
		db:       db,
		config:   config,
		logModel: models.NewLogModel(db),
	}
}

func (l *LogController) RecodeLog(ctx iris.Context) {
	var (
		path   string
		userId uint
		ip     string
		method string
	)

	userId = common.GetUserIdFromCtx(ctx)
	path = ctx.Path()
	ip = common.GetRealIp(ctx)
	method = ctx.Method()

	ZapLog().With(zap.Uint("userid", userId), zap.String("path", path), zap.String("ip", ip), zap.String("method", method)).Info("")
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

func (l *LogController) RecodeLoginLog(userId uint, ip string, device string) {
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(l.config.IpFind.Auth, ip)
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

	_, err = l.logModel.CreateLoginLog(userId, ip, country, city, device)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateLoginLog err")
		//		glog.Error(err.Error())
		return
	}
}

func (l *LogController) RecodeOperationLog(userId uint, ip string, operation string) {
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(l.config.IpFind.Auth, ip)
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

	_, err = l.logModel.CreateOperationLog(userId, operation, ip, country, city)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("CreateOperationLog err")
		//		glog.Error(err.Error())
		return
	}
}

func (l *LogController) GetLoginLog(ctx iris.Context) {
	var (
		userId uint
		err    error
		params struct {
			common.PageParams
		}

		respData []interface{}
	)

	common.GetParams(ctx, &params)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	userId = common.GetUserIdFromCtx(ctx)

	data, count, err := l.logModel.GetLoginLog(userId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetLoginLog err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "QUERY_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	respData = make([]interface{}, len(data))

	for k, v := range data {
		respData[k] = &struct {
			Id        uint   `json:"id"`
			Ip        string `json:"ip"`
			Country   string `json:"country"`
			City      string `json:"city"`
			Device    string `json:"device"`
			CreatedAt int64  `json:"created_at"`
		}{v.ID, v.Ip, v.Country, v.City, v.Device, v.CreatedAt}
	}

	ctx.JSON(common.NewResponse(ctx).Success().SetMsg("SUCCESS").SetLimitResult(respData, count, params.Page))
}

func (l *LogController) GetOperationLog(ctx iris.Context) {
	var (
		userId uint
		err    error
		params struct {
			common.PageParams
		}

		respData []interface{}
	)

	common.GetParams(ctx, &params)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	userId = common.GetUserIdFromCtx(ctx)

	data, count, err := l.logModel.GetOperationLog(userId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("GetOperationLog err")
		//		glog.Error(err.Error())
		ctx.JSON(common.NewErrorResponse(ctx, nil, "QUERY_FAILED", apibackend.BASERR_DATABASE_ERROR.Code()))
		return
	}

	respData = make([]interface{}, len(data))

	for k, v := range data {
		operation := i18n.Translate(ctx, v.Operation)
		if operation == "" {
			operation = v.Operation
		}
		respData[k] = &struct {
			Id        uint   `json:"id"`
			Operation string `json:"operation"`
			Ip        string `json:"ip"`
			Country   string `json:"country"`
			City      string `json:"city"`
			CreatedAt int64  `json:"created_at"`
		}{v.ID, operation, v.Ip, v.Country, v.City, v.CreatedAt}
	}

	ctx.JSON(common.NewResponse(ctx).Success().SetMsg("SUCCESS").SetLimitResult(respData, count, params.Page))
}
