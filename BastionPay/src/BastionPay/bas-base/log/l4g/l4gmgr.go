package l4gmgr

import (
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"sync"
	"time"
)

var LogLevelMap = map[string]l4g.Level{
	"FINEST":   l4g.FINEST,
	"FINE":     l4g.FINE,
	"DEBUG":    l4g.DEBUG,
	"TRACE":    l4g.TRACE,
	"INFO":     l4g.INFO,
	"WARNING":  l4g.WARNING,
	"ERROR":    l4g.ERROR,
	"CRITICAL": l4g.CRITICAL,

	"FNST": l4g.FINEST,
	"DEBG": l4g.DEBUG,
	"TRAC": l4g.TRACE,
	"WARN": l4g.WARNING,
	"EROR": l4g.ERROR,
	"CRIT": l4g.CRITICAL,
}

const MIN_AfterTime = 5 //至少5秒，不然没意义

type LogLevelCallBack func(newLv l4g.Level) error

var GlobalLogMgr LogMgr

type LogMgr struct {
	sync.Mutex
	mTimer     *time.Timer
	mLastLV    l4g.Level //记录上次日志级别
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
	close(this.mExitCH)
	this.mWaitGroup.Wait()
}

func (this *LogMgr) ResetLevel() error {
	this.stopLogTimer()
	this.Lock()
	defer this.Unlock()
	if err := this.setFilt(l4g.INFO); err != nil {
		return err
	}
	this.mLastLV = l4g.INFO
	this.mLastLvStr = "INFO"
	return nil
}

func (this *LogMgr) SetLevel(levelStr string, timeout int64) error {
	levelCode, ok := LogLevelMap[levelStr]
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
func (this *LogMgr) setLogTimer(levelCode l4g.Level, afterTime int64) {
	if levelCode >= l4g.INFO || afterTime <= MIN_AfterTime {
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
				if this.mLastLV >= l4g.INFO {
					this.Unlock()
					continue
				}
				if err := this.setFilt(l4g.INFO); err != nil {
					l4g.Error("setFilt err[%s]", err.Error())
					this.Unlock()
					continue
				}
				l4g.Info("last loglevel[%s] expired, set to Info", this.mLastLvStr)
				this.mLastLV = l4g.INFO
				this.mLastLvStr = "INFO"
				this.Unlock()
			case <-this.mExitCH:
				return
			}
		}

	}()
}

//设置所有设备的日志级别(不包括忽略设备)，并且调动第三方日志设置函数
func (this *LogMgr) setFilt(levelCode l4g.Level) error {
	for _, filt := range l4g.Global {
		filt.Level = levelCode
	}

	l4g.Info("Set Loglevel of Filt[%v] dev num[%d] ok", l4g.Global, len(l4g.Global))
	if this.mLogCallBack != nil {
		if err := this.mLogCallBack(levelCode); err != nil {
			return errors.New("LogLevelCallBack failed, " + err.Error())
		}
	}
	return nil
}

/**************************************其他方式加载日志配置（已废弃）************************************/
/*
type SConfig struct {
	XMLName   xml.Name   `xml:"logging"` // 指定最外层的标签为config
	Filters   []SFilter  `xml:"filter"` // 读取receivers标签下的内容，以结构方式获取
}

type SFilter struct {
	XMLName   xml.Name     `xml:"filter"`
	Enabled   bool         `xml:"enabled,attr"`
	Tag       string       `xml:"tag"`
	Type      string       `xml:"type"`
	Level     string       `xml:"level"`
	Propertys []SProperty  `xml:"property"`
}

type SProperty struct {
	Name      string       `xml:"name,attr"`
	Property  string       `xml:",chardata"`
}

func makeDir(path string) error {
	index := strings.LastIndex(path,"/")
	if index < 0 {
		return nil
	}

	path = path[0:index]
	fmt.Println(path)
	cmd := exec.Command("/bin/bash", "-c", "mkdir -p "+path)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}*/

/*
* 参数说明：
*          confPath是log.xml配置文件完整路径；
*          logPath 是 日志文件输出完整路径
* 其他问题：只试过macos和cetos系统，其它系统（windows）可能不支持
 */
/*
func LoadConfiguration(confPath, logPath string){
	if err := makeDir(logPath); err != nil {
		fmt.Fprintf(os.Stderr,"makedir %s err : %v",logPath, err)
		os.Exit(1)
	}
	mysys := runtime.GOOS
	cmdStr := ""
	if mysys == "darwin" {
		cmdStr = "sed -i '' 's:test.log:"+logPath+":g' " + confPath
	} else {
		cmdStr = "sed -i 's:test.log:"+logPath+":g' " + confPath
	}
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr,"cmd[%s] err[%v]",cmdStr, err)
		os.Exit(1)
	}
	l4g.LoadConfiguration(confPath)
}*/

/*
* 参数说明：
*          confPath是log.xml配置文件完整路径；
*          logPath 是 日志文件输出完整路径
* 其他问题：
*          配置文件注释将丢失
 */
/*
func LoadConfiguration2(confPath, logPath string){
	if err := makeDir(logPath); err != nil {
		fmt.Fprintf(os.Stderr,"makedir %s err : %v",logPath, err)
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Fprintf(os.Stderr,"read %s err : %v",confPath, err)
		os.Exit(1)
	}

	var sConfig SConfig
	err = xml.Unmarshal(content, &sConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr,"Unmarshal error: %v", err)
		os.Exit(1)
	}

	for i:=0; i<len(sConfig.Filters);i++ {
		f := &sConfig.Filters[i]
		if !f.Enabled || (f.Tag != "file") || (f.Type != "file"){
			continue
		}
		for j:=0; j < len(f.Propertys); j++ {
			if f.Propertys[j].Name == "filename" {
				f.Propertys[j].Property = logPath
				break
			}
		}

	}

	xmlOutPut, outPutErr := xml.MarshalIndent(sConfig, "", "\t")
	if outPutErr != nil {
		fmt.Fprintf(os.Stderr,"MarshalIndent err: %v", outPutErr)
		os.Exit(1)

	}

	headerBytes := []byte(xml.Header)
	xmlOutPutData := append(headerBytes, xmlOutPut...)
	err = ioutil.WriteFile(confPath, xmlOutPutData, os.ModeAppend)
	if err != nil {
		fmt.Fprintf(os.Stderr,"WriteFile %s err: %v", confPath, err)
		os.Exit(1)
	}

	l4g.LoadConfiguration(confPath)
}

func Close(){
	time.Sleep(time.Second * 2)
	l4g.Close()
}
*/
