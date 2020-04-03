package main

import (
	"BastionPay/bas-game/base"
	"BastionPay/bas-game/type"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	//"github.com/bugsnag/bugsnag-go/errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	//"net/http"
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
	"BastionPay/bas-game/config"
	//"BastionPay/bas-filetransfer-srv/api"
	"BastionPay/bas-game/comsumer"
	"BastionPay/bas-game/db"
)

type WebServer struct {
	mIris  *iris.Application
	wsconn base.WsCon
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

	this.wsconn = base.WsCon{}
	url := "ws://iot.bigeapp.com:1883/ws"
	err = this.wsconn.Init(url, SendPingHander, RecvPingHander, RecvPongHander)
	if err != nil {
		return err
	}

	this.wsconn.Start()

	comsumer.GTasker.Init(&this.wsconn)

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
	comsumer.GTasker.Start()
	// go a.get()
}

func (this *WebServer) get() {

	for true {

		sum1 := sum()
		sumFloat1, _ := strconv.ParseFloat(sum1, 64)
		putAllMoneyToRedis(sum1)
		fmt.Println("first balance...", sumFloat1)

		time.Sleep(time.Second * time.Duration(6))

		sum2 := sum()
		sumFloat2, _ := strconv.ParseFloat(sum2, 64)
		redisAllMoney, _ := getAllMoneyToRedis()
		putAllMoneyToRedis(sum2)

		fmt.Println("second balance...", sumFloat2)
		fmt.Println("compare..", sumFloat2 > redisAllMoney, sumFloat2, redisAllMoney)
		if sumFloat2 > redisAllMoney {
			//发送ws
			message2 := []byte("{\"type\":\"admin.coinup\",\"upnumber\":\"2\",\"devid\":\"860344040771835\"}")
			this.wsconn.Send(1, message2)

			for i := 0; i < 3; i++ {
				msg, err := this.wsconn.Recv()
				if err != nil {
					log.Println("err", err)
					continue
				}
				Msg := new(_type.MsgRcv)
				json.Unmarshal(msg.Data, Msg)

				if Msg.Message == "success" && Msg.State == "coinup end" && Msg.Type == "coinup" {
					log.Println("up coin succ...")
					fmt.Println("up coin succ...")
					break
				}
			}
			fmt.Println("send over, next...")
			continue

		} else {
			fmt.Println("no paly, next...")
			continue
		}
	}
}

//var addr = flag.String("addr", "iot.bigeapp.com:1883", "http service address")
var addr = "iot.bigeapp.com:1883"

func send() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
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
	c.Close()

	done := make(chan struct{})

	func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				fmt.Println("recv  message...", string(message))
				return
			}
			log.Printf("recv: %s", message)
		}

	}()

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	//
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
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var value string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
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
		fmt.Println("----send4-----")
	}
	ZapLog().Info("get balance succ")

	return value
}

func SendPingHander() []byte {

	return []byte("ping")
}

func RecvPingHander(str string) []byte {

	return []byte(str)
}

func RecvPongHander(str string) []byte {

	return []byte(str)
}

func putAllMoneyToRedis(amount string) error {
	key := "Game_AllMoney"
	ZapLog().With(zap.String("key", string(key))).Debug("redis put")
	_, err := db.GRedis.Do("SET", key, amount)
	if err != nil {
		return err
	}
	return nil
}

func getAllMoneyToRedis() (float64, error) {
	key := "Game_AllMoney"
	ZapLog().With(zap.String("key", string(key))).Debug("redis put")
	allMoney, err := db.GRedis.Do("GET", key)
	if err != nil {
		return 0, err
	}
	sumFloat1, err := strconv.ParseFloat(string(allMoney.([]byte)), 64)
	return sumFloat1, err
}
