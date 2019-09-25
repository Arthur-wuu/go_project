package main

import (

	"BastionPay/pay-user-merchant-api/config"
	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/pay-user-merchant-api/db"
	"BastionPay/pay-user-merchant-api/models"
	"math/rand"
	"time"
	"BastionPay/pay-user-merchant-api/common"
)

func Loader() error {
	if err := db.GDbMgr.Init(&db.DbOptions{
		Host:        config.GConfig.Db.Host,
		Port:        config.GConfig.Db.Port,
		User:        config.GConfig.Db.User,
		Pass:        config.GConfig.Db.Pwd,
		DbName:      config.GConfig.Db.Dbname,
		MaxIdleConn: config.GConfig.Db.Max_idle_conn,
		MaxOpenConn: config.GConfig.Db.Max_open_conn,
	}); err != nil {
		return err
	}
	if err := common.GRedis.Init(config.GConfig.Redis.Host, config.GConfig.Redis.Port, config.GConfig.Redis.Password, config.GConfig.Redis.Database); err != nil {
		return err
	}
	models.InitDbTable()
	db.GCache.Init()

	//sdk_notify_mail.GNotifySdk.Init(config.GConfig.BasNotify.Addr, "pay-user-merchant-api")
	//if err := controllers.Init(); err != nil {
	//	return err
	//}
	rand.Seed(time.Now().Unix())
	return nil
}

func UnLoader() {

}
