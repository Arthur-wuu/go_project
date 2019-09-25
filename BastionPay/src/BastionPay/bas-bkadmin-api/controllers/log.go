package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/i18n"
	"time"
)

var GLogController LogController

type LogController struct {
	db       *gorm.DB
	logModel *models.LogModel
	config   *tools.Config
}

func NewLogController(config *tools.Config) *LogController {

	GLogController.logModel = models.NewLogModel()
	GLogController.config = config

	//删除操作日志
	go GLogController.Start()

	return &GLogController
}

func (l *LogController) Start() {
	defer utils.PanicPrint()
	for {
		l.logModel.DeleteLoginLog(l.config.OperateLog.RemainDays)
		l.logModel.DeleteOperatLog(l.config.OperateLog.RemainDays)
		time.Sleep(time.Second * time.Duration(l.config.OperateLog.CleanIntvl))
	}

}

func (l *LogController) RecodeLog(ctx iris.Context) {
	var (
		path   string
		userId uint
		ip     string
		method string
	)
	userId = uint(utils.NewUtils().GetValueUserId(ctx))
	if userId <= 0 {
		//		l4g.Debug("RecodeLog userId is 0, %s", method+path)
		return
	}
	path = ctx.Path()
	ip = common.GetRealIp(ctx)
	method = ctx.Method()

	go l.RecodeOperationLog(userId, ip, method+path)
	return

	//go l.RecodeOperationLog(userId, ip, method+path)
	//return

	l4g.Debug("userid:[%d], path:[%s], ip:[%s], method:[%s]", userId, path, ip, method)

	if method == "POST" {
		switch path {
		case "/v1/account/register":
			go l.RecodeOperationLog(userId, ip, "/v1/account/register")
			break
		case "/v1/account/login":
			go l.RecodeLoginLog(userId, "ip 123", "/v1/account/login")
			break
		case "/v1/access/add-access":
			go l.RecodeOperationLog(userId, ip, "/v1/access/add-access")
			break
		case "/v1/role/add-role":
			go l.RecodeOperationLog(userId, ip, "/v1/role/add-role")
			break
		case "/v1/user-role/set-user-role":
			go l.RecodeOperationLog(userId, ip, "/v1/user-role/set-user-role")
			break
		case "/v1/role-access/set-role-access":
			go l.RecodeOperationLog(userId, ip, "/v1/role-access/set-role-access")
			break
		case "/v1/ga/verify":
			go l.RecodeOperationLog(userId, ip, "/v1/ga/verify")
			break
		case "/v1/ga/bind-verify":
			go l.RecodeOperationLog(userId, ip, "/v1/ga/bind-verify")
			break
		}
	}

	if method == "GET" {
		switch path {
		case "/v1/account/search":
			go l.RecodeOperationLog(userId, ip, "/v1/account/search")
			break
		case "/v1/account/user-info":
			go l.RecodeOperationLog(userId, ip, "/v1/account/user-info")
			break
		case "/v1/account/batch-user-by-ids":
			go l.RecodeOperationLog(userId, "ip 123", "/v1/account/batch-user-by-ids")
			break
		case "/v1/access/search":
			go l.RecodeOperationLog(userId, ip, "/v1/access/search")
			break
		case "/v1/access/search-user-pertain-access":
			go l.RecodeOperationLog(userId, ip, "/v1/access/search-user-pertain-access")
			break
		case "/v1/role/search":
			go l.RecodeOperationLog(userId, ip, "/v1/role/search")
			break
		case "/v1/user-role/search-user-role":
			go l.RecodeOperationLog(userId, ip, "/v1/user-role/search-user-role")
			break
		case "/v1/role-access/search":
			go l.RecodeOperationLog(userId, ip, "/v1/role-access/search")
			break
		case "/v1/ga/bind":
			go l.RecodeOperationLog(userId, ip, "/v1/ga/bind")
			break
		}
	}

	if method == "PUT" {
		switch path {
		case "/v1/account/update":
			go l.RecodeOperationLog(userId, ip, "/v1/account/update")
			break
		case "/v1/account/disabled":
			go l.RecodeOperationLog(userId, ip, "/v1/account/disabled")
			break
		case "/v1/account/set-admin":
			go l.RecodeOperationLog(userId, ip, "/v1/account/set-admin")
			break
		case "/v1/account/before-change-password":
			go l.RecodeOperationLog(userId, ip, "/v1/access/before-change-password")
			break
		case "/v1/account/after-change-password":
			go l.RecodeOperationLog(userId, ip, "/v1/account/after-change-password")
			break
		case "/v1/account/change-user-password":
			go l.RecodeOperationLog(userId, ip, "/v1/account/change-user-password")
			break
		case "/v1/access/update":
			go l.RecodeOperationLog(userId, ip, "/v1/access/update")
			break
		case "/v1/role/update":
			go l.RecodeOperationLog(userId, ip, "/v1/role/update")
			break
		case "/v1/role/disabled":
			go l.RecodeOperationLog(userId, ip, "/v1/role/disabled")
			break
		}
	}

	if method == "DELETE" {
		switch path {
		case "/v1/account/delete":
			go l.RecodeOperationLog(userId, ip, "/v1/account/delete")
			break
		case "/v1/account/logout":
			go l.RecodeOperationLog(userId, ip, "/v1/account/logout")
			break
		case "/v1/access/delete":
			go l.RecodeOperationLog(userId, ip, "/v1/access/delete")
			break
		case "/v1/role/delete":
			go l.RecodeOperationLog(userId, ip, "/v1/role/delete")
			break
		}
	}

	if method == "ANY" {
		switch path {
		case "/v1/upload/coinlogo":
			go l.RecodeOperationLog(userId, ip, "/v1/upload/coinlogo")
			break
		case "/v1/upload/notice":
			go l.RecodeOperationLog(userId, ip, "/v1/upload/notice")
			break
		case "/v1/upload/notify":
			go l.RecodeOperationLog(userId, ip, "/v1/upload/notify")
			break
		}
	}
}

func (l *LogController) RecodeLoginLog(userId uint, ip string, device string) {
	defer utils.PanicPrint()
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(l.config.IpFind.Auth, ip)
	if err != nil {
		l4g.Error("IpLocation err", err.Error())
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
		l4g.Error("CreateLoginLog err", err.Error())
		//		glog.Error(err.Error())
		return
	}
}

func (l *LogController) RecodeOperationLog(userId uint, ip string, operation string) {
	defer utils.PanicPrint()
	var (
		country string
		city    string
	)
	ipInfo, err := utils.IpLocation(l.config.IpFind.Auth, ip)
	if err != nil {
		l4g.Error("IpLocation err", err.Error())
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
		l4g.Error("CreateOperationLog err", err.Error())
		//		glog.Error(err.Error())
		return
	}
}

func (l *LogController) GetLoginLog(ctx iris.Context) {
	var (
		//userId uint
		err    error
		params struct {
			common.PageParams
			UserId uint `params:"user_id"`
		}

		respData []interface{}
	)

	common.GetParams(ctx, &params)
	common.GetParams(ctx, &params.PageParams)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	//userId = common.GetUserIdFromCtx(ctx)

	data, count, err := l.logModel.GetLoginLog(params.UserId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		l4g.Error("GetLoginLog err", err.Error())
		//		glog.Error(err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
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

	ctx.JSON(NewResponse(0, "").SetLimitResult(respData, count, params.Page))
}

func (l *LogController) GetOperationLog(ctx iris.Context) {
	var (
		//userId uint
		err    error
		params struct {
			common.PageParams
			UserId uint `params:"user_id"`
		}

		respData []interface{}
	)

	common.GetParams(ctx, &params)
	common.GetParams(ctx, &params.PageParams)

	if params.Limit <= 0 || params.Limit >= 100 {
		params.Limit = 10
	}

	if params.Page <= 1 {
		params.Page = 1
	}

	//userId = common.GetUserIdFromCtx(ctx)

	data, count, err := l.logModel.GetOperationLog(params.UserId, params.Limit, params.Limit*(params.Page-1))
	if err != nil {
		l4g.Error("GetOperationLog err", err.Error())
		//		glog.Error(err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
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

	ctx.JSON(NewResponse(0, "").SetLimitResult(respData, count, params.Page))
}
