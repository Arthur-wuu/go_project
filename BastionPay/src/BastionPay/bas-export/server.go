package main

import (
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"database/sql"
	"encoding/csv"
	"fmt"
	. "github.com/BastionPay/bas-export/config"
	"github.com/BastionPay/bas-export/quote"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

var (
	tables         = make([]string, 0)
	dataSourceName = ""
)

const (
	ErrCode_Success    = 0
	ErrCode_Param      = 10001
	ErrCode_InerServer = 10002
)

type WebServer struct {
	mIris       *iris.Application
	mExportData quote.DataMgr
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
	//err := this.mExportData.Init()
	//if err != nil {
	//	ZapLog().With(zap.Error(err)).Error("quote Init err")
	//	panic("quote Init err")
	//}
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

	err := this.mIris.Run(iris.Addr(":" + GConfig.Server.Port)) //阻塞模式
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
	return this.mExportData.Stop()
}

/********************内部接口************************/
func (a *WebServer) controller() {
	app := a.mIris
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
		v1.Get("/data/export", a.handleExport) //数据导出
		v1.Any("/", a.defaultRoot)
	}
}

func (this *WebServer) defaultRoot(ctx iris.Context) {
	resMsg := apiquote.NewResMsg(ErrCode_Success, "")
	ctx.JSON(resMsg)
}

//带上appid吧，适合以后做次数限制
func (this *WebServer) handleExport(ctx iris.Context) {
	defer PanicPrint()
	ZapLog().Debug("start handleTicker ")
	table := strings.ToUpper(ctx.URLParam("table"))

	fmt.Println(ctx.Params().Len(), ctx.Values().Len())
	//判断coin是数字还是字符
	if len(table) == 0 {
		ZapLog().With(zap.String("table", table)).Error("param table name err")
		ctx.JSON(*apiquote.NewResMsg(ErrCode_Param, "param fail"))
		return
	}
	limit := strings.ToUpper(ctx.URLParam("limit"))
	//	to := ctx.Params().Get("to")
	if len(limit) == 0 {
		limit = "1000"
	}

	table = strings.TrimSpace(table)
	limit = strings.TrimSpace(limit)
	table = strings.TrimRight(table, ",")
	limit = strings.TrimRight(limit, ",")
	lim, err := strconv.Atoi(limit)
	//fromArr := strings.Split(from, ",")
	//toArr := strings.Split(to, ",")
	resMsg := apiquote.NewResMsg(ErrCode_Success, "")

	dataSourceName = this.mExportData.Init()
	tables = append(tables, strings.Split(table, ",")...)
	count := len(tables)
	ch := make(chan bool, count)

	db, err := sql.Open("mysql", dataSourceName)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	for _, table := range tables {
		go querySQL(db, table, ch, lim)
	}

	for i := 0; i < count; i++ {
		<-ch
	}
	fmt.Println("Done!")

	ctx.JSON(resMsg)
	ZapLog().Debug("deal handleTicker ok")
	fmt.Println("done over")
}

func querySQL(db *sql.DB, table string, ch chan bool, lim int) {
	fmt.Println("开始处理：", table)
	rows, err := db.Query(fmt.Sprintf("SELECT * from %s limit %d", table, lim))

	if err != nil {
		panic(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	//values：一行的所有值,把每一行的各个字段放到values中，values长度==列数
	values := make([]sql.RawBytes, len(columns))
	// print(len(values))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	//存所有行的内容totalValues
	totalValues := make([][]string, 0)
	for rows.Next() {

		//存每一行的内容
		var s []string

		//把每行的内容添加到scanArgs，也添加到了values
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for _, v := range values {
			s = append(s, string(v))
			// print(len(s))
		}
		totalValues = append(totalValues, s)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	writeToCSV(table+".csv", columns, totalValues)
	ch <- true
}

//writeToCSV
func writeToCSV(file string, columns []string, totalValues [][]string) {
	f, err := os.Create(file)
	// fmt.Println(columns)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	//f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	for i, row := range totalValues {
		//第一次写列名+第一行数据
		if i == 0 {
			w.Write(columns)
			w.Write(row)
		} else {
			w.Write(row)
		}
	}
	w.Flush()
	fmt.Println("处理完毕：", file)
}
