package models

import (
	"BastionPay/pay-user-merchant-api/common"
	"BastionPay/pay-user-merchant-api/db"
)

type LogOperation struct {
	Id            *int       `json:"id,omitempty" gorm:"AUTO_INCREMENT:1;column:id;primary_key;not null"` //加上type:int(11)后AUTO_INCREMENT无效
	UserId        int       `json:"user_id,omitempty" gorm:"column:user_id;type:int(11)"`
	Operation     string     `json:"operation,omitempty" gorm:"column:operation;type:varchar(255)"`
	Ip            string     `json:"ip,omitempty" gorm:"column:ip;type:varchar(255)"`
	Country       string     `json:"country,omitempty" gorm:"column:country;type:varchar(255)"`
	City          string     `json:"city,omitempty" gorm:"column:city;type:varchar(255)"`
	CreatedAt     *int64 	 `json:"created_at,omitempty" gorm:"column:created_at;type:bigint(20)"`
	UpdatedAt     *int64 	 `json:"updated_at,omitempty" gorm:"column:updated_at;type:bigint(20)"`
	DeletedAt     *int64 	 `json:"deleted_at,omitempty" gorm:"column:deleted_at;type:bigint(20)"`
}

func (this *LogOperation) TableName() string {
	return "log_operation"
}

func (this *LogOperation) ParseAdd(userId int, ip, country, city, device string) *LogOperation {
	this.UserId = userId
	this.Ip = ip
	this.Country = country
	this.City = city
	return this
}


func (this *LogOperation) Add() ( error) {
	return db.GDbMgr.Get().Create(this).Error
}

func (this *LogOperation) List(userId uint, limit int, skip int) ([]*LogOperation, int, error) {
	var (
		data  []*LogOperation
		count int
		err   error
	)

	// 最近一周
	time := common.NowTimestamp() - 1000*60*60*24*180
	err = db.GDbMgr.Get().Model(&LogOperation{}).Where("user_id = ? AND created_at > ?", userId, time).
		Order("created_at desc").Limit(limit).Offset(skip).Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.GDbMgr.Get().Model(&LogOperation{}).Where("user_id = ? AND created_at > ?", userId, time).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}