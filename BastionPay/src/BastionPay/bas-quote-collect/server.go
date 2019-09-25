package main

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	. "BastionPay/bas-quote-collect/config"
	"BastionPay/bas-quote-collect/controllers"
	"BastionPay/bas-quote-collect/quote"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
)

const (
	//	ErrCode_Success    = 0
	ErrCode_Param      = 10001
	ErrCode_InerServer = 10002
)

type controller struct {
	Controllers controllers.Controllers
}

type WebServer struct {
	mIris       *iris.Application
	mQuote      quote.QuoteMgr
	Controllers controllers.Controllers
}

func NewWebServer() *WebServer {
	web := new(WebServer)
	if err := web.Init(); err != nil {
		ZapLog().With(zap.Error(err)).Error("quote Init err")
		panic("quote Init err")
	}
	return web
}

func (this *WebServer) Init() error {
	err := this.mQuote.Init()
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("quote Init err")
		panic("quote Init err")
	}
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	if GConfig.Server.Debug {
		app.Any("/debug/pprof/{action:path}", pprof.New())
	}
	this.mIris = app
	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Run() error {
	err := this.mQuote.Start()
	if err != nil {
		ZapLog().Sugar().Error("Quote start err[%v]", err)
		panic("Quote Start err")
	}
	err = this.mIris.Run(iris.Addr(":" + GConfig.Server.Port)) //阻塞模式
	if err != nil {
		if err == iris.ErrServerClosed {
			ZapLog().Sugar().Infof("Iris Run[%d] Stoped[%v]", GConfig.Server.Port, err)
		} else {
			ZapLog().Sugar().Errorf("Iris Run[%d] err[%v]", GConfig.Server.Port, err)
		}
	}
	return nil
}

func (this *WebServer) Stop() error {
	return this.mQuote.Stop()
}

/********************内部接口************************/
func (this *WebServer) controller() {
	app := this.mIris
	//app.UseGlobal(interceptorCtrl.Interceptor)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data":    "",
		})
	})

	v1 := app.Party("/v1", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	{
		quoteCtl := controllers.NewQuoteCtl(&this.mQuote)
		codeCtl := controllers.NewCodeCtl(&this.mQuote)
		kxianCtl := controllers.NewKXianCtl(&this.mQuote)
		fileCtl := controllers.NewFileTransferCtl()
		//v1.Get("/coin/mquote", this.handleTicker) //准备废弃
		v1.Get("/coin/quote", quoteCtl.Ticker) //各种币==>各种法币
		v1.Get("/coin/code", codeCtl.ListSymbols)
		v1.Post("/coin/set", codeCtl.SetSymbol)
		v1.Get("/coin/kxian", kxianCtl.GetKXian)
		v1.Post("/filetransfer/status/add", fileCtl.AddStatus)
		v1.Post("/filetransfer/status/update", fileCtl.UpdateStatus)
		v1.Get("/filetransfer/status/get", fileCtl.GetStatus)
		v1.Any("/", this.defaultRoot)
	}
}

func (this *WebServer) defaultRoot(ctx iris.Context) {
	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	ctx.JSON(resMsg)
}
