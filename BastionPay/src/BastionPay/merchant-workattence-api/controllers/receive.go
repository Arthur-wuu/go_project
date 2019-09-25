package controllers

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/models"
	"encoding/json"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"io/ioutil"
)

type (
	Receive struct {
		Controllers
	}
)

func (this *Receive) CheckinRecord(ctx iris.Context) {
	bodybuf, err := ioutil.ReadAll(ctx.Request().Body)

	if err != nil {
		ZapLog().Error("request buf read err", zap.Error(err))
		return
	}

	param := new(api.PushCheckin)
	paramSlice := make([]api.PushCheckin, 0)
	err = json.Unmarshal(bodybuf, param)

	if err != nil {
		err = json.Unmarshal(bodybuf, &paramSlice)

		if err != nil {
			ZapLog().Error("request buf array unmarshal err", zap.Error(err))
			return
		}

		returnData := make([]string, len(paramSlice))

		for k, v := range paramSlice {
			if *v.Id == "5" {
				rcd, err := new(models.CheckinRecord).Add(&v)

				if err != nil {
					ZapLog().Error("checkin record add err", zap.Error(err))
					return
				}

				go this.SendAward(rcd, common.New().DayBeginTimestamp())
			}

			returnData[k] = *v.Id
		}

		this.ResponsePush(ctx, returnData)
	} else {
		var returnData = [1]string{*param.Id}

		if *param.Id != "5" {
			this.ResponsePush(ctx, returnData)
		} else {
			rcd, err := new(models.CheckinRecord).Add(param)

			if err != nil {
				ZapLog().Error("checkin record add err", zap.Error(err))
				return
			}

			go this.SendAward(rcd, common.New().DayBeginTimestamp())
			this.ResponsePush(ctx, returnData)
		}
	}
}

func (this *Receive) SendAward(rcd *models.CheckinRecord, bTime int64) {
	arm := new(models.AwardRecord)
	//check today award had sended
	num, err := arm.SendCheckDay(*rcd.StaffId, bTime, *rcd.CheckinAt)

	if err != nil || num > 0 {
		return
	}

	//get accout info
	act, err := new(models.AccountMap).GetByStaffId(*rcd.StaffId)

	if err != nil {
		return
	}

	//add send record and send award
	coin := config.GConfig.Award.Checkin.Coin
	symbol := config.GConfig.Award.Checkin.Symbol
	awd, err := arm.AddAuto(*rcd.Id, *act, *rcd.StaffId, coin, symbol, 1)

	if err != nil {
		return
	}

	models.SChan <- models.SendChan{
		Coin:       coin,
		Symbol:     symbol,
		MerchantId: config.GConfig.Award.MerchantId,
		Times:      config.GConfig.Award.SendTimes,
		AwardId:    *awd.Id,
		AccountId:  *awd.AccId,
	}
}
