package models

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
	"BastionPay/bas-notify/db"
	"BastionPay/bas-notify/models/table"
	"go.uber.org/zap"
)

const (
	Notify_Level_Normal = 0
	Notify_Level_Low    = 0
	Notify_Level_High   = 1
)

const (
	Notify_Type_Sms  = 0
	Notify_Type_Mail = 1
	Notify_Type_Ding = 2
)

const (
	Notify_AliveMode_Dead = 0
	Notify_AliveMode_Live = 1
)

const (
	SMSPlatform_AWS        = 0
	SMSPlatform_CHUANGLAN  = 1
	SMSPlatform_TWL        = 2
	SMSPlatform_Nexmo      = 3
	SMSPlatform_YunTongXun = 4
)

func InitDbTable() {
	log.ZapLog().Info("start InitDbTable")
	if true || !config.GConfig.Db.Debug {
		log.ZapLog().Info("end InitDbTable")
		return
	}

	err := db.GDbMgr.Get().Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;").AutoMigrate(&table.Template{}, &table.TemplateGroup{}, &table.History{}).Error
	if err != nil {
		log.ZapLog().Error("AutoMigrate err", zap.Error(err))
	}
}
