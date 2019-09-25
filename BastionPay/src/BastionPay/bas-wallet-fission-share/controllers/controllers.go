package controllers

import (
	"BastionPay/bas-wallet-fission-share/common"
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
)

type (
	Controllers struct {
	}

	Response struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message,omitempty"`
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
