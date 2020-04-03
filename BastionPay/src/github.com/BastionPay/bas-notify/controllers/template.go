package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/models"
	"encoding/base64"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Template struct {
	Controllers
}

func (this *Template) Add(ctx iris.Context) {
	param := new(models.TemplateAdd)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	uflag, err := param.Unique()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	if !uflag {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "exist err")
		ZapLog().Error("exist err", zap.Error(err))
		return
	}
	if param.Content != nil {
		decodeBytes, err := base64.StdEncoding.DecodeString(*param.Content)
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "content Base64 decode err")
			ZapLog().Error("content Base64 decode err", zap.Error(err))
			return
		}
		content := string(decodeBytes)
		param.Content = &content
	}
	//if param.SmsPlatform == nil {
	//	oldGroup,err := new(models.TemplateGroup).GetAlive()
	//	if err != nil {
	//		ZapLog().Error( "database err", zap.Error(err))
	//	}
	//	if oldGroup != nil {
	//		param.SmsPlatform = oldGroup.SmsPlatform
	//	}
	//}
	if err := param.Add(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	this.Response(ctx, nil)
}

func (this *Template) Adds(ctx iris.Context) {
	params := make([]*models.TemplateAdd, 0)
	err := Tools.ShouldBindJSON(ctx, params)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	for i := 0; i < len(params); i++ {
		param := params[i]
		uflag, err := param.Unique()
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
		if !uflag {
			this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "exist err")
			ZapLog().Error("exist err", zap.Error(err))
			return
		}
		if param.Content != nil {
			decodeBytes, err := base64.StdEncoding.DecodeString(*param.Content)
			if err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "content Base64 decode err")
				ZapLog().Error("content Base64 decode err", zap.Error(err))
				return
			}
			content := string(decodeBytes)
			param.Content = &content
		}
		if err := param.Add(); err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
	}
	this.Response(ctx, nil)
}

func (this *Template) Update(ctx iris.Context) {
	param := new(models.Template)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.Content != nil {
		decodeBytes, err := base64.StdEncoding.DecodeString(*param.Content)
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "content Base64 decode err")
			ZapLog().Error("content Base64 decode err", zap.Error(err))
			return
		}
		content := string(decodeBytes)
		param.Content = &content
	}
	if err := param.Update(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	this.Response(ctx, nil)
}

func (this *Template) Updates(ctx iris.Context) {
	params := make([]*models.Template, 0)
	err := Tools.ShouldBindJSON(ctx, params)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	for i := 0; i < len(params); i++ {
		param := params[i]
		if param.Content != nil {
			decodeBytes, err := base64.StdEncoding.DecodeString(*param.Content)
			if err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "content Base64 decode err")
				ZapLog().Error("content Base64 decode err", zap.Error(err))
				return
			}
			content := string(decodeBytes)
			param.Content = &content
		}
		if err := param.Update(); err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
	}

	this.Response(ctx, nil)
}

func (this *Template) Saves(ctx iris.Context) {
	params := make([]*models.TemplateSave, 0)

	err := ctx.ReadJSON(&params)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	for i := 0; i < len(params); i++ {
		param := params[i]
		ok, err := govalidator.ValidateStruct(param)
		if !ok || err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
			ZapLog().Error("param err", zap.Error(err))
			return
		}
		if param.Content != nil {
			decodeBytes, err := base64.StdEncoding.DecodeString(*param.Content)
			if err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "content Base64 decode err")
				ZapLog().Error("content Base64 decode err", zap.Error(err))
				return
			}
			content := string(decodeBytes)
			param.Content = &content
		}
		if param.Id == nil {
			if param.GroupId == nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "groupid param err")
				ZapLog().Error("groupid param err")
				return
			}
			modParam := &models.TemplateAdd{
				Content:          param.Content,
				Name:             param.Name,
				Title:            param.Title,
				Type:             param.Type,
				Lang:             param.Lang,
				GroupId:          param.GroupId,
				YuntongxinTempId: param.YuntongxinTempId,
			}
			if err := modParam.Add(); err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
				ZapLog().Error("database err", zap.Error(err))
				return
			}
		} else {
			modParam := &models.Template{
				Id:      param.Id,
				Content: param.Content,
				Name:    param.Name,
				Title:   param.Title,
				Type:    param.Type,
				//Lang: param.Lang,
				//GroupId: param.GroupId,
				YuntongxinTempId: param.YuntongxinTempId,
			}
			if err := modParam.Update(); err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
				ZapLog().Error("database err", zap.Error(err))
				return
			}
		}

	}

	this.Response(ctx, nil)
}

//func (this *Template) List(ctx iris.Context) {
//
//}

func (this *Template) Gets(ctx iris.Context) {
	gidStr := ctx.URLParam("groupid")
	if len(gidStr) == 0 {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err")
		return
	}
	gidStr = strings.TrimSpace(gidStr)
	id, err := strconv.ParseUint(gidStr, 10, 32)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	arr, err := new(models.Template).GetsByGId(int(id))
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	for i := 0; i < len(arr); i++ {
		if arr[i].Content == nil {
			continue
		}
		encodeString := base64.StdEncoding.EncodeToString([]byte(*arr[i].Content))
		arr[i].Content = &encodeString
	}
	this.ResponseTemps(ctx, arr)
}
