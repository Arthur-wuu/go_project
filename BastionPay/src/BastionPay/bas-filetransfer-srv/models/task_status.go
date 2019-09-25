package models

import (
	"BastionPay/bas-filetransfer-srv/db"
	"encoding/json"
)

const (
	TASKSTATUS_Doing   = 0
	TASKSTATUS_Finish  = 1
	TASKSTATUS_FAIL    = 2
	TASKSTATUS_InQueue = 3
	TASKSTATUS_EXPIRE  = 4
	TASKSTATUS_NOFIND  = 5
)

type TaskStatusInfoForm struct {
	Order_id string `form:"order_id" valid:"required"`
}

type TaskStatusInfo struct {
	Order_id     string `json:"order_id,omitempty"  valid:"required"`
	Order_status int    `json:"order_status,omitempty"  valid:"required"`
	File_url     string `json:"file_url,omitempty"`
	Msg          string `json:"msg,omitempty"`
	File_name    string `json:"file_name,omitempty"`
	Data         string `json:"data, omitempty"`
	UserParam    string `json:"user_param,omitempty"`
}

func NewTaskStatusInfo(orderId, file_url string, orderStatus int, msg string) *TaskStatusInfo {
	return &TaskStatusInfo{
		Order_id:     orderId,
		Order_status: orderStatus,
		File_url:     file_url,
		Msg:          msg,
	}
}

func (this *TaskStatusInfoForm) Get() (*TaskStatusInfo, error) {
	data, err := db.GRedis.Get(EXPORT_Status_KeyPrefix + this.Order_id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &TaskStatusInfo{
			Order_id:     this.Order_id,
			Order_status: TASKSTATUS_NOFIND,
			Msg:          "nofind",
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
	if err != nil {
		return err
	}
	if err = db.GRedis.Setex(EXPORT_Status_KeyPrefix+this.Order_id, string(content), 15*60); err != nil {
		return err
	}
	return nil
}

func (this *TaskStatusInfo) Update() error {
	content, err := json.Marshal(this)
	if err != nil {
		return err
	}
	if err = db.GRedis.Setex(EXPORT_Status_KeyPrefix+this.Order_id, string(content), 5*60); err != nil {
		return err
	}
	return nil
}

//func (this *TaskStatusInfo) genSignUrl(url ,filename string) ( error) {
//	stat := sdk_aws_sts.Statement{
//		Action: "s3:GetObject",
//		Effect: "Allow",
//		Resource: "arn:aws:s3:::"+config.GConfig.Aws.FileBucket +"/"+filename,
//	}
//	cred, err := db.GStsSdk.GetFederationToken(5*60, "", []sdk_aws_sts.Statement{stat})
//	if err != nil {
//		return err
//	}
//	newCred := credentials.NewStaticCredentials(*cred.AccessKeyId, *cred.SecretAccessKey, *cred.SessionToken)
//
//	httpReq,err := http.NewRequest("GET" , url, nil)
//	if err != nil {
//		return err
//	}
//	sign := v4.NewSigner(newCred)
//	hdd,err := sign.Sign(httpReq, nil, "s3", config.GConfig.Aws.FileRegion, time.Now())
//	if err != nil {
//		return err
//	}
//	fmt.Println(hdd)
//	return nil
//}
