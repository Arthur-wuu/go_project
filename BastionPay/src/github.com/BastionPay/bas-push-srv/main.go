package main

import (
	//"github.com/BastionPay/bas-service/base/service"
	service "github.com/BastionPay/bas-base/service2"
	"BastionPay/bas-push-srv/handler"
	"fmt"
	"context"
	"time"
	l4g "github.com/alecthomas/log4go"
	"BastionPay/bas-push-srv/db"
	"github.com/BastionPay/bas-base/config"
	"github.com/BastionPay/bas-api/utils"
	"os"
	"os/signal"
)

const PushSrvConfig = "push.json"

func closeLog()  {
	time.Sleep(time.Second * 3)
	defer l4g.Close()
}

func main() {
	cfgDir := config.GetBastionPayConfigDir()

	l4g.LoadConfiguration(cfgDir + "/log.xml")
	defer closeLog()

	defer utils.PanicPrint()

	cfgPath := cfgDir + "/" + PushSrvConfig
	db.Init(cfgPath)

	accountDir := cfgDir + "/" + config.BastionPayAccountDirName
	handler.PushInstance().Init(cfgPath, accountDir)

	// create service node
	fmt.Println("config path:", cfgPath)
	nodeInstance, err := service.NewServiceNode(cfgPath)
	if nodeInstance == nil || err != nil{
		l4g.Error("Create service node failed: %s", err.Error())
		return
	}

	// register apis
	service.RegisterNodeApi(nodeInstance, handler.PushInstance())

	// start service node
	ctx, cancel := context.WithCancel(context.Background())
	service.StartNode(ctx, nodeInstance)

	fmt.Println("Press Ctrl+c to quit...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	cancel()
	l4g.Info("Waiting all routine quit...")
	service.StopNode(nodeInstance)
	l4g.Info("All routine is quit...")
}