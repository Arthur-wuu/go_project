package main

import (
	"BastionPay/merchant-teammanage-api/config"
	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/merchant-teammanage-api/db"
	"BastionPay/merchant-teammanage-api/models"
	"math/rand"
	"time"
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
	//if err := db.GRedis.Init(config.GConfig.Redis.Host, config.GConfig.Redis.Port, config.GConfig.Redis.Password, config.GConfig.Redis.Database); err != nil {
	//	return err
	//}
	models.InitDbTable()
	//db.GCache.Init()
	//db.GCache.SetActivityFunc(new(models.Activity).InnerGetByUuid)
	//db.GCache.SetShareInfoFunc(new(models.ShareInfo).InnerGetByAcId)
	//db.GCache.SetPageFunc(new(models.Page).InnerGetByAcId)
	//db.GCache.SetRobberFunc(new(models.Robber).InnerGetRedIdAndPhone)
	//db.GCache.SetPageShareInfoFunc(new(models.PageShareInfo).InnerGetByAcId)
	//db.GCache.SetSponsorAkFunc(new(models.Sponsor).InnerGetSponsorByAk)
	//db.GCache.SetFissionApiActyList(new(fisson_api.ReActivity).InnerGet)
	//db.GCache.SetActyYqlFunc()

	//db.GCache.SetUserLevelFunc(new(models.UserLevelSearch).InnerSearch)
	//db.GCache.SetLevelRuleFunc(new(models.RuleListSearch).InnerSearch)
	//pusher.GTasker.Init()
	//pusher.GTasker.Start()
	//sdk_notify_mail.GNotifySdk.Init(config.GConfig.BasNotify.Addr, "merchant-teammanage-api")
	//if err := controllers.Init(); err != nil {
	//	return err
	//}
	rand.Seed(time.Now().Unix())
	return nil
}

func UnLoader() {
	//pusher.GTasker.Stop()
}
