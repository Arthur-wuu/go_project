package controllers

import (
	"BastionPay/merchant-workattence-api/common"
	"github.com/kataras/iris/context"
)

type (
	Controllers struct {
	}

	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	ResponsePush struct {
		Status int         `json:"status"`
		Info   string      `json:"info"`
		Data   interface{} `json:"data,omitempty"`
	}
)

var (
	Tools *common.Tools
)

func init() {
	Tools = common.New()
}

func (c *Controllers) Response(
	ctx context.Context,
	data interface{}) {

	ctx.JSON(
		Response{
			Code:    0,
			Message: "Success",
			Data:    data,
		})
}

func (c *Controllers) ExceptionSerive(
	ctx context.Context,
	code int,
	message string) {

	ctx.JSON(
		Response{
			Code:    code,
			Message: message,
		})
}

func (c *Controllers) ExceptionSeriveWithData(
	ctx context.Context,
	code int,
	data interface{}) {

	ctx.JSON(
		Response{
			Code: code,
			Data: data,
		})
}

func (c *Controllers) ResponsePush(ctx context.Context, data interface{}) {
	ctx.JSON(ResponsePush{
		Status: 1,
		Info:   "ok",
		Data:   data,
	})
}
