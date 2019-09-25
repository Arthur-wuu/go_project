package models

import (
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/jinzhu/gorm"
)

type LogModel struct {
	conn *gorm.DB
}

func NewLogModel(conn *gorm.DB) *LogModel {
	return &LogModel{conn: conn}
}

func (l *LogModel) CreateLoginLog(userId uint, ip string, country string, city string, device string) (uint, error) {
	data := LogLogin{
		UserId:  userId,
		Ip:      ip,
		Country: country,
		City:    city,
		Device:  device,
	}

	err := l.conn.Debug().Save(&data).Error
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

	err := l.conn.Debug().Save(&data).Error
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (l *LogModel) GetLoginLog(userId uint, limit int, skip int) ([]*LogLogin, int, error) {
	var (
		data  []*LogLogin
		count int
		err   error
	)

	// 最近一周
	time := common.NowTimestamp() - 1000*60*60*24*7
	err = l.conn.Debug().Model(&LogLogin{}).Where("user_id = ? AND created_at > ?", userId, time).
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
	time := common.NowTimestamp() - 1000*60*60*24*7
	err = l.conn.Debug().Model(&LogOperation{}).Where("user_id = ? AND created_at > ?", userId, time).
		Order("created_at desc").Count(&count).Limit(limit).Offset(skip).Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}
