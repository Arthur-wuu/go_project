package models

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/db"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

func InitDbTable() {
	log.ZapLog().Info("start InitDbTable")

	if !config.GConfig.Db.Debug {
		log.ZapLog().Info("end InitDbTable")
		return
	}

	err := db.GDbMgr.Get().Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;").AutoMigrate(&AccountMap{}, &AwardRecord{}, &CheckinRecord{}, &StaffMotivation{}, &RubbishClassify{}, &RubbishClassifyRecord{}).Error

	if err != nil {
		log.ZapLog().Error("AutoMigrate err", zap.Error(err))
	}

	log.ZapLog().Info("end InitDbTable")
}
