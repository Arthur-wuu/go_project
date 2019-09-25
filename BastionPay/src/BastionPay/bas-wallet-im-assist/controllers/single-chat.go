package controllers

import (
	"fmt"
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-im-assist/api"
	"BastionPay/bas-wallet-im-assist/comsumer"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	SingleChat struct {
		Controllers
	}
)

func (this *SingleChat) Handle(ctx iris.Context) {

	param := new(api.SingleChat)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	fmt.Println("param", param)

	comsumer.GTasker.SendSChat(param)

	this.Response(ctx, "handle single chat succ")
}
