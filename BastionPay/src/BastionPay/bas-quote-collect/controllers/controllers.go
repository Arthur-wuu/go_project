package controllers

import (
	"BastionPay/bas-quote-collect/common"
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
)

type (
	Controllers struct {
	}

	Response struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
		Data    interface{} `json:"data"`
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
