package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/models"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"time"
)

type (
	WorkAttendance struct {
		Controllers
	}

	ResponseOvertimeAwardList struct {
		Date string         `json:"date,omitempty"`
		List *common.Result `json:"list,omitempty"`
	}
)

func (this *WorkAttendance) OvertimeAwardListForBack(ctx iris.Context) {
	param := new(api.OvertimeAwardList)
	var stime int64

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	if param.Datetime == nil {
		stime = common.New().DayBeginTimestamp()
	} else {
		stime = common.New().DateToDayBeginTimestamp(*param.Datetime)
	}

	awd, err := new(models.AwardRecord).GetOvertimeList(stime, param.Page, param.Size)

	if err != nil {
		ZapLog().Error("overtime award list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if awd == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	aList := awd.List.(*[]*models.ResponseOvertimeList)
	ckModel := new(models.CheckinRecord)

	for _, v := range *aList {
		lCheckR, err := ckModel.GetRecordById(v.CheckinId)

		if err != nil || lCheckR == nil {
			v.Lasttime = ""
		} else {
			v.Lasttime = time.Unix(*lCheckR.CheckinAt, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02 15:04:05")
		}

		fCheckR, err := ckModel.GetEarliestCheckinRecord(v.UserId)

		if err != nil || fCheckR == nil {
			v.Firsttime = ""
		} else {
			v.Firsttime = time.Unix(*fCheckR.CheckinAt, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02 15:04:05")
		}

		if v.Firsttime != "" && v.Lasttime != "" {
			v.Duration = *lCheckR.CheckinAt - *fCheckR.CheckinAt - 37800
		} else {
			v.Duration = 0
		}
	}

	response := new(ResponseOvertimeAwardList)
	response.List = awd
	response.Date = time.Unix(stime, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02")
	this.Response(ctx, response)
}
