package controllers

import (
	"github.com/kataras/iris"
	//"BastionPay/bas-api/apibackend"
	//. "BastionPay/bas-base/log/zap"
	//"go.uber.org/zap"
)

type (
	CallBacker struct {
		Controllers
		mUserStatus UserStatus
		mSingleChat SingleChat
	}
)

func (this *CallBacker) Handle(ctx iris.Context) {
	cmd := ctx.URLParam("CallbackCommand")

	switch cmd {
	case "State.StateChange":
		this.mUserStatus.Handle(ctx)
	case "C2C.CallbackAfterSendMsg":
		this.mSingleChat.Handle(ctx)
	default:
		return
	}
}
