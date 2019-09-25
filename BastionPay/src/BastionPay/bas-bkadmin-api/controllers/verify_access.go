package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/models/redis"
	"github.com/BastionPay/bas-bkadmin-api/services/account"
	"github.com/BastionPay/bas-bkadmin-api/services/rbac"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"net/http"
	"strings"
	"time"
)

type (
	Verify struct {
		Controllers
	}
)

func (this *Verify) VerifyAccess(ctx iris.Context) {
	l4g.Debug("start deal VerifyAccess  path[%s]token[%s]", ctx.Path(), ctx.GetHeader("token"))
	var user *models.Account
	var err error
	var token string
	if this.needToken(ctx.Path()) {
		token = ctx.GetHeader("token")
		user, err = new(account.Account).GetUserInfoByToken(token)
		if err == nil || user != nil {
			ctx.Values().Set(token, user)
		}
	}
	l4g.Debug("username[%s] user[%v] path[%s]", utils.GetValueUserName(ctx), user, ctx.Path())

	if user != nil &&
		user.GoogleSecret != "" &&
		!user.IsGauth && !this.checkGA(ctx.Path()) {
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: "need GA"})
		//		this.Status(ctx, http.StatusPreconditionFailed)
		l4g.Error("StatusPreconditionFailed username[%s] user[%v] path[%s] token[%s] err", utils.GetValueUserName(ctx), user, ctx.Path(), ctx.GetHeader("token"))
		return
	}

	//是否免检uri
	verifyAccess := rbac.VerifyAccess{}
	if verifyAccess.GetIgnoreUri().Ignore(ctx.Path()) {
		ctx.Next()
		l4g.Debug("deal VerifyAccess username[%s] ok", utils.GetValueUserName(ctx))
		return
	}

	if token == "" || user == nil {
		ctx.JSON(Response{Code: apibackend.BASERR_TOKEN_INVALID.Code(), Message: "no token or user"})
		//this.Status(ctx, http.StatusUnauthorized)
		l4g.Error("StatusUnauthorized username[%s] token[%s] user[%v] err", utils.GetValueUserName(ctx), token, user)
		return
	}

	redis.RedisClient.Expire(ctx.GetHeader("token"), this.Config.System.Expire*time.Second)

	verifyAccess.UserId = user.Id

	//是否是超级管理员
	if user.IsAdmin == 1 {
		ctx.Next()
		l4g.Debug("deal VerifyAccess username[%s] ok", utils.GetValueUserName(ctx))
		return
	}

	//获取用户权限列表 ,下面两步可能有bug，权限出错，但是无任何错误日志
	err = verifyAccess.GetUserAccessList()
	if err != nil {
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		//		this.Status(ctx, http.StatusForbidden)
		l4g.Error("StatusForbidden username[%s] verifyAccess[%v] err[%s]", utils.GetValueUserName(ctx), verifyAccess, err.Error())
		return
	}

	ctx.Values().Set("user_access_list", verifyAccess.UserAccessList)
	ok := verifyAccess.GetUserList(user.Id, ctx.Path())
	if ok {
		ctx.Next()
		l4g.Debug("deal VerifyAccess username[%s] ok", utils.GetValueUserName(ctx))
		return
	}

	ctx.JSON(Response{Code: apibackend.BASERR_UNAUTHORIZED_METHOD.Code(), Message: "have no permission"})
	//this.Status(ctx, http.StatusForbidden)
	l4g.Error("StatusForbidden username[%s] verifyAccess[%v] user.Id[%s] path[%s] err[GetUserList not ok]", utils.GetValueUserName(ctx), verifyAccess, user.Id, ctx.Path())
	return
}

func (this *Verify) Status(ctx iris.Context, code int) {
	//		ctx.StatusCode(code)
	ctx.JSON(Response{
		Code:    code,
		Message: http.StatusText(code),
	})
}

func (this *Verify) checkGA(path string) bool {
	gaPath := []string{
		"/v1/ga/bind-verify",
		"/v1/ga/verify",
	}

	for _, v := range gaPath {
		if strings.EqualFold(v, path) {
			return true
		}
	}

	return false
}

func (this *Verify) needToken(path string) bool {
	if path == "/v1/account/login" {
		return false
	}
	if path == "/v1/bastionpay/admin/sp_get_asset_attribute" {
		return false
	}
	if path == "/v1/bastionpay/admin/support_assets" {
		return false
	}
	return true
}
