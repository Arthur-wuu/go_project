package controllers

import (
	"BastionPay/bas-wallet-fission-share/models"
	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-fission-share/api"
	"BastionPay/bas-wallet-fission-share/db"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"time"
)

type (
	Red struct {
		Controllers
	}
)

func (this *Red) Add(ctx iris.Context) {
	//事物：判断活动存在，分配币数，创建红包
	param := new(api.RedAdd)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//一个事物
	tx := db.GDbMgr.Get().Begin()
	acty, err := new(models.Activity).TxGetById(tx, *param.ActivityId)
	if err != nil {
		ZapLog().Error("databse err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		tx.Rollback()
		return
	}
	if acty == nil {
		ZapLog().Error("nofind activity")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NOFOUND.Code(), "nofind activity")
		tx.Rollback()
		return
	}
	nowTime := time.Now().Unix()
	if *acty.OnAt > nowTime {
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NO_ON.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_NO_ON.OriginDesc())
		ZapLog().Error("acty notStarted err")
		tx.Rollback()
		return
		//未开始
	}
	if *acty.OffAt <= nowTime {
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_HAS_OFF.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_HAS_OFF.OriginDesc())
		//已经结束
		ZapLog().Error("acty have Over err")
		tx.Rollback()
		return
	}
	if *acty.RemainRed <= 0 {
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_ZERO.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_ZERO.OriginDesc())
		ZapLog().Error("no reds err")
		tx.Rollback()
		return
		//红包完
	}

	//判断用户有否重复创建红包，订单id是否存在
	param2 := new(models.Red).ParseAdd(param, acty)
	uFlag, err := param2.TxUnique(tx)
	if err != nil {
		ZapLog().Error("datebase err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		tx.Rollback()
		return
	}
	if !uFlag {
		ZapLog().Error("orderid exist err")
		this.ExceptionSerive(ctx, apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_EXISTS.Code(), apibackend.BASERR_ACTIVITY_FISSIONSHARE_ACTIVITY_RED_EXISTS.OriginDesc())
		tx.Rollback()
		return
	}
	_, err = param2.TxAdd(tx)
	if err != nil {
		ZapLog().Error("Add err"+fmt.Sprintf("%d", *param2.ActivityId), zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		tx.Rollback()
		return
	}

	_, err = acty.TxUpdate(tx)
	if err != nil {
		ZapLog().Error("ReduceRed err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		tx.Rollback()
		return
	}
	tx.Commit()
	//事物结束
	this.Response(ctx, param2)
}

func (this *Red) ListforBack(ctx iris.Context) {
	param := new(api.RedList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Red).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *Red) ListForFront(ctx iris.Context) {
	param := new(api.RedList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	data, err := new(models.Red).ParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("List err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}
