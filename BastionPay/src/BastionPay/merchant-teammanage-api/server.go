package main

import (
	"BastionPay/merchant-teammanage-api/config"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"

	. "BastionPay/bas-base/log/zap"
)

type WebServer struct {
	mIris   *iris.Application
	mBkIris *iris.Application
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
	this.mIris = iris.New()
	this.mIris.Use(recover.New())
	this.mIris.Use(logger.New())
	if config.GConfig.Server.Debug {
		this.mIris.Any("/debug/pprof/{action:path}", pprof.New())
	}

	this.mBkIris = iris.New()
	this.mBkIris.Use(recover.New())
	this.mBkIris.Use(logger.New())
	if config.GConfig.Server.Debug {
		this.mBkIris.Any("/debug/pprof/{action:path}", pprof.New())
	}
	//err := db.GRedis.Init(config.GConfig.Redis.Host,
	//	config.GConfig.Redis.Port, config.GConfig.Redis.Password,
	//	config.GConfig.Redis.Database)
	//if err != nil {
	//	return err
	//}

	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Run() error {
	ZapLog().Info("WebServer Run with port[" + config.GConfig.Server.Port + "]")

	//go this.bkRun()

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

func (this *WebServer) Stop() error { //这里要处理下，全部锁得再看看，还有就是qid
	return nil
}

/********************内部接口************************/
func (a *WebServer) controller() {
	a.routes()
}
