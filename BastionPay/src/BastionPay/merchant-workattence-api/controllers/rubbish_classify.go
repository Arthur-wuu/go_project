package controllers

import (
	"BastionPay/bas-api/apibackend"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/api"
	"BastionPay/merchant-workattence-api/common"
	"BastionPay/merchant-workattence-api/models"
	"BastionPay/merchant-workattence-api/modules"
	"encoding/json"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

type (
	RubbishClassify struct {
		Controllers
	}

	ResponseDepartmenInfo struct {
		Id           int                      `json:"id,omitempty"`
		Date         string                   `json:"date,omitempty"`
		TotalNumbers int                      `json:"total_numbers,omitempty"`
		Grade        []RubbishScore           `json:"grade,omitempty"`
		Department   []*models.DepartmentInfo `json:"department,omitempty"`
	}

	RubbishScore struct {
		Desc  string `json:"desc,omitempty"`
		Score int    `json:"score"`
	}

	ResponseSendDetail struct {
		Id           int            `json:"id,omitempty"`
		Date         string         `json:"date,omitempty"`
		TotalNumbers int            `json:"total_numbers,omitempty"`
		Grade        []RubbishScore `json:"grade,omitempty"`
		List         *common.Result `json:"list,omitempty"`
	}
)

func (this *RubbishClassify) SendAward(ctx iris.Context) {
	param := new(api.RubbishClassifyAwardSend)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)

	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	go new(modules.RClassify).SendAward(*param.Id)
	this.Response(ctx, map[string]string{"message": "sending rewards in the background"})

}

func (this *RubbishClassify) ListForBack(ctx iris.Context) {
	param := new(api.BkRubbishClassifyAwardList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.RubbishClassify).BkParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("rubbish classify award list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}
	if vip == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *RubbishClassify) DepartmentListForBack(ctx iris.Context) {
	acModel := new(models.AccountMap)
	depart, err := acModel.GetDepartmentInfo()

	if err != nil {
		ZapLog().Error("get department info err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if len(depart) == 0 {
		ZapLog().Error("department info len 0", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_DATA_NULL.Code(), apibackend.BASERR_OBJECT_DATA_NULL.Desc())
		return
	}

	tCount, err := acModel.GetValidEmployerCount()

	if err != nil {
		ZapLog().Error("get employer count err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	departResponse := new(ResponseDepartmenInfo)
	departResponse.Department = depart
	departResponse.TotalNumbers = tCount
	departResponse.Date = common.New().NowDataString()
	departResponse.Grade = []RubbishScore{{Desc: "差", Score: 0}, {Desc: "良", Score: 1}, {Desc: "优", Score: 2}}
	this.Response(ctx, departResponse)
}

func (this *RubbishClassify) AddForBack(ctx iris.Context) {
	param := new(api.ResponseDepartmenInfo)
	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)

	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//数据库添加
	data, err := new(models.RubbishClassifyRecord).Add(param)

	if err != nil {
		ZapLog().Error("rubbish classify add err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *RubbishClassify) RecordListForBack(ctx iris.Context) {
	param := new(api.BkRubbishClassifyList)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	vip, err := new(models.RubbishClassifyRecord).BkParseList(param).List(param.Page, param.Size)
	if err != nil {
		ZapLog().Error("rubbish classify record list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if vip == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	this.Response(ctx, vip)
}

func (this *RubbishClassify) GetRecordForBackEdit(ctx iris.Context) {
	param := new(api.BkEditRubbishClassify)

	err := Tools.ShouldBindJSON(ctx, param)
	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	rcd, err := new(models.RubbishClassifyRecord).GetById(*param.Id)

	if err != nil {
		ZapLog().Error("get classify record err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if rcd == nil {
		ZapLog().Error("classify record not found", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	var depart []*models.DepartmentInfo
	err = json.Unmarshal([]byte(*rcd.DepartmentInfo), &depart)

	if err != nil {
		ZapLog().Error("department info unmarshal err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_INCORRECT_FORMAT.Code(), apibackend.BASERR_INCORRECT_FORMAT.Desc())
		return
	}

	departResponse := new(ResponseDepartmenInfo)
	departResponse.Id = *rcd.Id
	departResponse.Department = depart
	departResponse.TotalNumbers = *rcd.TotalNumbers
	departResponse.Date = *rcd.ScoreDate
	departResponse.Grade = []RubbishScore{{Desc: "差", Score: 0}, {Desc: "良", Score: 1}, {Desc: "优", Score: 2}}
	this.Response(ctx, departResponse)
}

func (this *RubbishClassify) RecordUpdateForBack(ctx iris.Context) {
	param := new(api.ResponseDepartmenInfo)
	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)

	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	//数据库添加
	data, err := new(models.RubbishClassifyRecord).Update(param)

	if err != nil {
		ZapLog().Error("rubbish classify update err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), err.Error())
		return
	}

	this.Response(ctx, data)
}

func (this *RubbishClassify) SendDetailForBack(ctx iris.Context) {
	param := new(api.RubbishClassifyAwardSendList)

	//参数检测
	err := Tools.ShouldBindJSON(ctx, param)

	if err != nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc())
		ZapLog().Error("param err", zap.Error(err))
		return
	}

	rcd, err := new(models.RubbishClassifyRecord).GetById(param.Id)

	if err != nil {
		ZapLog().Error("get rubbish classify record  err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	if rcd == nil {
		this.ExceptionSerive(ctx, apibackend.BASERR_OBJECT_NOT_FOUND.Code(), apibackend.BASERR_OBJECT_NOT_FOUND.Desc())
		return
	}

	awdList, err := new(models.RubbishClassify).AwardListByRecordId(param.Id, param.Page, param.Size)

	if err != nil {
		ZapLog().Error("get rubbish classify award list err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_DATABASE_ERROR.Code(), apibackend.BASERR_DATABASE_ERROR.Desc())
		return
	}

	response := new(ResponseSendDetail)
	response.Id = *rcd.Id
	response.Date = *rcd.ScoreDate
	response.TotalNumbers = *rcd.TotalNumbers
	response.Grade = []RubbishScore{{Desc: "差", Score: 0}, {Desc: "良", Score: 1}, {Desc: "优", Score: 2}}
	response.List = awdList
	this.Response(ctx, response)
}
