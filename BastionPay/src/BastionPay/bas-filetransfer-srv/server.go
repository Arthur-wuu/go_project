package main

import (
	"BastionPay/bas-filetransfer-srv/config"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	//"BastionPay/bas-filetransfer-srv/api"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/consumer"
	"BastionPay/bas-filetransfer-srv/db"
	"BastionPay/bas-tools/sdk.aws.s3"
	"BastionPay/bas-tools/sdk.aws.sts"
)

type WebServer struct {
	mIris *iris.Application
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
	err := db.GRedis.Init(config.GConfig.Redis.Host,
		config.GConfig.Redis.Port, config.GConfig.Redis.Password,
		config.GConfig.Redis.Database)
	if err != nil {
		return err
	}
	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Run() error {
	ZapLog().Info("WebServer Run with port[" + config.GConfig.Server.Port + "]")
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
	db.Init()
	db.GFileS3Sdk = sdk_aws_s3.NewS3Sdk(config.GConfig.Aws.FileRegion, config.GConfig.Aws.AccessKeyId, config.GConfig.Aws.AccessKey, config.GConfig.Aws.AccessToken)
	db.GStsSdk = sdk_aws_sts.NewStsSdk(config.GConfig.Aws.FileRegion,
		config.GConfig.Aws.AccessKeyId,
		config.GConfig.Aws.AccessKey, "")
	a.routes()
	consumer.GWorker.Start()
}
