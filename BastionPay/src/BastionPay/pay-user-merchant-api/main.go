package main

import (
	"BastionPay/bas-base/config"
	. "BastionPay/bas-base/log/zap"
	. "BastionPay/pay-user-merchant-api/config"
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

//var confFile = flag.String("c", "config.yml", "conf file.")

//var (
//	Meta = ""
//)

func main() {
	laxFlag := config.NewLaxFlagDefault()
	cfgPath := laxFlag.String("conf_path", "config.yaml", "config path")
	logPath := laxFlag.String("log_path", "zap.conf", "log conf path")
	laxFlag.LaxParseDefault()
	fmt.Printf("command param: conf_path=%s, log_path=%s\n", *cfgPath, *logPath)
	LoadConfig(*cfgPath)
	LoadZapConfig(*logPath)
	ZapLog().Sugar().Infof("Config Content[%v]", GConfig)
	defer ZapClose()
	defer PanicPrint()

	if err := Loader(); err != nil {
		ZapLog().Sugar().Error("Loader panic....%v", err)
		return
	}

	srv := NewWebServer()
	ZapLog().Sugar().Info("WebServer Start Runing...")
	srv.Start() //阻塞
	time.Sleep(time.Second * 2)
	srv.Stop()
	UnLoader()
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//	<-c
	ZapLog().Sugar().Info("WebServer Stop")
}
