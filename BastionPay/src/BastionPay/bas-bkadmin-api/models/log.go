package models

import (
	l4g "github.com/alecthomas/log4go"
	"time"
)

// 登录日志
type LogLogin struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"default:null"`
	UserId    uint
	Ip        string `gorm:"type:varchar(255)"`
	Country   string `gorm:"type:varchar(255)"`
	City      string `gorm:"type:varchar(255)"`
	Device    string `gorm:"type:varchar(10)"`
}

// 操作日志
type LogOperation struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"default:null"`
	UserId    uint
	Operation string `gorm:"type:varchar(255)"`
	Ip        string `gorm:"type:varchar(255)"`
	Country   string `gorm:"type:varchar(255)"`
	City      string `gorm:"type:varchar(255)"`
}

type LogModel struct {
}

func NewLogModel() *LogModel {
	return &LogModel{}
}

func (l *LogModel) CreateLoginLog(userId uint, ip string, country string, city string, device string) (uint, error) {

	data := LogLogin{
		UserId:  userId,
		Ip:      ip,
		Country: country,
		City:    city,
		Device:  device,
	}

	err := DB.Save(&data).Error
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (l *LogModel) CreateOperationLog(userId uint, operation string, ip string, country string, city string) (uint, error) {
	data := LogOperation{
		UserId:    userId,
		Operation: operation,
		Ip:        ip,
		Country:   country,
		City:      city,
	}

	err := DB.Save(&data).Error
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

//add 删除bkadmin 中14天之前的数据
func (l *LogModel) DeleteLoginLog(remainDays int64) {
	//14天
	if remainDays < 1 {
		remainDays = 1
	}
	timestamp := NowTimestamp() - 1000*60*60*24*remainDays
	err = DB.Delete(&LogLogin{}, " created_at < ?", timestamp).Error

	if err != nil {
		l4g.Error("delete loginlog err[%s]", err.Error())
	}

	//time.Sleep(time.Second * time.Duration(60*60*24))
}

func (l *LogModel) DeleteOperatLog(remainDays int64) {
	if remainDays < 1 {
		remainDays = 1
	}
	timestamp := NowTimestamp() - 1000*60*60*24*remainDays
	err = DB.Delete(LogOperation{}, " created_at < ?", timestamp).Error
	if err != nil {
		l4g.Error("delete operatelog err[%s]", err.Error())
	}
	//time.Sleep(time.Second * time.Duration(60*60*24))
}

func (l *LogModel) GetLoginLog(userId uint, limit int, skip int) ([]*LogLogin, int, error) {
	var (
		data  []*LogLogin
		count int
		err   error
	)

	// 最近一周
	time := NowTimestamp() - 1000*60*60*24*7
	err = DB.Model(&LogLogin{}).Where("user_id = ? AND created_at > ?", userId, time).
		Order("created_at desc").Count(&count).Limit(limit).Offset(skip).Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (l *LogModel) GetOperationLog(userId uint, limit int, skip int) ([]*LogOperation, int, error) {
	var (
		data  []*LogOperation
		count int
		err   error
	)

	// 最近一周
	time := NowTimestamp() - 1000*60*60*24*7
	err = DB.Model(&LogOperation{}).Where("user_id = ? AND created_at > ?", userId, time).
		Order("created_at desc").Count(&count).Limit(limit).Offset(skip).Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (l *LogModel) timeStamp() (s int64) {
	s = time.Now().Unix()
	return s
}
