package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/db"
	"BastionPay/bas-wallet-fission-share/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	Robber struct {
		Controllers
	}
)

//抢红包
func (this *Robber) Rob(ctx iris.Context) {
	//判断用户是否已抢，已抢则返回老数据
	//未抢，则抢
	param := new(api.RobberAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	tx := db.GDbMgr.Get().Begin()
	redUser, err := new(models.Robber).TxGetBy(tx, *param.RedId, *param.CountryCode, *param.Phone)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		tx.Rollback()
		return
	}
	if redUser != nil {
		this.ExceptionSeriveWithData(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_RED_ROBBER_EXISTS.Code(), redUser)
		//this.Response(ctx, redUser)
		tx.Rollback()
		return
	}

	red, err := new(models.Red).TxGetById(tx, *param.RedId)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
		tx.Rollback()
		return
	}

	if red == nil {
		ZapLog().Error("nofind err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_NOFOUND.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_NOFOUND.OriginDesc())
		tx.Rollback()
		return
	}

	if *red.RemainRob <= 0 {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO.OriginDesc())
		tx.Rollback()
		return
	}

	momey, ok := GenRedMoney(red)
	if !ok {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_RED_ROB_ZERO.OriginDesc())
		tx.Rollback()
		return
	}

	param.Coin = momey
	param2 := new(models.Robber).ParseAdd(param, red)
	_, err = param2.TxAdd(tx)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "has over")
		tx.Rollback()
		return
	}
	_, err = red.TxUpdate(tx)
	if err != nil {
		ZapLog().Error("Add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), "has over")
		tx.Rollback()
		return
	}
	tx.Commit()
	//事物结束
	this.Response(ctx, param2)

}

func (this *Robber) ListForBack(ctx iris.Context) {
	param := new(api.RobberList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Robber).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *Robber) ListForFront(ctx iris.Context) {
	param := new(api.RobberList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Robber).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}
