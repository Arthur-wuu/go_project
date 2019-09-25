package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/models"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"time"
)

type (
	StaffMotivation struct {
		Controllers
	}

	ResponseAwardList struct {
		Date string         `json:"date,omitempty"`
		List *common.Result `json:"list,omitempty"`
	}
)

func (this *StaffMotivation) DayListForBack(ctx iris.Context) {
	param := new(api.StaffMotivationList)
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

	awd, err := new(models.StaffMotivation).GetDayAwardList(stime, param.Page, param.Size)

	if err != nil {
		ZapLog().Error("staff motivation award list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if awd == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	aList := awd.List.(*[]*models.StaffMotivationAwardResponse)

	for _, v := range *aList {
		v.Years = fmt.Sprintf("%.2f", float32(common.New().DayBeginTimestamp()-v.HiredAt)/(365*43200))
	}

	response := new(ResponseAwardList)
	response.List = awd
	response.Date = time.Unix(stime, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02")
	this.Response(ctx, response)
}

func (this *StaffMotivation) TotalListForBack(ctx iris.Context) {
	param := new(api.StaffMotivationList)
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

	awd, err := new(models.StaffMotivation).GetTotalAwardList(stime, param.Page, param.Size)

	if err != nil {
		ZapLog().Error("staff motivation award list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if awd == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	aList := awd.List.(*[]*models.StaffMotivationAwardResponse)

	for _, v := range *aList {
		v.Years = fmt.Sprintf("%.2f", float32(common.New().DayBeginTimestamp()-v.HiredAt)/(365*43200))
	}

	response := new(ResponseAwardList)
	response.List = awd
	response.Date = time.Unix(stime, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02")
	this.Response(ctx, response)
}
