package db

//import (
//	"BastionPay/bas-tools/sdk.aws.s3"
//	"BastionPay/bas-tools/sdk.aws.sts"
//)

//const(
//	TASKSTATUS_InQueue = 0
//	TASKSTATUS_Doing = 1
//	TASKSTATUS_Finish = 2
//	TASKSTATUS_EXPIRE = 3
//	TASKSTATUS_NOFIND = 4
//	TASKSTATUS_FAIL = 5
//)

//var (
//	GFileS3Sdk *sdk_aws_s3.S3Sdk
//	GStsSdk    *sdk_aws_sts.StsSdk
//)

//type TaskStatusInfo struct{
//	Order_id string        `json:"order_id,omitempty"`
//	Order_status int       `json:"order_status,omitempty"`
//	File_url     string    `json:"file_url,omitempty"`
//	Msg          string    `json:"msg,omitempty"`
//	File_name    string    `json:"file_name,omitempty"`
//}

//func NewTaskStatusInfo(orderId, file_url string, orderStatus int, msg string) *TaskStatusInfo {
//	return &TaskStatusInfo{
//		Order_id: orderId,
//		Order_status:  orderStatus,
//		File_url: file_url,
//		Msg: msg,
//	}
//}

type DbOptions struct {
	Host        string
	Port        string
	User        string
	Pass        string
	DbName      string
	MaxIdleConn int
	MaxOpenConn int
}

type RedisOptions struct {
	Network     string
	Host        string
	Port        string
	Password    string
	Database    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Prefix      string
}
