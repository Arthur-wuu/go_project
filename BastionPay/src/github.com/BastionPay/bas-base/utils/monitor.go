package utils

import (
	basel4g "BastionPay/bas-base/log/l4g"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/kataras/iris/core/errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ConfigReload interface {
	Reload() error
}

var GlobalMonitor Monitor

type Monitor struct {
	mConfigReload ConfigReload
}

func (this *Monitor) Start(addr string) error {
	l4g.Info("monitor addr[%s]", addr)
	run := func() error {
		serveMux := http.NewServeMux()
		//日志相关
		serveMux.HandleFunc("/log/loglevel/set", this.handleLogLevelSet) //日志纵向分类
		serveMux.HandleFunc("/log/loglevel/reset", this.handleLogLevelReset)
		serveMux.HandleFunc("/log/loglevel/watch", this.handleLogLevelWatch) //日志横向分类
		//配置相关
		serveMux.HandleFunc("/config/reload", this.handleConfigReload)
		//其它

		err := http.ListenAndServe(addr, serveMux)
		if err != nil {
			l4g.Error("monitor ListenAndServe[%s] err[%s]", addr, err.Error())
			return err
		}
		return nil
	}
	go run()
	l4g.Info("monitor addr[%s] start ok", addr)
	return nil
}

func (this *Monitor) Close() error {
	return nil
}

func (this *Monitor) RegConfigRerload(f ConfigReload) {
	this.mConfigReload = f
}

/****************************内部***************************/
//这个monitor可不能影响主程序的功能，必须捕获所有panic
func (this *Monitor) monRecover() {
	if err := recover(); err != nil {
		l4g.Error("panic err[%v]", err)
	}
}

func (this *Monitor) handleLogLevelSet(w http.ResponseWriter, req *http.Request) {
	defer this.monRecover()
	lv := req.FormValue("level")
	timeoutStr := req.FormValue("timeout")
	l4g.Info("handleLogLevelSet level[%s] timeout[%s] start", lv, timeoutStr)

	seconds, err := this.timeStrToSeconds(timeoutStr)
	if err != nil {
		l4g.Error("timeStrToSeconds lv[%s] timeout[%s] err[%s]", lv, timeoutStr, err.Error())
		io.WriteString(w, "timeout format err")
		return
	}

	if err = basel4g.GlobalLogMgr.SetLevel(lv, seconds); err != nil {
		io.WriteString(w, fmt.Sprintf("SetLevel err[%s]", err.Error()))
		l4g.Error("SetLevel lv[%s] timeout[%d] err[%s]", lv, seconds, err.Error())
		return
	}
	io.WriteString(w, fmt.Sprintf("ok and timeout is %d seconds", seconds))
	l4g.Info("handleLogLevelSet level[%s] timeout[%d] ok", lv, seconds)
	this.logTest()
}

func (this *Monitor) handleLogLevelReset(w http.ResponseWriter, req *http.Request) {
	defer this.monRecover()
	l4g.Info("handleLogLevelReset start")
	if err := basel4g.GlobalLogMgr.ResetLevel(); err != nil {
		l4g.Error("ResetLevel err[%s]", err.Error())
		io.WriteString(w, fmt.Sprintf("ResetLevel err[%s]", err.Error()))
		return
	}
	io.WriteString(w, "ok")
	l4g.Info("handleLogLevelReset ok")
	this.logTest()
}

func (this *Monitor) handleLogLevelWatch(w http.ResponseWriter, req *http.Request) {
	defer this.monRecover()
	l4g.Info("handleLogLevelWatch start")
	//elemInt := req.FormValue("int")
	//elemStr := req.FormValue("str")
	io.WriteString(w, "剧本暂无，敬请期待...")
	l4g.Info("handleLogLevelWatch ok")
}

func (this *Monitor) handleConfigReload(w http.ResponseWriter, req *http.Request) {
	defer this.monRecover()
	l4g.Info("handleConfigReload start")
	if this.mConfigReload == nil {
		io.WriteString(w, "not set ConfigReload interface")
		return
	}
	if err := this.mConfigReload.Reload(); err != nil {
		io.WriteString(w, fmt.Sprintf("ConfigReload err[%s]", err.Error()))
		return
	}
	io.WriteString(w, "ok")
	l4g.Info("handleConfigReload ok")
}

func (this *Monitor) timeStrToSeconds(str string) (int64, error) {
	if len(str) == 0 {
		return 24 * 60 * 60, nil
	}
	lowStr := strings.ToLower(str)
	index := strings.LastIndexAny(lowStr, "0123456789")
	intStr := lowStr[0 : index+1]
	seconds, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, err
	}
	charStr := lowStr[index+1:]
	if len(charStr) == 0 {
		return seconds, nil
	}
	switch charStr[0] {
	case 's':
		break
	case 'h':
		seconds *= 60 * 60
		break
	case 'm':
		seconds *= 60
		break
	case 'd':
		seconds *= 60 * 60 * 24
		break
	default:
		return 0, errors.New("Unknow format")
	}
	return seconds, nil
}

func (this *Monitor) logTest() {
	l4g.Finest("just test")
	l4g.Fine("just test")
	l4g.Debug("just test")
	l4g.Trace("just test")
	l4g.Info("just test")
	l4g.Warn("just test")
	l4g.Error("just test")
	l4g.Critical("just test")
}
