package main

import (
	//"BastionPay/bas-service/base/service"
	"BastionPay/bas-api/utils"
	"BastionPay/bas-base/config"
	service "BastionPay/bas-base/service2"
	"context"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ServiceGatewayConfig = "gateway.json"

func closeLog() {
	time.Sleep(time.Second * 3)
	l4g.Close()
}

func main() {
	laxFlag := config.NewLaxFlagDefault()
	cfgDir := laxFlag.String("conf_path", config.GetBastionPayConfigDir(), "config path")
	logPath := laxFlag.String("log_path", config.GetBastionPayConfigDir()+"/log.xml", "log conf path")
	laxFlag.LaxParseDefault()
	fmt.Printf("commandline param: conf_path=%s, log_path=%s\n", *cfgDir, *logPath)

	l4g.LoadConfiguration(*logPath)
	defer closeLog()
	defer utils.PanicPrint()

	cfgPath := *cfgDir + "/" + ServiceGatewayConfig
	fmt.Println("config path:", cfgPath)

	// create service center
	gatewayInstance, err := service.NewServiceGateway(cfgPath)
	if gatewayInstance == nil || err != nil {
		l4g.Error("Create service center failed: %s", err.Error())
		return
	}

	// start service center
	ctx, cancel := context.WithCancel(context.Background())
	service.StartCenter(ctx, gatewayInstance)

	fmt.Println("Press Ctrl+c to quit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	l4g.Info("Waiting all routine quit...")
	service.StopCenter(gatewayInstance)
	l4g.Info("All routine is quit...")
}
