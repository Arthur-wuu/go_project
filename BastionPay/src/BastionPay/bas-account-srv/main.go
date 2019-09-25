package main

import (
	//"BastionPay/bas-service/base/service"
	"BastionPay/bas-account-srv/db"
	"BastionPay/bas-account-srv/handler"
	"BastionPay/bas-account-srv/install"
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

const AccountSrvConfig = "account.json"

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

	cfgPath := *cfgDir + "/" + AccountSrvConfig
	db.Init(cfgPath)
	InitBasNotify(cfgPath)

	accountDir := *cfgDir + "/" + config.BastionPayAccountDirName
	err := os.MkdirAll(accountDir, os.ModePerm)
	if err != nil && os.IsExist(err) == false {
		l4g.Error("Create dir failed：%s - %s", accountDir, err.Error())
		return
	}

	err = install.InstallBastionPay(accountDir)
	if err != nil {
		l4g.Error("Install super wallet failed: %s", err.Error())
		return
	}

	// create service node
	l4g.Info("config path: %s", cfgPath)
	nodeInstance, err := service.NewServiceNode(cfgPath)
	if nodeInstance == nil || err != nil {
		l4g.Error("Create service node failed: %s", err.Error())
		return
	}

	// init
	handler.AccountInstance().Init(accountDir, nodeInstance, GetAuditeTemplateName(cfgPath), GetUserFrozen1ForAdminTemp(cfgPath), GetUserFrozen1ForUserTemp(cfgPath))

	// register APIs
	service.RegisterNodeApi(nodeInstance, handler.AccountInstance())

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
