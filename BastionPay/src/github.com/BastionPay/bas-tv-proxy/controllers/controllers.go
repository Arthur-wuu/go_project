package controllers

import (
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
	"BastionPay/bas-tv-proxy/api"
)

type Controllers struct {

}

func (c *Controllers) Response(
	ctx context.Context,
	data interface{}) {

	res := api.NewResponse(ctx.URLParam("qid"))
	res.Data = data

	ctx.JSON(res)
}

func (c *Controllers) ExceptionSerive(
	ctx context.Context,
	code int32) {
	res := api.NewErrResponse(ctx.URLParam("qid"), code)

	ctx.JSON(res)
}
