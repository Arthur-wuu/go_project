package controllers

import (
	"BastionPay/bas-wallet-im-assist/common"
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
)

type (
	Controllers struct {
	}

	Response struct {
		ActionStatus string `json:"ActionStatus"`
		ErrorCode    int    `json:"ErrorCode"`
		ErrorInfo    string `json:"ErrorInfo"`
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
	info string) {

	ctx.JSON(
		Response{
			ErrorCode:    0,
			ErrorInfo:    info,
			ActionStatus: "OK",
		})
}

func (c *Controllers) ExceptionSerive(
	ctx context.Context,
	message string) {

	ctx.JSON(
		Response{
			ErrorCode:    1,
			ErrorInfo:    message,
			ActionStatus: "FAIL",
		})
}
