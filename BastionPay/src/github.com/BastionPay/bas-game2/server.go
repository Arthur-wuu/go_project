package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"strconv"
	"time"
	//"github.com/kataras/iris/websocket"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"

	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	//"go.uber.org/zap"
	"BastionPay/bas-game2/config"
	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/bas-game2/db"
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

	go a.get()
}

func (this *WebServer) get() {

	for true {

		//userDetail1 := new(table.UserDetail)
		//err := db.GDbMgr.Get().Model(&table.UserDetail{}).Where("USER_ID=?", 35).Order("ORDER_ID desc").Last(userDetail1).Error

		sum1 := sum()
		sumFloat1, _ := strconv.ParseFloat(sum1, 64)

		time.Sleep(time.Second * time.Duration(6))

		sum2 := sum()
		sumFloat2, _ := strconv.ParseFloat(sum2, 64)

		if sumFloat2 > sumFloat1 {
			//发送ws
			c()
			fmt.Println("send over, next...")
			continue

		} else {
			fmt.Println("no paly, next...")
			continue
		}
	}
}

var addr = flag.String("addr", "iot.bigeapp.com:1883", "http service address")

func c() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	message2 := []byte("{\"type\":\"admin.coinup\",\"upnumber\":\"2\",\"devid\":\"860344040771835\"}")
	err = c.WriteMessage(1, message2)

	if err != nil {
		log.Fatal("message err:", err)
	}
	defer c.Close()

	//done := make(chan struct{})
	//
	// func() {
	//	defer close(done)
	//	for {
	//		_, message, err := c.ReadMessage()
	//		if err != nil {
	//			log.Println("read:", err)
	//			return
	//		}
	//		log.Printf("recv: %s", message)
	//	}
	//}()

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	//for {
	//	select {
	//	case <-done:
	//		return
	//	case t := <-ticker.C:
	//		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	//		if err != nil {
	//			log.Println("write:", err)
	//			return
	//		}
	//	case <-interrupt:
	//		log.Println("interrupt")
	//
	//		// Cleanly close the connection by sending a close message and then
	//		// waiting (with timeout) for the server to close the connection.
	//		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			log.Println("write close:", err)
	//			return
	//		}
	//		select {
	//		case <-done:
	//		case <-time.After(time.Second):
	//		}
	//		return
	//	}
	//}
}

func sum() string {

	rows, err := db.GDbMgr.Get().Table("USER_ACC").Select("sum(BALANCE) as total").Where("USER_ID=?", 35).Group("USER_ID").Rows()
	if err != nil {
		fmt.Println("err,***", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var value string
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "nil"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	return value
}
