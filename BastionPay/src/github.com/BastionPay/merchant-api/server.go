package main

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/comsumer"
	"BastionPay/merchant-api/config"
	"BastionPay/merchant-api/device"
	"BastionPay/merchant-api/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/merchant-api/db"
)

type WebServer struct {
	mIris  *iris.Application
}

func NewWebServer() *WebServer {
	web := new(WebServer)
	if err := web.Init(); err != nil {
		ZapLog().With(zap.Error(err)).Error("webServer Init err")
		panic("webServer Init err")
	}
	return web
}

func (this *WebServer) Init() error {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	if config.GConfig.Server.Debug {
		app.Any("/debug/pprof/{action:path}", pprof.New())
	}
	this.mIris = app
	//err := db.GRedis.Init(config.GConfig.Redis.Host,
	//	config.GConfig.Redis.Port, config.GConfig.Redis.Password,
	//	config.GConfig.Redis.Database)
	//if err != nil {
	//	return err
	//}
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
	models.InitDbTable()
	//models.Init()
	//models.InitMerchantConfig()
	//models.InitWhiteList()
	// init 各个配置文件中的机器，可是如果有些机器关着呢，重来？
	device.GDeviceMgr.Init()
	//comsumer.GTasker.Start()

	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Run() error {
	ZapLog().Info("WebServer Run with port["+config.GConfig.Server.Port+"]")
	//go controllers.LoadVipAuth()
	err := this.mIris.Run(iris.Addr(":" + config.GConfig.Server.Port)) //阻塞模式
	if err != nil {
		if err == iris.ErrServerClosed {
			ZapLog().Sugar().Infof("Iris Run[%v] Stoped[%v]", config.GConfig.Server.Port, err)
		} else {
			ZapLog().Sugar().Errorf("Iris Run[%v] err[%v]", config.GConfig.Server.Port, err)
		}
	}
	return nil
}

func (this *WebServer) Stop() error {//这里要处理下，全部锁得再看看，还有就是qid
	return nil
}

/********************内部接口************************/
func (a *WebServer) controller() {

	comsumer.GLoginTasker.Start()
	a.routes()
}