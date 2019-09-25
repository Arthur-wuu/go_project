package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
)

var zapCfg zap.Config
var GobalZapLog *zap.Logger //log rotation must use another progress,such as logrotate

var ZapLogLevelMap = map[string]zapcore.Level{
	"DEBUG":   zap.DebugLevel,
	"debug":   zap.DebugLevel,
	"DEBG":    zap.DebugLevel,
	"debg":    zap.DebugLevel,
	"INFO":    zap.InfoLevel,
	"info":    zap.InfoLevel,
	"WARN":    zap.WarnLevel,
	"warn":    zap.WarnLevel,
	"WARNING": zap.WarnLevel,
	"warning": zap.WarnLevel,
	"EROR":    zap.ErrorLevel,
	"eror":    zap.ErrorLevel,
	"ERROR":   zap.ErrorLevel,
	"error":   zap.ErrorLevel,
}

func LoadZapConfig(filename string) {
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %s for reading: %v\n", filename, err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %s: %v\n", filename, err)
		os.Exit(1)
	}

	if err := json.Unmarshal(contents, &zapCfg); err != nil {
		panic(err)
	}

	GobalZapLog, err = zapCfg.Build()
	if err != nil {
		panic(err)
	}
	GlobalLogMgr.Start()
}

func ZapClose() {
	if GobalZapLog == nil {
		return
	}
	GobalZapLog.Sync()
	GlobalLogMgr.Stop()
}

func ZapLog() *zap.Logger {
	return GobalZapLog
}

func SetLevel(lvStr string) error {
	if GobalZapLog == nil {
		return nil
	}
	lvCode, ok := ZapLogLevelMap[lvStr]
	if !ok {
		return errors.New("unknow " + lvStr)
	}
	zapCfg.Level.SetLevel(lvCode)
	return nil
}

/****************************************************************/
