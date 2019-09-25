package controllers

import (
	"github.com/BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/services/access"
	"github.com/BastionPay/bas-bkadmin-api/services/rbac"
	"github.com/BastionPay/bas-bkadmin-api/utils"
	l4g "github.com/alecthomas/log4go"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
)

type Access struct {
	Controllers
}

func (this *Access) AddAccess(ctx iris.Context) {
	l4g.Debug("start deal AddAccess username[%s]", utils.GetValueUserName(ctx))
	var param *access.Access

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

	result, err := param.AddAccess()
	if err != nil {
		l4g.Error("AddAccess username[%s] param[%v] err:%s", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal AddAccess username[%s] ok, result=%v", utils.GetValueUserName(ctx), result)
	ctx.Next()
}

func (this *Access) Search(ctx iris.Context) {
	l4g.Debug("start deal Search username[%s]", utils.GetValueUserName(ctx))
	param := &access.AccessList{
		Name:     ctx.FormValue("name"),
		ParentId: Tools.ParseInt(ctx.FormValue("parent_id"), 0),
		Page:     Tools.ParseInt(ctx.FormValue("page"), 1),
		Size:     Tools.ParseInt(ctx.FormValue("size"), SEARCHSIZE),
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)
	isBool, err := govalidator.ValidateStruct(param)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err:%s", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	result, err := param.Search()
	if err != nil {
		l4g.Error("Search username[%s] body[%v] err:%s", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		return
	}

	ctx.JSON(Response{Data: result})
	l4g.Debug("deal Search username[%s] ok, result=%v", utils.GetValueUserName(ctx), result)
}

func (this *Access) Delete(ctx iris.Context) {
	l4g.Debug("start deal Delete username[%s]", utils.GetValueUserName(ctx))

	var param *access.AccessUpdateList

	err := ctx.ReadJSON(&param)
	if err != nil {
		l4g.Error("Context[%s] ReadJSON err:%s", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)

	isBool, err := govalidator.ValidateStruct(param)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] param[%v] err:%s", utils.GetValueUserName(ctx), param, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	//标识为删除
	param.Valid = 1
	param.ConnBegin = models.DB.Begin()

	ok := param.UpdateAccess()
	if ok {
		err = new(rbac.RoleAccess).DeleteAccessById(param.Id)
		if err == nil {
			param.ConnBegin.Commit()
			ctx.JSON(Response{Message: "update success"})
			l4g.Debug("deal Delete username[%s] ok", utils.GetValueUserName(ctx))
			return
		} else {
			l4g.Error("DeleteAccessById[%d] username[%s] err[%s]", param.Id, utils.GetValueUserName(ctx), err.Error())
		}
	} else {
		l4g.Error("UpdateAccess param[%v] username[%s] err", param, utils.GetValueUserName(ctx))
	}

	param.ConnBegin.Rollback()
	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "update failure"})
	ctx.Next()
}

func (this *Access) Update(ctx iris.Context) {
	l4g.Debug("start deal Update username[%s]", utils.GetValueUserName(ctx))
	param := access.AccessUpdateList{}

	err := Tools.ShouldBindJSON(ctx, &param)
	if err != nil {
		l4g.Error("ShouldBindJSON username[%s] err[%s]", utils.GetValueUserName(ctx), err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}
	l4g.Debug("username[%s] param[%v]", utils.GetValueUserName(ctx), param)
	if param.Id == param.ParentId {
		l4g.Error("Id[%d] same as ParentId[%d] username[%s] err", param.Id, param.ParentId, utils.GetValueUserName(ctx))
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: errors.New("id and parent_id cannot be the same")})
		return
	}

	param.ConnBegin = models.DB
	ok := param.UpdateAccess()
	if !ok {
		l4g.Error("UpdateAccess param[%v] username[%s] err", param, utils.GetValueUserName(ctx))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: "update failure"})
		return
	}

	ctx.JSON(Response{Message: "update success"})
	l4g.Debug("deal Update username[%s] ok", utils.GetValueUserName(ctx))
	ctx.Next()
}

func (this *Access) SearchUserPertainAccess(ctx iris.Context) {
	l4g.Debug("start deal SearchUserPertainAccess username[%s]", utils.GetValueUserName(ctx))

	userId := Tools.ParseInt(ctx.FormValue("user_id"), 0)

	userPertainAccess := rbac.UserPertainAccess{
		UserId: userId,
	}
	l4g.Debug("username[%s] userPertainAccess[%v]", utils.GetValueUserName(ctx), userPertainAccess)

	isBool, err := govalidator.ValidateStruct(userPertainAccess)
	if err != nil && !isBool {
		l4g.Error("ValidateStruct username[%s] userPertainAccess[%v] err:%s", utils.GetValueUserName(ctx), userPertainAccess, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_INVALID_PARAMETER.Code(), Message: err.Error()})
		return
	}

	verify := rbac.VerifyAccess{}
	user, err := User.GetValueUserInfo(ctx)
	if err != nil {
		l4g.Error("ValidateStruct username[%s] userPertainAccess[%v] err:%s", utils.GetValueUserName(ctx), userPertainAccess, err.Error())
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: err.Error()})
		return
	}

	if user.IsAdmin == 1 {
		accessResult := access.AccessList{
			Size: 10000,
		}

		result, err := accessResult.Search()
		if err == nil {
			l4g.Error("deal SearchUserPertainAccess username[%s] ok, result[%v] ", utils.GetValueUserName(ctx), result)
			ctx.JSON(Response{Data: result.List})
			return
		}
	}

	if userId != user.Id {
		l4g.Error("no right to operate other users userId[%d] != user.Id[%d] username[%s] err", userId, user.Id, utils.GetValueUserName(ctx))
		ctx.JSON(Response{Code: apibackend.BASERR_UNKNOWN_BUG.Code(), Message: "no right to operate other users"})
		return
	}

	var userAccess []*models.Access
	body := ctx.Values().Get("user_access_list")
	if body == nil {
		verify.UserId = userId
		verify.GetUserAccessList()
		userAccess = verify.UserAccessList
	} else {
		userAccess = body.([]*models.Access)
	}

	ctx.JSON(Response{Data: verify.SearchUserPertainAccess(userAccess)})
	l4g.Debug("deal SearchUserPertainAccess username[%s] ok ", utils.GetValueUserName(ctx))
}
