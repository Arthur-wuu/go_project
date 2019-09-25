package main

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote/collect"
	. "BastionPay/bas-quote/config"
	"BastionPay/bas-quote/utils"
	"BastionPay/bas-quote/controllers"
	"BastionPay/bas-quote/quote"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	"strings"
)


type WebServer struct {
	mIris  *iris.Application
	mQuote quote.QuoteMgr
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

	block := controllers.Inspection{}
	v1 := app.Party("/api/v1", crs ,block.VerifyIsBlock ).AllowMethods(iris.MethodOptions)
	{
		quoteCtl := controllers.NewQuoteCtl(&a.mQuote)
		fxCtl := controllers.NewFxCtl(&a.mQuote)
		codeCtl := controllers.NewCodeCtl(&a.mQuote)
		kxianCtl := controllers.NewKXianCtl(&a.mQuote)
		v1.Get("/coin/quote", quoteCtl.Ticker)     //各种币==>各种法币
		v1.Get("/coin/code", codeCtl.ListSymbols)
		v1.Get("/coin/kxian", kxianCtl.GetKXian)       // 新 kxain 各种币到各种法币,和几个主流币"BTC,ETH,XRP,LTC,BCH"
		v1.Get("/coin/huilv",a.handleHuilv)        // 查看 法币 到 法币
		v1.Get("/coin/fx",fxCtl.Ticker)        // 查看 法币 到 法币
		v1.Get("/coin/list",codeCtl.ListCoinAndFx)        // list数据货币和法币列表
		v1.Get("/coin/exchange", quoteCtl.Exchange)  // 法币到 数字货币
		v1.Any("/", a.defaultRoot)
		//v1.Get("/coin/_quote", a.handleTicker2)    // 原 quote 改了个名字
		//v1.Get("/coin/quote", a.handleTicker)      // 新 quote 提供查询借口，各种币到各种法币
		//v1.Get("/coin/_kxian", a.handleKXian2)     // 原 kxian 改了个名字
	}
}

func (this *WebServer) defaultRoot(ctx iris.Context) {
	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")
	ctx.JSON(resMsg)
}


//法币到法币
func (this *WebServer) handleHuilv(ctx iris.Context) {
	defer PanicPrint()
	ZapLog().Debug("start handleHuilv ")

	//判断coin是数字还是字符
	to := strings.ToUpper(ctx.URLParam("to"))

	if len(to) == 0 {
		ZapLog().With(zap.String("to", to)).Error("param err")
		ctx.JSON(*apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc()))
		return
	}

	to = strings.TrimSpace(to)
	to = strings.TrimRight(to, ",")
	toArr := strings.Split(to, ",")
	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")

		QuoteDetailInfo := resMsg.GenQuoteDetailInfo()

			for j:=0; j < len(toArr); j++{
				if len(toArr[j]) == 0 {
					continue
				}
				moneyInfo := new(collect.MoneyInfo)
				moneyInfo3, err := this.mQuote.GetQuoteHuilv( toArr[j])
				if err != nil {
					ZapLog().Error("get qt_USD_to[i] err", zap.Error(err), zap.String("huilv", toArr[j]))
					continue
				}

				moneyInfo.SetPrice((moneyInfo3.GetPrice()))
				moneyInfo.SetSymbol(toArr[j])
				moneyInfo.SetLast_updated(moneyInfo3.GetLast_updated())
				QuoteDetailInfo.AddMoneyInfo(utils.ToApiMoneyInfo(moneyInfo))
			}
	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleTicker ok")
}


