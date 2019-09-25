package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/models"
	"github.com/kataras/iris/context"
	"go.uber.org/zap"
)

type (
	Bulletin struct {
		Controllers
	}
)

func (b *Bulletin) Export(ctx context.Context) {
	taskExport := new(models.TaskExportInfo)

	err := Tools.ShouldBindJSON(ctx, taskExport)
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}
	if err := taskExport.PreProduce(); err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	status, err := taskExport.InQueue()
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		ZapLog().Error("InQueue err", zap.Error(err))
		return
	}

	b.Response(ctx, status)
	return
}

func (b *Bulletin) Cancel(ctx context.Context) {
	var cancelInfo models.TaskCancelInfo

	err := Tools.ShouldBindQuery(ctx, &cancelInfo)
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = cancelInfo.Cancel()
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		ZapLog().Error("Cancel err", zap.Error(err))
		return
	}
	b.Response(ctx, "")
}

func (b *Bulletin) GetStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfoForm

	err := Tools.ShouldBindQuery(ctx, &statusInfo)
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	result, err := statusInfo.Get()
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "get info fail")
		ZapLog().Error("getstatus err", zap.Error(err))
		return
	}

	b.Response(ctx, result)
}

func (b *Bulletin) AddStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfo

	err := Tools.ShouldBindQuery(ctx, &statusInfo)
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = statusInfo.Add()
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "add info fail")
		ZapLog().Error("AddStatus err", zap.Error(err))
		return
	}

	b.Response(ctx, "")
}

func (b *Bulletin) UpdateStatus(ctx context.Context) {
	var statusInfo models.TaskStatusInfo

	err := Tools.ShouldBindQuery(ctx, &statusInfo)
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), err.Error())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	err = statusInfo.Update()
	if err != nil {
		b.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "UpdateStatus info fail")
		ZapLog().Error("UpdateStatus err", zap.Error(err))
		return
	}

	b.Response(ctx, "")
}
