package controllers

import (
	"encoding/json"
	"github.com/BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/models/redis"
	"github.com/BastionPay/bas-bkadmin-api/services/account"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"time"
)

type Login struct {
	Controllers
}

func (this *Login) Login(ctx iris.Context) {
	l4g.Debug("start deal Login username[%s]", utils.GetValueUserName(ctx))
	var userLogin account.UserLogin
	param := &account.Login{}

	userLogin = param
	err := ctx.ReadJSON(param)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err:%s", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)

	isBool, err := govalidator.ValidateStruct(param)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	var user *models.Account

	switch {
	case govalidator.IsEmail(param.Name):
		user, err = userLogin.GetUserInfoByEmail(param.Name)
		if err != nil {
			l4g.Error("GetUserInfoByEmail username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
			ctx.JSON(Response{Code: apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), Message: err.Error()})
			return
		}

		break
	case Tools.IsMobile(param.Name):
		user, err = userLogin.GetUserInfoByMobile(param.Name)
		if err != nil {
			l4g.Error("GetUserInfoByMobile username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
			ctx.JSON(Response{Code: apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), Message: err.Error()})
			return
		}

		break
	default:
		user, err = userLogin.GetUserInfoByName(param.Name)
		if err != nil {
			l4g.Error("GetUserInfoByName username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
			ctx.JSON(Response{Code: apibackend.BASERR_ACCOUNT_NOT_FOUND.Code(), Message: err.Error()})
			return
		}
	}

	if user.Status == 0 {
		l4g.Error("user.Status==0 username[%s] user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_BLOCK_ACCOUNT.Code(), Message: "user has been disabled"})
		return
	}

	err = this.checkLogin(user, param.Password)
	if err != nil {
		l4g.Error("checkLogin username[%s] user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INCORRECT_PWD.Code(), Message: err.Error()})
		return
	}

	user.Token = Tools.GenerateUserLoginToken(user.Id)
	body, err := json.Marshal(user)
	if err != nil {
		l4g.Error("Marshal username[%s] user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATA_PACK_ERROR.Code(), Message: err.Error()})
		return
	}

	_, err = redis.RedisClient.Set(user.Token, body, this.Config.System.Expire*time.Second).Result()
	if err != nil {
		l4g.Error("RedisClient username[%s]user[%v] err[%s]", utils.GetValueUserName(ctx), user, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	//loginLog := access.UserLoginLog{
	//	UserId: user.Id,
	//	Ip:     common.GetRealIp(ctx),
	//}
	//loginLog.SaveLog()
	go GLogController.RecodeLoginLog(uint(user.Id), common.GetRealIp(ctx), "web")

	switch true {
	case user.GoogleSecret != "":
		user.GoogleSecret = "1"
		break
	case user.GoogleSecret == "":
		user.GoogleSecret = "0"
		break
	}

	//处理管理员的密码，手机号码，邮箱
	user.Email = common.SecretEmail(user.Email)
	user.Mobile = common.SecretPhone(user.Mobile)
	user.Password = common.SecretPassword(user.Password)

	ctx.JSON(Response{Data: user})
	l4g.Debug("deal Login username[%s] ok, result[%v]", utils.GetValueUserName(ctx), user)
}

func (this *Login) Logout(ctx iris.Context) {
	l4g.Debug("start deal Logout username[%s]", utils.GetValueUserName(ctx))
	token := ctx.GetHeader("token")
	l4g.Debug("username[%s] token[%s]", utils.GetValueUserName(ctx), token)

	if err := redis.RedisClient.Del(token).Err(); err != nil {
		l4g.Error("RedisClient username[%s]token[%v] err[%s]", utils.GetValueUserName(ctx), token, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Message: errors.New("logout successfully").Error()})
	l4g.Debug("deal Logout username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *Login) checkLogin(data *models.Account, password string) error {
	if data == nil {
		return errors.New("User does not exist")
	}

	password = Tools.MD5(password)
	if !Tools.CheckPassword(password, data.Password) {
		return errors.New("Incorrect user password")
	}

	return nil
}
