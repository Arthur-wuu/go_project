package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/services/rbac"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
)

type (
	UserRole struct {
		Controllers
	}
)

func (this *UserRole) SetUserRule(ctx iris.Context) {
	l4g.Debug("start deal SetUserRule username[%s]", utils.GetValueUserName(ctx))
	param := rbac.UserRole{}

	err := ctx.ReadJSON(&param)
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

	err = param.SetUserRule()
	if err != nil {
		l4g.Error("SetUserRule username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_NOT_FOUND.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Message: "set user role success"})
	l4g.Debug("deal SetUserRule username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *UserRole) SearchUserRole(ctx iris.Context) {
	l4g.Debug("start deal SearchUserRole username[%s]", utils.GetValueUserName(ctx))
	param := &rbac.UserRoleList{
		UserId: Tools.ParseInt(ctx.FormValue("user_id"), 0),
		Page:   Tools.ParseInt(ctx.FormValue("page"), 1),
		Size:   Tools.ParseInt(ctx.FormValue("size"), SEARCHSIZE),
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)

	isBool, err := govalidator.ValidateStruct(param)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	result, err := param.SearchUserRole()
	if err != nil {
		l4g.Error("SearchUserRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal SearchUserRole username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
}
