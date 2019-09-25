package main

import (
	basel4g "BastionPay/bas-base/log/l4g"
	basutils "BastionPay/bas-base/utils"
	"BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	l4g "github.com/alecthomas/log4go"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	//laxFlag := config.NewLaxFlagDefault()
	//cfgDir := laxFlag.String("conf_path", config.GetBastionPayConfigDir(), "config path")
	//logPath := laxFlag.String("log_path", config.GetBastionPayConfigDir()+"/log.xml", "log conf path")
	//laxFlag.LaxParseDefault()
	//fmt.Printf("commandline param: conf_path=%s, log_path=%s\n", *cfgDir, *logPath)

	command := Command{}
	com := command.NewCli()
	com.Cli.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "run Commands",
			Action: func(c *cli.Context) {
				conf := tools.Analyze(command.ConfigPath)

				svc := Service{
					Config: conf,
				}
				basel4g.LoadConfig(command.LogConfPath)
				l4g.Info("conf[%v]", *conf)
				basutils.GlobalMonitor.Start(conf.Monitor.Addr)
				models.GlobalNotifyMgr.Init(conf)
				models.GlobalNotifyMgr.Start()
				//			svc.RunLogrus()
				l4g.Info("start svc...")
				svc.Run()
			},
		},
	}

	defer basel4g.Close()
	defer models.GlobalNotifyMgr.Close()

	if err := com.Cli.Run(os.Args); err != nil {
		l4g.Error(err, "run server errors")
	}
	l4g.Info("Stoped.....")
	time.Sleep(time.Second * 1)
}
