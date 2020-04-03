package controllers

import (
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/common"
	"github.com/kataras/iris/context"
	//"gopkg.exa.center/blockshine-ex/api-article/config"
)

type (
	Controllers struct {
	}

	Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	PriceResponse struct {
		Code     int         `json:"code"`
		Message  string      `json:"message"`
		Data     interface{} `json:"data,omitempty"`
		UsdPrice interface{} `json:"usd_price,omitempty"`
	}
)

var (
	Tools *common.Tools
)

func init() {
	Tools = common.New()
}

func (c *Controllers) Response(ctx context.Context, data interface{}) {

	ctx.JSON(
		Response{
			Code:    0,
			Message: "Success",
			Data:    data,
		})
}

func (c *Controllers) RefundResponse(ctx context.Context, refundRes baspay.RefundRes) {
	ctx.JSON(
		Response{
			Code:    0,
			Message: "Success",
			Data:    refundRes.Data,
		})
}

func (c *Controllers) PriceResponse(ctx context.Context, data, UsdPrice interface{}) {

	ctx.JSON(
		PriceResponse{
			Code:     0,
			Message:  "Success",
			Data:     data,
			UsdPrice: UsdPrice,
		})
}

func (c *Controllers) PosResponse(ctx context.Context, posRes baspay.PosRes) {

	ctx.JSON(
		posRes)
}

func (c *Controllers) PosOrdersRes(ctx context.Context, posOrderRes baspay.PosOrdersRes) {

	ctx.JSON(
		posOrderRes)
}

func (c *Controllers) CoffeeResponse(ctx context.Context, coffeeStruct CoffeeStruct) {

	ctx.JSON(
		coffeeStruct)
}

func (c *Controllers) CoffeeStatusRes(ctx context.Context, coffeeStatus CoffeeStatus) {
	ctx.JSON(
		coffeeStatus)
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

//type coffeeStruct struct {
//	orderid  string  `json:"orderid,omitempty"`
//	torderid string  `json:"torderid,omitempty"`
//	code     string  `json:"code,omitempty"`
//	msg      string  `json:"msg,omitempty"`
//	twocode  string  `json:"twocode,omitempty"`
//}
