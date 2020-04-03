package main

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/api"
	"BastionPay/bas-tv-proxy/base"
	"BastionPay/bas-tv-proxy/common"
	"BastionPay/bas-tv-proxy/config"
	"BastionPay/bas-tv-proxy/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/websocket"
	"go.uber.org/zap"
)

const (
//ErrCode_Success    = 0
//ErrCode_Param      = 10001
//ErrCode_InerServer = 10002
)

type WebServer struct {
	mIris *iris.Application
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
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	if config.GConfig.Server.Debug {
		app.Any("/debug/pprof/{action:path}", pprof.New())
	}
	this.mIris = app
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
	app := a.mIris
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

	/********Ctrl**********/ //有依赖关系的，GKBSpirit最后初始化
	controllers.GCoinMerit.Init(&config.GConfig)
	controllers.GBtcExa.Init(&config.GConfig)
	controllers.GKBSpirit.Init()

	/********http***********/
	v1 := app.Party("/api/v1/tvproxy", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	{
		v1.Get("/coinmerit/quote/kxian", controllers.GCoinMerit.HandleHttpKXian)
		v1.Get("/coinmerit/quote/market", controllers.GCoinMerit.HandleHttpExa)
		v1.Get("/coinmerit/quote/objs", controllers.GCoinMerit.HandleHttpObjList)
		v1.Get("/btcexa/quote/kxian", controllers.GBtcExa.HandleHttpKXian)
		v1.Get("/btcexa/quote/objs", controllers.GBtcExa.HandleHttpObjList)
		v1.Get("/quote/kxian", controllers.GQuoteCtrl.HandleHttpKXian)
		v1.Get("/quote/kbspirit", controllers.GKBSpirit.HandleGetObjs)
		v1.Any("/", a.defaultRoot)
	}

	//Ws 设置
	Ws := websocket.New(websocket.Config{})
	app.Get("/api/v1/tvproxy/ws", Ws.Handler())
	base.GWsServerMgr.Init(&config.GConfig)
	base.GWsServerMgr.RegPackHander(common.PackPushMsg, common.PackResMsg)
	base.GWsServerMgr.RegNewRequesterHandler(common.NewRequester)
	Ws.OnConnection(base.GWsServerMgr.HandleWsConnection)

	base.GWsServerMgr.RegHandler("/api/v1/tvproxy/coinmerit/quote/kxian", controllers.GCoinMerit.HandleWsKXian)
	base.GWsServerMgr.RegHandler("/api/v1/tvproxy/btcexa/quote/kxian", controllers.GBtcExa.HandleWsKXian)
	base.GWsServerMgr.RegHandler("/api/v1/tvproxy/quote/kxian", controllers.GQuoteCtrl.HandleWsKXian)
}

func (this *WebServer) defaultRoot(ctx iris.Context) {
	resMsg := api.NewResponse("")
	ctx.JSON(resMsg)
}
