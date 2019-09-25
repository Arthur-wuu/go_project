package main

import (
"BastionPay/bas-base/config"
. "BastionPay/bas-base/log/zap"
. "BastionPay/merchant-api/config"
"fmt"
"go.uber.org/zap"
"runtime/debug"
)

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

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

	srv := NewWebServer()
	ZapLog().Sugar().Info("WebServer Start Runing...")
	srv.Run() //阻塞
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//	<-c
	ZapLog().Sugar().Info("WebServer Stop")
}
