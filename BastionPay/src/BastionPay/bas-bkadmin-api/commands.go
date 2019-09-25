package main

import (
	"github.com/urfave/cli"
)

type Command struct {
	Cli         *cli.App
	ConfigPath  string
	RbacConfig  string
	LogConfPath string
}

func (com *Command) NewCli() *Command {
	com.Cli = cli.NewApp()
	com.Cli.Name = "admin"
	com.Cli.Usage = "admin server"
	com.Cli.Version = "1.0.0"

	com.Cli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       string("config/local/config.yaml"),
			Destination: &com.ConfigPath,
			Usage:       "config file",
		},
		cli.StringFlag{
			Name:        "rbac-config,r",
			Value:       string("config/local/rbac_model.yaml"),
			Destination: &com.RbacConfig,
			Usage:       "rbac model config file",
		},
		cli.StringFlag{
			Name:        "log,l",
			Value:       string("config/local/log.xml"),
			Destination: &com.LogConfPath,
			Usage:       "log conf path",
		},
	}

	return com
}
