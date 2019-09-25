package controllers

import (
	"fmt"
	"github.com/kataras/iris"
)

type (
	Form struct {
		Controllers
	}
)

type FormParam struct {
	MerchantTradeNo string `form:"merchant_order_no"`
	MerchantId      string `form:"merchant_id"`
	Sign            string `form:"sign"`
	Amount          string `form:"amount"`
	Assets          string `form:"assets"`
	ProductName     string `form:"product_name"`
	ReturnUrl       string `form:"return_url"`
	ShowUrl         string `form:"merchant_id"`
}

func (this *Form) FormAction(ctx iris.Context) {

	//form := new(FormParam)
	//err:=Tools.ShouldBindQuery(ctx, form)
	ctx.Request().ParseForm()

	fmt.Println(ctx.Request().Form)

	//fmt.Println("**form id**",form.MerchantId)
	//fmt.Println("**form no**",form.MerchantTradeNo)

	mid := ctx.Request().Form.Get("merchant_id")
	mno := ctx.Request().Form.Get("merchant_order_no")

	fmt.Println("id:", mid)
	fmt.Println("no:", mno)

	ctx.Redirect("https://m.bastionpay.io/transitional/" + mid + "/" + mno)

	//_, err := base.HttpSend("https://test-m.bastionpay.io/transitional?merchant_id="+form.MerchantId+"&merchant_order_no="+form.MerchantTradeNo, nil, "GET",nil)
	//
	//if err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "direct to page err")
	//	ZapLog().Error( "direct to page err",zap.Error(err))
	//	return
	//}
	//this.Response(ctx, "succ")
	return
}
