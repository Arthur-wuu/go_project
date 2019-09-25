package logrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type (
	Params map[string]interface{}

	logConfig struct {
		fileName string
		logPath  string
		Debug    bool
	}
)

var (
	logruser = logrus.New()
	path     = "/data/logs/golang/"
	fileName = "kyc_error.log"
	Debug    bool
)

func New(path, name string, debug bool) *logrus.Logger {
	config := &logConfig{
		logPath:  path,
		fileName: name,
		Debug:    debug,
	}

	config.setConfig()

	return logruser
}

func (log *logConfig) setConfig() {
	if log.logPath != "" {
		path = log.logPath
	}

	Debug = log.Debug

	_, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Fatal(err, "log module mkdir dir errer")
		}
	}

	if log.fileName == "" {
		log.fileName = fileName
	}

	fileNames := fmt.Sprintf("%s/%s", path, log.fileName)
	file, err := os.OpenFile(fileNames, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		logruser.Out = file
	} else {
		Fatal(err, "log module open file errer")
	}

	logruser.Level = logrus.InfoLevel
	logruser.Formatter = &logrus.JSONFormatter{}
}

func Info(err error, message string) {
	if Debug == true {
		logruser.WithFields(logrus.Fields{
			"errors": err,
		}).Info(message)
	}
}

func Warn(err error, message string) {
	logruser.WithFields(logrus.Fields{
		"errors": err,
	}).Warn(message)
}

func Fatal(err error, message string) {
	logruser.WithFields(logrus.Fields{
		"errors": err,
	}).Fatal(message)
}
