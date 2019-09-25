package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/models"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
)

func NewFileTransferCtl() *FileTransferCtl {
	return &FileTransferCtl{}
}

type FileTransferCtl struct {
	Controllers
}

func (this *FileTransferCtl) AddStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfo

	err := Tools.ShouldBindJSON(ctx, &statusInfo)
	if err != nil {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = statusInfo.Add()
	if err != nil {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "add info fail")
		ZapLog().Error("AddStatus err", zap.Error(err))
		return
	}

	this.Controllers.Response(ctx, nil)
}

func (this *FileTransferCtl) UpdateStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfo

	err := Tools.ShouldBindJSON(ctx, &statusInfo)
	if err != nil {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = statusInfo.Update()
	if err != nil {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "UpdateStatus info fail")
		ZapLog().Error("UpdateStatus err", zap.Error(err))
		return
	}

	this.Controllers.Response(ctx, nil)
}

func (this *FileTransferCtl) GetStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfoForm

	statusInfo.QueryId = ctx.URLParam("query_id")
	ok, err := govalidator.ValidateStruct(&statusInfo)
	if err != nil || !ok {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Any("error", err), zap.Bool("ok", ok))
		return
	}

	result, err := statusInfo.Get()
	if err != nil {
		this.Controllers.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "get info fail")
		ZapLog().Error("getstatus err", zap.Error(err))
		return
	}

	this.Controllers.Response(ctx, result)
}
