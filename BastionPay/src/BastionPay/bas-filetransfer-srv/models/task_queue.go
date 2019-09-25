package models

import (
	"BastionPay/bas-filetransfer-srv/common"
	"BastionPay/bas-filetransfer-srv/db"

	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/config"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)

const (
	EXPORT_List_Key         = "filetransfer_export_list"
	EXPORT_Status_KeyPrefix = "filetransfer_status_"
	EXPORT_File_Exd         = ".zip"
)

type TaskExportInfo struct {
	Sql          []byte `json:"sql" valid:"required"` //这里要用[]byte, json自动base64
	Dbname       string `json:"dbname" valid:"required"`
	PageSize     int    `json:"pagesize"`
	WaitExpireAt int64  `json:"wait_expire_at"` //排队到期时间
	MaxOpTime    int64  `json:"max_op_time"`    //操作时间
	PageHalt     int    `json:"pagehalt"`
	Max_lines    uint64 `json:"max_lines"`
	OrderId      string `json:"order_id"`
	FileName     string `json:"file_name"`
	FileFormat   string `json:"file_format"`
	Caller       string `json:"caller"`
}

func (this *TaskExportInfo) PreProduce() error {
	sqlUpper := strings.ToUpper(string(this.Sql))
	if ok := strings.Contains(sqlUpper, "LIMIT"); ok {
		return errors.New("sql have LIMIT")
	}
	if ok := strings.Contains(sqlUpper, "OFFSET"); ok {
		return errors.New("sql have OFFSET")
	}
	taskConf := config.GConfig.Task
	if this.PageSize <= 2 || this.PageSize >= taskConf.MaxPage {
		this.PageSize = taskConf.MaxPage
	}
	expire := this.WaitExpireAt
	if this.WaitExpireAt > 100000 {
		expire = this.WaitExpireAt - time.Now().Unix()
	}
	if expire <= 3 || expire >= taskConf.MaxWaitTime {
		expire = taskConf.MaxWaitTime
	}
	this.WaitExpireAt = expire + time.Now().Unix()
	if this.MaxOpTime <= 3 || this.MaxOpTime >= taskConf.MaxOpTime {
		this.MaxOpTime = taskConf.MaxOpTime
	}
	if this.PageHalt <= 0 {
		this.PageHalt = 0
	}
	if this.Max_lines >= taskConf.MaxRecords {
		this.Max_lines = taskConf.MaxRecords
	}
	if len(this.FileFormat) == 0 {
		this.FileFormat = "xlsx"
	}
	this.FileFormat = strings.ToLower(this.FileFormat)
	if !(this.FileFormat == "xlsx" || this.FileFormat == "csv") {
		return errors.New("file_format not support")
	}
	return nil
}

func (this *TaskExportInfo) InQueue() (*TaskStatusInfo, error) {
	_, ok := db.GDbMgrs[this.Dbname]
	if !ok {
		return nil, errors.New("db not in config")
	}

	//生成文件名
	this.FileName = this.genFileName(this.Caller, string(this.Sql))
	ZapLog().Info("wait In Queue", zap.String("filename", this.FileName), zap.String("sql", string(this.Sql)))

	//判断文件是否存在于s3
	if config.GConfig.Task.FileUseExist {
		fileStatus, err := db.GFileS3Sdk.Status(config.GConfig.Aws.FileBucket, this.FileName+EXPORT_File_Exd)
		if err != nil {
			ZapLog().Error("s3 status err", zap.Error(err))
			return nil, err
		}
		if fileStatus.Exist() && !fileStatus.Timeout() {
			newAddr, err := this.genSignUrl(fileStatus.Addr(), this.FileName+EXPORT_File_Exd)
			if err != nil {
				ZapLog().Error("genSignUrl err", zap.Error(err), zap.String("orderid", this.OrderId))
				return nil, err
			}
			return &TaskStatusInfo{
				Order_id:     "275XGS34OPWMZP567890XAS3",
				Order_status: TASKSTATUS_Finish,
				File_url:     newAddr,
			}, nil
		}
	}

	//队列是否过长
	overFlag, err := this.overTaskList()
	if overFlag {
		ZapLog().Warn("Queue is too busy")
		return nil, errors.New("too long queue")
	} else if err != nil {
		ZapLog().Warn("redis llen Queue err", zap.Error(err))
		return nil, err
	}

	//入队
	taskStatusInfo := NewTaskStatusInfo(common.New().GenerateUuid(), "", TASKSTATUS_InQueue, "")
	this.OrderId = taskStatusInfo.Order_id

	if err = this.AddTaskStatus(taskStatusInfo, this.WaitExpireAt+this.MaxOpTime); err != nil {
		ZapLog().Error("AddTaskStatus err", zap.Error(err))
		return nil, err
	}

	ZapLog().Info("push Task in queue ok", zap.String("filename", this.FileName), zap.String("orderid", this.OrderId))
	if err = this.pushTask(); err != nil {
		ZapLog().Error("pushTask err", zap.Error(err))
		return nil, err
	}
	return taskStatusInfo, nil
}

func (this *TaskExportInfo) AddTaskStatus(info *TaskStatusInfo, expireat int64) error {
	content, err := json.Marshal(info)
	if err != nil {
		return err
	}
	expire := time.Now().Unix() - expireat
	if expire <= 5 {
		expire = 3 * 60
		//return errors.New("expire timeout")
	}
	expire += config.GConfig.Task.StatusKeepTime
	_, err = db.GRedis.Do("setex", EXPORT_Status_KeyPrefix+this.OrderId, expire, string(content))
	if err != nil {
		return err
	}
	return err
}

func (this *TaskExportInfo) updateTaskStatus(file_url string, status int, msg string) error {
	info := NewTaskStatusInfo(this.OrderId, file_url, status, msg)
	content, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = db.GRedis.Do("set", EXPORT_Status_KeyPrefix+this.OrderId, string(content))
	return err
}

func (this *TaskExportInfo) existTaskStatus(orderid string) error {
	return nil
}

func (this *TaskExportInfo) pushTask() error {
	taskInfoContent, err := json.Marshal(this)
	if err != nil {
		return err
	}
	_, err = db.GRedis.Do("lpush", EXPORT_List_Key, taskInfoContent)
	if err != nil {
		return err
	}
	return err
}

func (this *TaskExportInfo) overTaskList() (bool, error) {
	queueLen, err := db.GRedis.Llen(EXPORT_List_Key)
	if err != nil {
		return false, err
	}
	if queueLen > config.GConfig.Task.MaxWaitLen {
		return true, nil
	}
	return false, nil
}

func (this *TaskExportInfo) genFileName(str ...string) string {
	bigStr := "" //排序的目的是保证sql中字段顺序 颠倒，导致生成的文件名称不一样
	for i := 0; i < len(str); i++ {
		bigArr := make([]string, 0)
		arr := strings.Split(str[i], " ")
		for j := 0; j < len(arr); j++ {
			subarr := strings.Split(arr[i], ",")
			bigArr = append(bigArr, subarr...)
		}
		sort.Strings(bigArr)
		bigStr += strings.Join(bigArr, " ") + " "
	}
	return fmt.Sprintf("%X", md5.Sum([]byte(bigStr)))
}
