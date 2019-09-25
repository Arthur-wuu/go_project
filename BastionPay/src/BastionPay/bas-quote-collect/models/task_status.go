package models

import (
	"BastionPay/bas-quote-collect/db"
	"encoding/json"
	"fmt"
)

const (
	TASKSTATUS_Doing        = 0
	TASKSTATUS_Finish       = 1
	TASKSTATUS_FAIL         = 2
	TASKSTATUS_InQueue      = 3
	TASKSTATUS_EXPIRE       = 4
	TASKSTATUS_NOFIND       = 5
	EXPORT_Status_KeyPrefix = "dt"
)

type TaskStatusInfoForm struct {
	QueryId string `form:"query_id" valid:"required"`
	Lang    string `form:"language" valid:"-"`
}

type TaskStatusInfo struct {
	QueryId string `json:"query_id,omitempty"  valid:"required"`
	Status  int    `json:"status,omitempty"  valid:"optional"` // 0准备中  1成功  2失败
	//File_url     string    `json:"file_url,omitempty"`
	//Msg          string    `json:"msg,omitempty"`
	//File_name    string    `json:"file_name,omitempty"`
	Data      string `json:"data, omitempty" valid:"optional"`
	UserParam string `json:"user_param,omitempty" valid:"optional"`
}

func NewTaskStatusInfo(query_id, file_url string, queryStatus int, msg string) *TaskStatusInfo {
	return &TaskStatusInfo{
		QueryId: query_id,
		Status:  queryStatus,
	}
}

func (this *TaskStatusInfoForm) Get() (*TaskStatusInfo, error) {
	data, err := db.GRedis.Get(EXPORT_Status_KeyPrefix + this.QueryId)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &TaskStatusInfo{
			QueryId: this.QueryId,
			Status:  TASKSTATUS_NOFIND,
		}, nil
	}

	result := new(TaskStatusInfo)
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	//if result.Order_status != db.TASKSTATUS_Finish {
	//	return result, nil
	//}
	//this.genSignUrl(result.File_url, result.File_name)
	return result, err
}

func (this *TaskStatusInfo) Add() error {
	content, err := json.Marshal(this)
	fmt.Println("cont	", content, string(content))
	if err != nil {
		return err
	}
	if err = db.GRedis.Setex(EXPORT_Status_KeyPrefix+this.QueryId, string(content), 15*60); err != nil {
		return err
	}
	return nil
}

func (this *TaskStatusInfo) Update() error {
	content, err := json.Marshal(this)
	if err != nil {
		return err
	}
	if err = db.GRedis.Setex(EXPORT_Status_KeyPrefix+this.QueryId, string(content), 5*60); err != nil {
		return err
	}
	return nil
}
