package controllers

import (
	"BastionPay/bas-wallet-im-assist/db"
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-im-assist/api"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	UserStatus struct {
		Controllers
	}
)

func (this *UserStatus) Handle(ctx iris.Context) {

	param := new(api.UserStatus)

	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	accountName := param.UserStatusInfo.ToAccount
	accountStatu := param.UserStatusInfo.Action

	//存cache  key:用户名   value:状态 Logout/Login
	db.GCache.SetAccountCache(accountName, accountStatu)

	this.Response(ctx, "set account statu succ")
}
