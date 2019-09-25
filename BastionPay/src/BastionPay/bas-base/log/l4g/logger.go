package l4gmgr

import (
	l4g "github.com/alecthomas/log4go"
	"time"
)

func LoadConfig(path string) {
	l4g.LoadConfiguration(path)
	GlobalLogMgr.Start()
}

func Close() {
	GlobalLogMgr.Stop()
	time.Sleep(time.Second * 2)
	l4g.Close()
}
