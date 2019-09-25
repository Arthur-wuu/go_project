package controllers

import (
	"github.com/kataras/iris/context"
	"BastionPay/pay-user-merchant-api/common"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
	"github.com/kataras/iris/middleware/i18n"
	"BastionPay/pay-user-merchant-api/config"
)

type (
	Controllers struct {

	}

	OrderResponse struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	Response struct {
		Status struct {
			Code int `json:"code"`
			Message string `json:"msg"`
		} `json:"status"`
		Result interface{} `json:"result,omitempty"`
	}
)

func Init() error {
	phoneSmsLimiter =  common.NewBusLimiter(&common.GRedis, "buslimit_sms_phone_", config.GPreConfig.PhoneSmsLimits)
	ipSmsLimiter = common.NewBusLimiter(&common.GRedis, "buslimit_sms_ip_", config.GPreConfig.IpSmsLimits)
	phoneEmailLimiter = common.NewBusLimiter(&common.GRedis, "buslimit_email_phone_", config.GPreConfig.PhoneSmsLimits)
	ipEmailLimiter  = common.NewBusLimiter(&common.GRedis, "buslimit_email_ip_", config.GPreConfig.IpSmsLimits)

	if err := phoneSmsLimiter.Init(); err != nil {
		return err
	}
	if err := ipSmsLimiter.Init(); err != nil {
		return err
	}
	if err := phoneEmailLimiter.Init(); err != nil {
		return err
	}
	if err := ipEmailLimiter.Init(); err != nil {
		return err
	}
	return nil
}

var (
	Tools *common.Tools
)

func init() {
	Tools = common.New()
}

func (c *Controllers) Response(
	ctx context.Context,
	data interface{}) {
	res := &Response{
		Result:    data,
	}
	res.Status.Code = 0
	res.Status.Message = "Success"

	ctx.JSON(res)

}

func (c *Controllers) ExceptionSerive(
	ctx context.Context,
	code int,
	message string) {

	res := new(Response)
	res.Status.Code = code
	res.Status.Message = message
	ctx.JSON(res)
}

func (c *Controllers) ExceptionSeriveWithParams(
	ctx context.Context,
	code int,
	message string, params ...interface{}) {

	res := new(Response)
	res.Status.Code = code
	res.Status.Message = message

	res.Status.Message = i18n.Translate(ctx, message, params...)
	if res.Status.Message == "" {
		res.Status.Message = message
	}
	ctx.JSON(res)
}

func (r *Controllers) ResponsePage(ctx context.Context, result interface{}, total interface{}, page interface{})  {

	tmp := &struct {
		Total interface{} `json:"total"`
		Page  interface{} `json:"page"`
		Data  interface{} `json:"data"`
	}{total, page, result}

	res := &Response{
		Result:    tmp,
	}
	res.Status.Code = 0
	res.Status.Message = "Success"

	ctx.JSON(res)
}