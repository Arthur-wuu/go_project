package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/base"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Form struct {
		Controllers
	}
)

type FormParam struct {
	MerchantTradeNo string `form:"merchant_trade_no"`
	MerchantId      string `form:"merchant_id"`
}

//这个暂时没用，处理form的放在加解密里面了
func (this *Form) FormAction(ctx iris.Context) {

	form := new(FormParam)
	Tools.ShouldBindQuery(ctx, form)

	fmt.Println("form id", form.MerchantId)
	fmt.Println("form no", form.MerchantTradeNo)

	_, err := base.HttpSend("https://test-m.bastionpay.io/transitional?merchant_id="+form.MerchantId+"&merchant_order_no="+form.MerchantTradeNo, nil, "GET", nil)

	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "direct to page err")
		ZapLog().Error("direct to page err", zap.Error(err))
		return
	}
	this.Response(ctx, "succ")
}
