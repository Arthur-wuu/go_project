package controllers

import (
	"BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/rbac"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
)

type (
	Role struct {
		Controllers
	}
)

func (this *Role) AddRule(ctx iris.Context) {
	l4g.Debug("start deal AddRule username[%s]", utils.GetValueUserName(ctx))
	param := rbac.Role{}

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

	result, err := param.AddRole()
	if err != nil {
		l4g.Error("AddRole username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_OBJECT_EXISTS.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal AddRule username[%s] ok, result=%v", utils.GetValueUserName(ctx), result)
	ctx.Next()
}

func (this *Role) Search(ctx iris.Context) {
	l4g.Debug("start deal Search username[%s]", utils.GetValueUserName(ctx))

	param := &rbac.RoleList{
		Page: Tools.ParseInt(ctx.FormValue("page"), 1),
		Size: Tools.ParseInt(ctx.FormValue("size"), SEARCHSIZE),
	}

	param.Name = ctx.FormValue("name")
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)

	result, err := param.RoleList()
	if err != nil {
		l4g.Error("RoleList username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal Search username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
}

func (this *Role) Delete(ctx iris.Context) {
	l4g.Debug("start deal Delete username[%s]", utils.GetValueUserName(ctx))

	param := rbac.RoleDelete{}

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

	tx := models.DB.Begin()
	result, err := param.Delete()
	if err == nil {
		err = new(rbac.RoleAccess).DeleteAccessById(result.Id)
		if err == nil {
			err = new(rbac.UserRole).DeleteRole(result.Id)
			if err == nil {
				tx.Commit()
				ctx.JSON(Response{Data: result})
				l4g.Debug("deal Delete username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
				return
			} else {
				l4g.Error("DeleteRole username[%s] id[%d] err[%s]", utils.GetValueUserName(ctx), result.Id, err.Error())
			}
		} else {
			l4g.Error("DeleteAccessById username[%s] id[%d] err[%s]", utils.GetValueUserName(ctx), result.Id, err.Error())
		}
	} else {
		l4g.Error("Delete username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
	}

	tx.Rollback()
	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
	ctx.Next()
}

func (this *Role) Update(ctx iris.Context) {
	l4g.Debug("start deal Update username[%s]", utils.GetValueUserName(ctx))

	param := rbac.RoleUpdate{}

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

	result, err := param.Update()
	if err == nil {
		l4g.Debug("deal Update username[%s] ok, result[%v]", utils.GetValueUserName(ctx), result)
		ctx.JSON(Response{Data: result})
		return
	} else {
		l4g.Error("Update username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
	}

	models.DB.Rollback()
	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
	ctx.Next()
}

func (this *Role) Disabled(ctx iris.Context) {
	l4g.Debug("start deal Disabled username[%s]", utils.GetValueUserName(ctx))

	param := rbac.RoleUpdate{}

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

	result, err := param.Disabled()
	if err == nil {
		ctx.JSON(Response{Data: result})
		l4g.Debug("deal Disabled username[%s] ok, result=%v", utils.GetValueUserName(ctx), result)
		return
	} else {
		l4g.Error("Disabled username[%s] param[%v] err[%s]", utils.GetValueUserName(ctx), param, err.Error())
	}

	models.DB.Rollback()
	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
	ctx.Next()
}
