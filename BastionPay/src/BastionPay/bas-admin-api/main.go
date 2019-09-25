package main

import (
	"flag"
	"fmt"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/common"
	"github.com/BastionPay/bas-admin-api/config"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/BastionPay/bas-tools/sdk.notify.mail"
	"go.uber.org/zap"
)

var confFile = flag.String("c", "config.yml", "conf file.")

func main() {
	var (
		conf  *config.Config
		redis *common.Redis
		db    *common.Db
		err   error
	)

	flag.Parse()

	// 读取配置
	conf = new(config.Config)
	common.NewConfig(conf).Read(*confFile)

	fmt.Printf("Config[%v]\n", *conf)
	LoadZapConfig(conf.Server.Logpath)
	defer ZapClose()
	ZapLog().Sugar().Infof("Config[%s]", *conf)

	if err := sdk_notify_mail.GNotifySdk.Init(conf.Bas_notify.Addr, "bas-admin-api"); err != nil {
		ZapLog().Error("notify_sdk init err", zap.String("notifyAddr", conf.Bas_notify.Addr))
		return
	}

	bastionpay.Init(&conf.Wallet)
	ZapLog().Info("bastionpay init ok")
	// 连接redis
	redis = common.NewRedis(conf.Redis.Host, conf.Redis.Port, conf.Redis.Password, conf.Redis.Db)
	ZapLog().Info("NewRedis ok")
	// 连接DB
	db, err = common.NewDb(&common.DbOptions{
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.User,
		conf.Db.Password,
		conf.Db.DbName,
		conf.Db.MaxIdleConn,
		conf.Db.MaxOpenConn})
	if err != nil {
		ZapLog().Error("NewDb err " + err.Error())
		return
	}
	ZapLog().Info("NewDB ok")
	// 启动服务
	NewApp(conf, redis, db.GetConn())

}
