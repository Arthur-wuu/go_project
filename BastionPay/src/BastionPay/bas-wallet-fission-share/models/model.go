package models

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/bas-wallet-fission-share/config"
	"BastionPay/bas-wallet-fission-share/db"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

func InitDbTable() {
	log.ZapLog().Info("start InitDbTable")
	if !config.GConfig.Db.Debug {
		log.ZapLog().Info("end InitDbTable")
		return
	}
	err := db.GDbMgr.Get().Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1;").AutoMigrate(&Activity{}, &Red{}, &Robber{}, &Slogan{}).Error
	if err != nil {
		log.ZapLog().Error("AutoMigrate err", zap.Error(err))
	}
	err = db.GDbMgr.Get().Model(&Activity{}).AddIndex("idx_id_valid", "id", "valid").Error
	if err != nil {
		log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	}
	err = db.GDbMgr.Get().Model(&Red{}).AddIndex("idx_id_valid", "id", "valid").Error
	if err != nil {
		log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	}
	err = db.GDbMgr.Get().Model(&Robber{}).AddIndex("idx_country_code_phone", "country_code", "phone").Error
	if err != nil {
		log.ZapLog().Error("AddIndex err", zap.Error(err))
	}
	err = db.GDbMgr.Get().Model(&Slogan{}).AddIndex("idx_activity_id", "activity_id").Error
	if err != nil {
		log.ZapLog().Error("AddIndex err", zap.Error(err))
	}
	//err := db.GDbMgr.Get().AutoMigrate(&table.Access{}, &table.Mod{}, &table.ModAccessRule{}, &table.Rule{}, &table.UserVip{}, &table.Vip{}, &table.VipMod{}).Error
	//if err != nil {
	//	log.ZapLog().Error("AutoMigrate err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.UserVip{}).AddUniqueIndex("idx_uservip_userkey_valid","user_key", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.VipMod{}).AddIndex("idx_vipmod_vipid_valid","vip_id",  "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.Vip{}).AddUniqueIndex("idx_vip_level_valid", "level", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.Mod{}).AddUniqueIndex("idx_mod_id_valid", "id", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.ModAccessRule{}).AddIndex("idx_modaccessrule_modid_valid","mod_id", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.Access{}).AddUniqueIndex("idx_access_id_valid", "id", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&table.Rule{}).AddUniqueIndex("idx_rule_id_valid", "id", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	log.ZapLog().Info("end InitDbTable")
}
