package controllers

import (
	"github.com/kataras/iris"
	"BastionPay/bas-notify/models"
	"go.uber.org/zap"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
)

type History struct{
	Controllers
}

func (this *History) List(ctx iris.Context) {
	param := new(models.TemplateHistoryList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error( "param err", zap.Error(err))
		return
	}
	if param.Total_lines <= 0 {
		param.Total_lines, err = new(models.TemplateHistory).Count(param.GroupId)
		if err != nil {
			this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
			ZapLog().Error( "database err", zap.Error(err))
			return
		}
	}
	hisArr,err := param.List()
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "database err")
		ZapLog().Error( "database err", zap.Error(err))
		return
	}
	list := &TemplateHistoryList{
		Total_lines: param.Total_lines,
		Max_disp_lines: param.Max_disp_lines,
		Page_index: param.Page_index,
		TemplateHistorys: hisArr,
	}
	this.ResponseHistoryList(ctx, list)
}