package controllers

import (
	"BastionPay/merchant-teammanage-api/common"
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
