package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	"BastionPay/pay-user-merchant-api/config"

	. "BastionPay/bas-base/log/zap"
	"sync"
)

type WebServer struct {
	mIris  *iris.Application
	mBkIris  *iris.Application
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

	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Start() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		this.run()
	}()
	go func() {
		defer wg.Done()
		this.bkRun()
	}()
	wg.Wait()
	return nil
}

func (this *WebServer) run() error {
	ZapLog().Info("WebServer Run with port["+config.GConfig.Server.Port+"]")

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

func (this *WebServer) bkRun() error {
	ZapLog().Info("WebServer BkRun with port["+config.GConfig.Server.BkPort+"]")
	//go controllers.LoadVipAuth()
	err := this.mBkIris.Run(iris.Addr(":" + config.GConfig.Server.BkPort)) //阻塞模式
	if err != nil {
		if err == iris.ErrServerClosed {
			ZapLog().Sugar().Infof("Iris BkRun[%v] Stoped[%v]", config.GConfig.Server.Port, err)
		} else {
			ZapLog().Sugar().Errorf("Iris BkRun[%v] err[%v]", config.GConfig.Server.Port, err)
		}
	}
	return nil
}

func (this *WebServer) Stop() error {//这里要处理下，全部锁得再看看，还有就是qid
	return nil
}

/********************内部接口************************/
func (a *WebServer) controller() {
	a.routes()
	a.bkroutes()
}