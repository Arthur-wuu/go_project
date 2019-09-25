package controllers

import (
	"BastionPay/bas-api/apibackend"
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/models/redis"
	"github.com/BastionPay/bas-bkadmin-api/services/account"
	"github.com/BastionPay/bas-bkadmin-api/services/ga"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"time"
)

type (
	GA struct {
		Controllers
	}
)

func (this *GA) Bind(ctx iris.Context) {
	l4g.Debug("start deal Bind username[%s]", utils.GetValueUserName(ctx))
	user, err := User.GetValueUserInfo(ctx)
	if err != nil {
		l4g.Error("GetValueUserInfo username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] user[%v]", utils.GetValueUserName(ctx), user)

	if user.GoogleSecret != "" {
		l4g.Error("GoogleSecret username[%s] err[is nil]", utils.GetValueUserName(ctx))
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: "already bound"})
		return
	}

	gas := ga.NewGA(this.Config.System.CompanyName)

	gas.Generate(user.Email)

	_, err = redis.RedisClient.Set(this.generate(user.Id), gas.Secret, this.Config.System.GaExpire*time.Second).Result()
	if err != nil {
		l4g.Error("RedisClient username[%s]user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: gas})
	l4g.Debug("deal Bind username[%s] ok, result[%v]", utils.GetValueUserName(ctx), gas)
	ctx.Next()
}

func (this *GA) BindVerify(ctx iris.Context) {
	l4g.Debug("start deal BindVerify username[%s]", utils.GetValueUserName(ctx))

	var param ga.GA

	err := Tools.ShouldBindJSON(ctx, &param)
	if err != nil {
		l4g.Error("ShouldBindJSON err[%s]", err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)
	user, err := User.GetValueUserInfo(ctx)
	if err != nil {
		l4g.Error("GetValueUserInfo username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: err.Error()})
		return
	}

	user.GoogleSecret, err = redis.RedisClient.Get(this.generate(user.Id)).Result()
	if err != nil {
		l4g.Error("RedisClient username[%s]user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "verification expired"})
		return
	}

	ok, err := param.Verify(user.GoogleSecret, param.Code)
	if err != nil || !ok {
		l4g.Error("Verify username[%s] GoogleSecret[%s] Code[%s] param[%v] err[%v]", utils.GetValueUserName(ctx), user.GoogleSecret, param.Code, param, err)
		ctx.JSON(Response{Code: apibackend.BASERR_INCORRECT_GA_PWD.Code(), Message: "google verify fail"})
		return
	}

	bindUser := account.AccountUpdate{
		Id:           user.Id,
		GoogleSecret: user.GoogleSecret,
	}

	_, err = bindUser.UpdateUserInfo()
	if err != nil {
		l4g.Error("UpdateUserInfo username[%s] bindUser[%v] err[%s]", utils.GetValueUserName(ctx), bindUser, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "update user info fail"})
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		l4g.Error("Marshal username[%s] user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	err = redis.RedisClient.Set(ctx.GetHeader("token"), body, this.Config.System.Expire*time.Second).Err()
	if err != nil {
		l4g.Error("RedisClient username[%s]token[%v] err[%s]", utils.GetValueUserName(ctx), ctx.GetHeader("token"), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: "set cache user info fail"})
		return
	}

	ctx.JSON(Response{Message: "verification ok", Data: true})
	l4g.Debug("deal BindVerify username[%s] ok, result[%v]", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *GA) generate(userId int64) string {

	return fmt.Sprintf("googauth_%d", userId)
}

func (this *GA) Verify(ctx iris.Context) {
	l4g.Debug("start deal Verify username[%s]", utils.GetValueUserName(ctx))
	var param ga.GA

	err := Tools.ShouldBindJSON(ctx, &param)
	if err != nil {
		l4g.Error("ShouldBindJSON err[%s]", err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)

	user, err := User.GetValueUserInfo(ctx)
	if err != nil || user.GoogleSecret == "" {
		l4g.Error("GetValueUserInfo username[%s] param[%v] GoogleSecret[%s] err[%s]", utils.GetValueUserName(ctx), param, user.GoogleSecret, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_DATA_NOT_FOUND.Code(), Message: "get user info null"})
		return
	}

	ok, err := param.Verify(user.GoogleSecret, param.Code)
	if err != nil || !ok {
		l4g.Error("Verify username[%s] GoogleSecret[%s] Code[%s] param[%v] err[%v]", utils.GetValueUserName(ctx), user.GoogleSecret, param.Code, param, err)
		ctx.JSON(Response{Code: apibackend.BASERR_INCORRECT_GA_PWD.Code(), Message: "google verify fail"})
		return
	}

	user.IsGauth = true
	body, err := json.Marshal(user)
	if err != nil {
		l4g.Error("Marshal username[%s] user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	err = redis.RedisClient.Set(ctx.GetHeader("token"), body, this.Config.System.Expire*time.Second).Err()
	if err != nil {
		l4g.Error("RedisClient username[%s]token[%v] err[%s]", utils.GetValueUserName(ctx), ctx.GetHeader("token"), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "set cache user info fail"})
		return
	}

	ctx.JSON(Response{Message: "verification ok", Data: true})
	l4g.Debug("deal Verify username[%s] ok, result[%v]", utils.GetValueUserName(ctx))
}
