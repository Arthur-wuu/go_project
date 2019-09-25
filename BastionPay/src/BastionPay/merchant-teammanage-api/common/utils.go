package common

import (
	. "BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

const (
	CompanyName = "BASTIONPAY"
)

func GetRealIp(ctx iris.Context) string {
	if ctx == nil {
		return ""
	}
	var (
		ips string
	)
	ZapLog().With(zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP")),
		zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP")),
		zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP"))).
		Info("Get IP from header")
		//	glog.Info("Get IP from header X-Real-IP", ctx.GetHeader("X-Real-IP"))
		//	glog.Info("Get IP from header X-Forwarded-For", ctx.GetHeader("X-Forwarded-For"))
		//	glog.Info("Get IP from ctx.RemoteAddr()", ctx.RemoteAddr())

	ips = ctx.GetHeader("X-Forwarded-For")
	if ips == "" {
		ips = ctx.GetHeader("X-Real-IP")
	}
	if ips == "" {
		ips = ctx.RemoteAddr()
	}

	return strings.Split(ips, ",")[0]
}
