package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/db"
	"BastionPay/bas-notify/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
)

type TemplateGroup struct {
	Controllers
}

func (this *TemplateGroup) List(ctx iris.Context) {
	param := new(models.TemplateGroupList)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.Total_lines <= 0 {
		param.Total_lines, err = param.LikeCount()
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
	}
	groupArr, err := param.List()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	for i := 0; i < len(groupArr); i++ {
		groupArr[i].Langs, err = new(models.Template).GetLangByGId(*groupArr[i].Id)
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
	}
	list := &TemplateGroupList{
		Total_lines:    param.Total_lines,
		Max_disp_lines: param.Max_disp_lines,
		Page_index:     param.Page_index,
		TemplateGroups: groupArr,
	}
	this.ResponseGroupList(ctx, list)
}

func (this *TemplateGroup) Add(ctx iris.Context) {
	param := new(models.TemplateGroupAdd)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.Name != nil {
		*param.Name = strings.Replace(*param.Name, " ", "", -1)
	}
	if param.Alive == nil {
		param.Alive = new(int)
		*param.Alive = 1
	}
	flag, err := param.Unique()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	if !flag {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_EXISTS.Code(), "exist err")
		ZapLog().Error("exist err", zap.Error(err))
		return
	}
	if param.SmsPlatform == nil {
		oldGroup, err := new(models.TemplateGroup).GetAlive()
		if err != nil {
			ZapLog().Error("database err", zap.Error(err))
		}
		if oldGroup != nil {
			param.SmsPlatform = oldGroup.SmsPlatform
		}
	}

	if _, err := param.Add(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	this.Response(ctx, nil)
}

func (this *TemplateGroup) Update(ctx iris.Context) {
	param := new(models.TemplateGroup)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if err := param.Update(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	this.Response(ctx, nil)
}

func (this *TemplateGroup) Copys(ctx iris.Context) {
	param := new(models.TemplateGroupCopy)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	srcGroupTemp, err := new(models.TemplateGroup).GetByid(*param.Id)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}
	if srcGroupTemp == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), "nofind err")
		ZapLog().Error("nofind err")
		return
	}

	subTemps, err := new(models.Template).GetsByGId(*param.Id)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		return
	}

	for i := 0; i < len(param.SubName); i++ {
		groupTemp, err := new(models.TemplateGroup).GetByNameAndType(*srcGroupTemp.Name, *srcGroupTemp.Type, param.SubName[i])
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error("database err", zap.Error(err))
			return
		}
		if groupTemp != nil {
			modelParam := &models.TemplateGroup{
				Id: groupTemp.Id,
				//Name: groupTemp.Name,
				//SubName: param.SubName[i],
				Detail: srcGroupTemp.Detail,
				Alive:  srcGroupTemp.Alive,
				//Type: srcGroupTemp.Type,
				Author:      srcGroupTemp.Author,
				Editor:      srcGroupTemp.Editor,
				SmsPlatform: srcGroupTemp.SmsPlatform,
			}
			if err = modelParam.Update(); err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
				ZapLog().Error("database err", zap.Error(err))
				return
			}
		} else {
			modelParam := &models.TemplateGroupAdd{
				Name:        srcGroupTemp.Name,
				SubName:     param.SubName[i],
				Detail:      srcGroupTemp.Detail,
				Alive:       srcGroupTemp.Alive,
				Type:        srcGroupTemp.Type,
				Author:      srcGroupTemp.Author,
				Editor:      srcGroupTemp.Editor,
				SmsPlatform: srcGroupTemp.SmsPlatform,
			}
			if groupTemp, err = modelParam.Add(); err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
				ZapLog().Error("database err", zap.Error(err))
				return
			}
		}

		for j := 0; j < len(subTemps); j++ {
			temp, err := new(models.Template).GetByGIdAndLang(*groupTemp.Id, *subTemps[j].Lang)
			if err != nil {
				this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
				ZapLog().Error("database err", zap.Error(err))
				return
			}
			if temp != nil {
				modelParam := models.Template{
					Id:      temp.Id,
					Title:   subTemps[j].Title,
					Content: subTemps[j].Content,
					Type:    subTemps[j].Type,
				}
				if err = modelParam.Update(); err != nil {
					this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
					ZapLog().Error("database err", zap.Error(err))
					return
				}
			} else {
				modelParam := &models.TemplateAdd{
					Title:   subTemps[j].Title,
					Lang:    subTemps[j].Lang,
					Content: subTemps[j].Content,
					Type:    subTemps[j].Type,
					GroupId: groupTemp.Id,
				}
				if err = modelParam.Add(); err != nil {
					this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
					ZapLog().Error("database err", zap.Error(err))
					return
				}
			}
		}
	}

	this.Response(ctx, nil)
}

//func (this *TemplateGroup) Gets(ctx iris.Context) {
//	param := new(models.TemplateGroup)
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//}

func (this *TemplateGroup) Alive(ctx iris.Context) {
	param := new(models.TemplateGroup)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.Alive == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	//tx := db.GDbMgr.Get().Begin()
	if err := param.Update(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		//tx.Rollback()
		return
	}

	//if err := new(models.Template).TxAliveByGid(tx, *param.Id, *param.Alive); err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
	//	ZapLog().Error( "database err", zap.Error(err))
	//	tx.Rollback()
	//	return
	//}
	//tx.Commit()
	this.Response(ctx, nil)
}

func (this *TemplateGroup) Del(ctx iris.Context) {
	param := new(models.TemplateGroup)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	tx := db.GDbMgr.Get().Begin()
	if err := param.TxDel(tx); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		tx.Rollback()
		return
	}

	if err := new(models.Template).TxDelByGid(tx, *param.Id); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		tx.Rollback()
		return
	}
	tx.Commit()
	this.Response(ctx, nil)
}

func (this *TemplateGroup) SetRecipient(ctx iris.Context) {
	param := new(models.TemplateGroup)
	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.DefaultRecipient == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	*param.DefaultRecipient = strings.Replace(*param.DefaultRecipient, " ", "", -1)

	//tx := db.GDbMgr.Get().Begin()
	if err := param.Update(); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		//tx.Rollback()
		return
	}

	//tempParam := &models.Template{}
	//if err := tempParam.TxSetDefaultRecipient(tx, *param.Id, *param.DefaultRecipient); err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
	//	ZapLog().Error( "database err", zap.Error(err))
	//	tx.Rollback()
	//	return
	//}
	//tx.Commit()
	this.Response(ctx, nil)
}

func (this *TemplateGroup) SetSmsPlatom(ctx iris.Context) {
	param := new(models.TemplateGroup)
	err := ctx.ReadJSON(param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if param.SmsPlatform == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	//tx := db.GDbMgr.Get().Begin()
	//if err := param.Update(); err != nil {
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
	//	ZapLog().Error( "database err", zap.Error(err))
	//	tx.Rollback()
	//	return
	//}
	//
	//tempParam := &models.Template{}
	if err := param.SetAllSmsPlatform(*param.SmsPlatform); err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error("database err", zap.Error(err))
		//tx.Rollback()
		return
	}
	//tx.Commit()
	this.Response(ctx, nil)
}
