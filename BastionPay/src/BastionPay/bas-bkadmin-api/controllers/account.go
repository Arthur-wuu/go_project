package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/services/account"
	"github.com/BastionPay/bas-bkadmin-api/services/ga"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
)

type Account struct {
	Controllers
}

func (this *Account) Register(ctx iris.Context) {
	l4g.Debug("start deal Register")

	param := &account.Account{}

	err := ctx.ReadJSON(&param)
	if err != nil {
		l4g.Error("Context ReadJSON err[%s]", err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", param.Name, param)

	isBool, err := govalidator.ValidateStruct(param)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err[%s]", param.Name, param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	log := new(account.Login)
	if param.Name != "" {
		user, _ := log.GetUserInfoByName(param.Name)
		if user.Id > 0 {
			l4g.Error("GetUserInfoByName userId[%d]param[%v] is exist", user.Id, param)
			ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: "name user is exist"})
			return
		}
	}

	if param.Mobile != "" {
		user, _ := log.GetUserInfoByMobile(param.Mobile)
		if user.Id > 0 {
			l4g.Error("GetUserInfoByMobile userId[%d]param[%v] is exist", user.Id, param)
			ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: "mobile user is exist"})
			return
		}
	}

	if param.Email != "" {
		user, _ := log.GetUserInfoByEmail(param.Email)
		if user.Id > 0 {
			l4g.Error("GetUserInfoByEmail userId[%d]param[%v] is exist", user.Id, param)
			ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: "email user is exist"})
			return
		}
	}

	result, err := param.Register(param)
	if err != nil {
		l4g.Error("Register username[%s] param[%v] err[%s]", param.Name, param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal Register username[%s] ok, result[%v]", param.Name, result)
	ctx.Next()
}

func (this *Account) Search(ctx iris.Context) {
	l4g.Debug("start deal Search username[%s]", utils.GetValueUserName(ctx))

	var params account.AccountList

	err := common.GetParams(ctx, &params)
	if err != nil {
		l4g.Error("GetParams err[%s]", err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), params)

	isBool, err := govalidator.ValidateStruct(params)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	result, err := params.Search()
	if err != nil {
		l4g.Error("Search username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal Search username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
}

func (this *Account) Update(ctx iris.Context) {
	l4g.Debug("start deal Update username[%s]", utils.GetValueUserName(ctx))

	var accounts account.AccountUpdate

	err := ctx.ReadJSON(&accounts)
	if err != nil {
		l4g.Error("Context ReadJSON username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	isBool, err := govalidator.ValidateStruct(accounts)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err:%s", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	user, err := accounts.UpdateUserInfo()
	if err != nil {
		l4g.Error("UpdateUserInfo username[%s] accounts[%v] err[%s]", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	userValue, err := User.GetValueUserInfo(ctx)
	if err == nil {
		if userValue.Id != this.UserId {
			account.RemoveToken(Tools.GenerateUserLoginToken(this.UserId))
		}
	}
	user.Email = common.SecretEmail(user.Email)
	user.Mobile = common.SecretPhone(user.Mobile)
	user.Password = common.SecretPassword(user.Password)

	ctx.JSON(Response{Data: user})
	l4g.Debug("deal Update username[%s] ok, result[%v]", utils.GetValueUserName(ctx), user)
	ctx.Next()
}

func (this *Account) SetAdmin(ctx iris.Context) {
	l4g.Debug("start deal SetAdmin username[%s]", utils.GetValueUserName(ctx))
	var accounts account.AccountUpdate

	err := ctx.ReadJSON(&accounts)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err:%s", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] accounts[%v]", utils.GetValueUserName(ctx), accounts)

	isBool, err := govalidator.ValidateStruct(accounts)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] accounts[%v] err:%s", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	user, err := accounts.SetUserAdmin(ctx.GetHeader("token"))
	if err != nil {
		l4g.Error("SetUserAdmin token[%s] err[%s]", ctx.GetHeader("token"), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	switch true {
	case user.GoogleSecret != "":
		user.GoogleSecret = "1"
		break
	case user.GoogleSecret == "":
		user.GoogleSecret = "0"
		break
	}

	user.Email = common.SecretEmail(user.Email)
	user.Mobile = common.SecretPhone(user.Mobile)
	user.Password = common.SecretPassword(user.Password)

	ctx.JSON(Response{Data: user})
	l4g.Debug("deal SetAdmin username[%s] ok, result[%v]", utils.GetValueUserName(ctx), user)
	ctx.Next()
}

func (this *Account) Disabled(ctx iris.Context) {
	l4g.Debug("start deal Disabled username[%s]", utils.GetValueUserName(ctx))
	var accounts account.AccountUpdate

	err := ctx.ReadJSON(&accounts)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] accounts[%v]", utils.GetValueUserName(ctx), accounts)
	isBool, err := govalidator.ValidateStruct(accounts)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] accounts[%v] err[%s]", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	user, err := accounts.DisabledUser(Tools.GenerateUserLoginToken(accounts.Id))
	if err != nil {
		l4g.Error("DisabledUser username[%s] accounts.Id[%d] err[%s]", utils.GetValueUserName(ctx), accounts.Id, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	switch true {
	case user.GoogleSecret != "":
		user.GoogleSecret = "1"
		break
	case user.GoogleSecret == "":
		user.GoogleSecret = "0"
		break
	}

	user.Email = common.SecretEmail(user.Email)
	user.Mobile = common.SecretPhone(user.Mobile)
	user.Password = common.SecretPassword(user.Password)

	ctx.JSON(Response{Data: user})
	l4g.Debug("deal Disabled username[%s] ok, result[%v]", utils.GetValueUserName(ctx), user)
	ctx.Next()
}

func (this *Account) ChangeUserPassword(ctx iris.Context) {
	l4g.Debug("start deal ChangeUserPassword username[%s]", utils.GetValueUserName(ctx))
	pwd := account.ChangeUserPassword{}
	err := ctx.ReadJSON(&pwd)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] pwd[%v]", utils.GetValueUserName(ctx), pwd)
	isBool, err := govalidator.ValidateStruct(pwd)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] pwd[%v] err[%s]", utils.GetValueUserName(ctx), pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	if err = pwd.ChangeUserPasswords(); err != nil {
		l4g.Error("ChangeUserPasswords pwd[%v] err[%s]", pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_DATA_NOT_SAME.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Message: "change password success"})
	l4g.Debug("deal ChangeUserPassword username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *Account) ChangeBeforePassword(ctx iris.Context) {
	l4g.Debug("start deal ChangeBeforePassword username[%s]", utils.GetValueUserName(ctx))
	pwd := account.ChangeBeforePassword{}
	err := ctx.ReadJSON(&pwd)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] pwd[%v]", utils.GetValueUserName(ctx), pwd)

	isBool, err := govalidator.ValidateStruct(pwd)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] pwd[%v] err[%s]", utils.GetValueUserName(ctx), pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	userId := User.GetValueUserId(ctx)
	if userId <= 0 {
		l4g.Error("GetValueUserId pwd[%v] userId[%d] err", pwd, userId)
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: err.Error()})
		return
	}

	pwd.Id = userId
	uuid, err := pwd.ChangeBeforePassword()
	if err != nil {
		l4g.Error("ChangeBeforePassword pwd[%v] err[%s]", pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INCORRECT_PWD.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{
		Data:    uuid,
		Message: "success",
	})
	l4g.Debug("deal ChangeBeforePassword username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *Account) Delete(ctx iris.Context) {
	l4g.Debug("start deal Delete username[%s]", utils.GetValueUserName(ctx))

	var accounts account.AccountDel

	err := ctx.ReadJSON(&accounts)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] accounts[%v]", utils.GetValueUserName(ctx), accounts)

	isBool, err := govalidator.ValidateStruct(accounts)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] accounts[%v] err[%s]", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	user, err := accounts.Delete()
	if err != nil {
		l4g.Error("Delete username[%s] accounts[%v] err[%v]", utils.GetValueUserName(ctx), accounts, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	switch true {
	case user.GoogleSecret != "":
		user.GoogleSecret = "1"
		break
	case user.GoogleSecret == "":
		user.GoogleSecret = "0"
		break
	}
	//处理管理员邮箱，手机号码和密码
	user.Email = common.SecretEmail(user.Email)
	user.Mobile = common.SecretPhone(user.Mobile)
	user.Password = common.SecretPassword(user.Password)

	ctx.JSON(Response{Data: user})
	l4g.Debug("deal Delete username[%s] ok, user[%v]", utils.GetValueUserName(ctx), user)
	ctx.Next()
}

func (this *Account) BatchUserByIds(ctx iris.Context) {
	l4g.Debug("start deal BatchUserByIds username[%s]", utils.GetValueUserName(ctx))
	params := account.BatchUserByIds{
		Id:   ctx.FormValue("id"),
		Page: Tools.ParseInt(ctx.FormValue("page"), 1),
		Size: Tools.ParseInt(ctx.FormValue("size"), 100),
	}
	l4g.Debug("username[%s] params[%v]", utils.GetValueUserName(ctx), params)
	isBool, err := govalidator.ValidateStruct(params)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] params[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	result, err := params.BatchUserByIds()
	if err != nil {
		l4g.Error("BatchUserByIds username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal BatchUserByIds username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
	ctx.Next()
}

func (this *Account) GetUserInfo(ctx iris.Context) {
	l4g.Debug("start deal GetUserInfo username[%s]", utils.GetValueUserName(ctx))
	params := account.UserInfo{
		Id: Tools.ParseInt(ctx.FormValue("id"), 0),
	}
	l4g.Debug("username[%s] params[%v]", utils.GetValueUserName(ctx), params)
	isBool, err := govalidator.ValidateStruct(params)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] params[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	result, err := params.GetUserInfo()
	if err != nil {
		l4g.Error("GetUserInfo username[%s] params[%v] err[%s]", utils.GetValueUserName(ctx), params, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	switch true {
	case result.GoogleSecret != "":
		result.GoogleSecret = "1"
		break
	case result.GoogleSecret == "":
		result.GoogleSecret = "0"
		break
	}
	ctx.JSON(Response{Data: result})
	l4g.Debug("deal GetUserInfo username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
	ctx.Next()
}

func (this *Account) ChangeAfterPassword(ctx iris.Context) {
	l4g.Debug("start deal ChangeAfterPassword username[%s]", utils.GetValueUserName(ctx))
	pwd := account.ChangeAfterPassword{}
	err := ctx.ReadJSON(&pwd)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] pwd[%v]", utils.GetValueUserName(ctx), pwd)

	isBool, err := govalidator.ValidateStruct(pwd)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] pwd[%v] err[%s]", utils.GetValueUserName(ctx), pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	user, err := User.GetValueUserInfo(ctx)
	if err != nil || user.Id <= 0 {
		l4g.Error("GetValueUserInfo username[%s] pwd[%v] err[%s]", utils.GetValueUserName(ctx), pwd, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: err.Error()})
		return
	}

	ok, err := new(ga.GA).Verify(user.GoogleSecret, pwd.Code)
	if err != nil || !ok {
		l4g.Error("Verify username[%s] GoogleSecret[%s] Code[%s] err[%v]", utils.GetValueUserName(ctx), user.GoogleSecret, pwd.Code, err)
		ctx.JSON(Response{Code: apibackend.BASERR_INCORRECT_GA_PWD.Code(), Message: "google verify fail"})
		return
	}

	pwd.Id = user.Id
	isOk := pwd.ChangeAfterPassword()
	if isOk == false {
		l4g.Error("ChangeAfterPassword username[%s] pwd[%v] err", utils.GetValueUserName(ctx), pwd)
		ctx.JSON(Response{Code: apibackend.BASERR_SERVICE_UNKNOWN_ERROR.Code(), Message: "change password fail"})
		return
	}

	ctx.JSON(Response{
		Message: "change password success",
		Data:    true,
	})
	l4g.Debug("deal ChangeAfterPassword username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}
