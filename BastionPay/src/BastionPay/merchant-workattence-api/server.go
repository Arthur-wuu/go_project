package main

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-sdk/go-sdk"
	"BastionPay/merchant-workattence-api/config"
	"BastionPay/merchant-workattence-api/jobs"
	"BastionPay/merchant-workattence-api/modules"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	"strconv"

	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/merchant-workattence-api/db"
	"BastionPay/merchant-workattence-api/models"
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
	this.controller()
	this.SetGlobalVar()
	jobs.Schd.RunCron()
	db.GCache.SetDingtalkCache()
	go this.SendSignAward()
	go this.SendStaffMotivation()
	go this.SendRClassifyAward()
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

func (this *WebServer) SendSignAward() {
	PaySdk := new(go_sdk.PaySdk)

	for c := range models.SChan {
		ZapLog().Sugar().Infof("send award params info coin[%v] symbol[%v] merchantId[%v] receiveAct[%v]", c.Coin, c.Symbol, c.MerchantId, c.AccountId)
		sendStatus, _ := PaySdk.Transfer(fmt.Sprintf("%v", c.Coin), c.Symbol, strconv.Itoa(c.MerchantId), strconv.Itoa(c.AccountId), "checkin award", "")

		//send award false
		if sendStatus == "fail" && c.Times > 0 {
			ZapLog().Sugar().Infof("send award false chan info: %v", c)
			c.Times--
			models.SChan <- c
		} else if sendStatus == "succ" { //send award success
			ZapLog().Sugar().Infof("send award success record id: %v", c.AwardId)
			err := new(models.AwardRecord).SetTransferSuccess(c.AwardId)

			if err != nil {
				ZapLog().Sugar().Errorf("send award success db update err record id: %v err[%v]", c.AwardId, err)
			}
		}
	}
}

func (this *WebServer) SendStaffMotivation() {
	PaySdk := new(go_sdk.PaySdk)

	for c := range models.SfChan {
		ZapLog().Sugar().Infof("send years of employment award params info coin[%v] symbol[%v] merchantId[%v] receiveAct[%v]", c.Coin, c.Symbol, c.MerchantId, c.AccountId)
		sendStatus, _ := PaySdk.Transfer(fmt.Sprintf("%v", c.Coin), c.Symbol, strconv.Itoa(c.MerchantId), strconv.Itoa(c.AccountId), "years of employment award", "")

		//send award false
		if sendStatus == "fail" && c.Times > 0 {
			ZapLog().Sugar().Infof("send award false chan info: %v", c)
			c.Times--
			models.SfChan <- c
		} else if sendStatus == "succ" { //send award success
			ZapLog().Sugar().Infof("send years of employment award success record id: %v", c.AwardId)
			err := new(models.StaffMotivation).SetTransferSuccess(c.AwardId)

			if err != nil {
				ZapLog().Sugar().Errorf("send years of employment award success db update err record id: %v err[%v]", c.AwardId, err)
			}
		}
	}
}

func (this *WebServer) SendRClassifyAward() {
	PaySdk := new(go_sdk.PaySdk)

	for c := range models.RcChan {
		ZapLog().Sugar().Infof("send rubbish classify award params info coin[%v] symbol[%v] merchantId[%v] receiveAct[%v]", c.Coin, c.Symbol, c.MerchantId, c.AccountId)
		sendStatus, _ := PaySdk.Transfer(fmt.Sprintf("%v", c.Coin), c.Symbol, strconv.Itoa(c.MerchantId), strconv.Itoa(c.AccountId), "rubbish classify award", "")

		//send award false
		if sendStatus == "fail" && c.Times > 0 {
			ZapLog().Sugar().Infof("send award false chan info: %v", c)
			c.Times--
			models.RcChan <- c
		} else if sendStatus == "succ" { //send award success
			ZapLog().Sugar().Infof("send rubbish classifyt award success record id: %v", c.AwardId)
			err := new(models.RubbishClassify).SetTransferSuccess(c.AwardId)

			if err != nil {
				ZapLog().Sugar().Errorf("send rubbish classify award success db update err record id: %v err[%v]", c.AwardId, err)
			}
		}
	}
}

/********************内部接口************************/
func (a *WebServer) controller() {
	a.routes()
}

func (a *WebServer) SetGlobalVar() {
	models.SChan = make(chan models.SendChan, config.GConfig.Award.ChanLen)
	models.SfChan = make(chan models.StaffChan, config.GConfig.Company.ServiceAward.ChanLen)
	models.RcChan = make(chan models.RClassifyChan, config.GConfig.Company.RubbishClassify.ChanLen)
	models.CorpIdMap = map[int]string{config.GConfig.Company.Id[0]: modules.CORPIDRUOXI, config.GConfig.Company.Id[1]: modules.CORPID}
}
