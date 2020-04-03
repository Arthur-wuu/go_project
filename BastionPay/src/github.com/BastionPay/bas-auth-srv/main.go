package main

import (
	//"BastionPay/bas-service/base/service"
	"BastionPay/bas-api/utils"
	"BastionPay/bas-auth-srv/db"
	"BastionPay/bas-auth-srv/handler"
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

const AuthSrvConfig = "auth.json"

func closeLog() {
	time.Sleep(time.Second * 3)
	l4g.Close()
}

func main() {
	laxFlag := config.NewLaxFlagDefault()
	cfgDir := laxFlag.String("conf_path", config.GetBastionPayConfigDir(), "config path")
	logPath := laxFlag.String("log_path", config.GetBastionPayConfigDir()+"/log.xml", "log conf path")
	laxFlag.LaxParseDefault()
	fmt.Printf("command param: conf_path=%s, log_path=%s\n", *cfgDir, *logPath)

	l4g.LoadConfiguration(*logPath)
	defer closeLog()
	defer utils.PanicPrint()

	cfgPath := *cfgDir + "/" + AuthSrvConfig
	fmt.Println("config path:", cfgPath)
	db.Init(cfgPath)

	// create service node
	nodeInstance, err := service.NewServiceNode(cfgPath)
	if nodeInstance == nil || err != nil {
		l4g.Error("Create service node failed: %s", err.Error())
		return
	}

	accountDir := *cfgDir + "/" + config.BastionPayAccountDirName
	handler.AuthInstance().Init(cfgPath, accountDir, nodeInstance)

	// register apis
	service.RegisterNodeApi(nodeInstance, handler.AuthInstance())

	// start service node
	ctx, cancel := context.WithCancel(context.Background())
	service.StartNode(ctx, nodeInstance)

	fmt.Println("Press Ctrl+c to quit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	l4g.Info("Waiting all routine quit...")
	service.StopNode(nodeInstance)
	l4g.Info("All routine is quit...")
}
