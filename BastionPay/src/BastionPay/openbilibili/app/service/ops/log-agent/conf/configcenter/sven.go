package configcenter

import (
	"go-common/library/conf"
	"go-common/library/log"
)

var (
	// Conf conf
	Client  *conf.Client
	Version int
)

func InitConfigCenter() {
	var err error
	if Client, err = conf.New(); err != nil {
		panic(err)
	}

	// watch update and update Version
	Client.WatchAll()
	go func() {
		for range Client.Event() {
			log.Info("config reload")
			Version += 1
		}
	}()
}
