package models

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/db"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

func InitDbTable() {
	log.ZapLog().Info("start InitDbTable")
	if !config.GConfig.Db.Debug {
		log.ZapLog().Info("end InitDbTable")
		return
	}
	err := db.GDbMgr.Get().AutoMigrate(&PosTrade{}, &ReFund{}).Error
	if err != nil {
		log.ZapLog().Error("AutoMigrate err", zap.Error(err))
	}
	//err = db.GDbMgr.Get().Model(&Trade{}).AddUniqueIndex("idx_merchant_trade_no","merchant_trade_no", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
	//err = db.GDbMgr.Get().Model(&FundOut{}).AddUniqueIndex("idx_userlevel_userkey_valid","user_key", "valid").Error
	//if err != nil {
	//	log.ZapLog().Error("AddUniqueIndex err", zap.Error(err))
	//}
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
