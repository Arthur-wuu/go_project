package main

import (
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common/logrus"
	"github.com/BastionPay/bas-bkadmin-api/bastionpay"
	"github.com/BastionPay/bas-bkadmin-api/models"
	rd "github.com/BastionPay/bas-bkadmin-api/models/redis"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris"
	"github.com/robfig/cron"
)

const LOGFILENAME = "backend_serve_error.log"

func (this *Service) Run() {
	address := fmt.Sprintf("%s:%d",
		this.Config.System.Host,
		this.Config.System.Port)

	bastionpay.Init(&this.Config.Wallet)

	orm := models.New(&this.Config.Mysql).Connection()

	this.RedisClient = rd.New(&this.Config.Redis).Connection()

	iris.RegisterOnInterrupt(func() {
		orm.Close()
		this.RedisClient.Close()
	})

	this.routes()

	if err := this.App.Run(iris.Addr(address)); err != nil {
		l4g.Info("run addr[%s] err[%v]", address, err)
	}
}

func (this *Service) RunLogrus() {
	logrus.New(
		this.Config.System.LogPath,
		LOGFILENAME,
		this.Config.System.Debug)
}

func (this *Service) timer() {
	cronds := cron.New()

	cronds.AddFunc("1 * * * * *", func() {
		fmt.Println("test")
	})

	cronds.Start()
}
