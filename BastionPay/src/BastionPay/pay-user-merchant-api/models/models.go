package models

import (
	. "BastionPay/bas-base/log/zap"
	//"github.com/jinzhu/gorm"
	"BastionPay/pay-user-merchant-api/config"
)

func InitDbTable() {
	ZapLog().Info("start InitDbTable")
	if !config.GConfig.Db.Debug {
		ZapLog().Info("end InitDbTable")
		return
	}
	return
	//err := db.GDbMgr.Get().Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;").AutoMigrate(&LogLogin{},&LogOperation{}).Error
	//if err != nil {
	//	ZapLog().Error("AutoMigrate err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&Activity{}).AddIndex("idx_id_valid","id" ,"uuid", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	ZapLog().Info("End InitDbTable")
}