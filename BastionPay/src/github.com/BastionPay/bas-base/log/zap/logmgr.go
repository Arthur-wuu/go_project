package log

import (
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

const MIN_AfterTime = 5 //至少5秒，不然没意义

type LogLevelCallBack func(newLv zapcore.Level) error

var GlobalLogMgr LogMgr

type LogMgr struct {
	sync.Mutex
	mTimer     *time.Timer
	mLastLV    zapcore.Level //记录最新日志级别
	mLastLvStr string
	//	mIgnoreDev string               //忽略设备
	mLogCallBack LogLevelCallBack //自定义日志设置方法
	mExitCH      chan bool        //整个系统的关闭
	mWaitGroup   sync.WaitGroup
}

func (this *LogMgr) Start() {
	if this.mExitCH != nil {
		return
	}
	this.mLastLV = -1
	this.mExitCH = make(chan bool)
	this.mTimer = time.NewTimer(time.Hour * 10)
	this.logRecover()
}

func (this *LogMgr) Stop() {
	if this.mExitCH == nil {
		return
	}
	close(this.mExitCH)
	this.mWaitGroup.Wait()
}

func (this *LogMgr) ResetLevel() error {
	if this.mExitCH == nil {
		return nil
	}
	this.stopLogTimer()
	this.Lock()
	defer this.Unlock()
	if err := this.setFilt(zap.InfoLevel); err != nil {
		return err
	}
	this.mLastLV = zap.InfoLevel
	this.mLastLvStr = "INFO"
	return nil
}

func (this *LogMgr) SetLevel(levelStr string, timeout int64) error {
	if this.mExitCH == nil {
		return nil
	}
	levelCode, ok := ZapLogLevelMap[levelStr]
	if !ok {
		return fmt.Errorf("unknown level[%s]", levelStr)
	}

	if timeout <= MIN_AfterTime {
		return fmt.Errorf("too small timeout[%d]", timeout)
	}

	this.stopLogTimer()
	//	l4g.Info("setLevel[%s][%d][%d]", levelStr, expireAt, afterTime)
	this.Lock()
	defer this.Unlock()
	if levelCode == this.mLastLV {
		l4g.Info("same as last level[%s]", levelStr)
		this.setLogTimer(levelCode, timeout)
		return nil
	}

	if err := this.setFilt(levelCode); err != nil {
		return err
	}
	this.mLastLV = levelCode
	this.mLastLvStr = levelStr
	this.setLogTimer(levelCode, timeout)
	return nil
}

//设置第三方日志级别函数
func (this *LogMgr) SetLogLevelFunc(pFun LogLevelCallBack) error {
	this.mLogCallBack = pFun
	return nil
}

/*****************************内部接口********************************/
func (this *LogMgr) setLogTimer(levelCode zapcore.Level, afterTime int64) {
	if levelCode >= zap.InfoLevel || afterTime <= MIN_AfterTime {
		return
	}
	//	l4g.Info("setLogTimer[%d][%d]s", levelCode, afterTime)
	this.mTimer.Reset(time.Second * time.Duration(afterTime))
}

func (this *LogMgr) stopLogTimer() {
	this.mTimer.Stop()
	select {
	case <-this.mTimer.C:
	default:
	}
}

//日志定时恢复，Info以下级别定时恢复到Info
func (this *LogMgr) logRecover() {
	go func() {
		this.stopLogTimer()
		this.mWaitGroup.Add(1)
		defer this.mWaitGroup.Done()
		for {
			select {
			case <-this.mTimer.C:
				//				l4g.Info("time up to recover to info")
				this.Lock()
				if this.mLastLV >= zap.InfoLevel {
					this.Unlock()
					continue
				}
				if err := this.setFilt(zap.InfoLevel); err != nil {
					l4g.Error("setFilt err[%s]", err.Error())
					this.Unlock()
					continue
				}
				l4g.Info("last loglevel[%s] expired, set to Info", this.mLastLvStr)
				this.mLastLV = zap.InfoLevel
				this.mLastLvStr = "INFO"
				this.Unlock()
			case <-this.mExitCH:
				return
			}
		}

	}()
}

//设置所有设备的日志级别(不包括忽略设备)，并且调动第三方日志设置函数
func (this *LogMgr) setFilt(levelCode zapcore.Level) error {
	zapCfg.Level.SetLevel(levelCode)

	l4g.Info("Set Loglevel of Filt[%v] dev num[%d] ok", l4g.Global, len(l4g.Global))
	if this.mLogCallBack != nil {
		if err := this.mLogCallBack(levelCode); err != nil {
			return errors.New("LogLevelCallBack failed, " + err.Error())
		}
	}
	return nil
}
