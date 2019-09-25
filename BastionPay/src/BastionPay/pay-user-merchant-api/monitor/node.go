package monitor

import (
	"BastionPay/bas-monitor/common"
	"BastionPay/bas-monitor/rpc"
	"BastionPay/bas-monitor/defaultrpc"
	"BastionPay/bas-monitor/logger"
	"fmt"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
)

type MyLogger struct {
}

func (l *MyLogger) Debug(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Debug("bas_monitor", zap.String("text", text))
}
func (l *MyLogger) Info(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Info("bas_monitor", zap.String("text", text))
}

func (l *MyLogger) Trace(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Info("bas_monitor", zap.String("text", text))
}
func (l *MyLogger) Warns(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Warn("bas_monitor", zap.String("text", text))
}
func (l *MyLogger) Error(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Error("bas_monitor", zap.String("text", text))
}
func (l *MyLogger) Fatal(arg0 string, args ...interface{}) {
	text := fmt.Sprintf(arg0, args...)
	ZapLog().Fatal("bas_monitor", zap.String("text", text))
}

func NewNode(cfg *common.ConfigNode, meta string) (*rpc.Node, error) {
	logger.InitLogger(&MyLogger{})

	if cfg == nil {
		return nil, fmt.Errorf("cfg is nil")
	}

	nodeInst, err := rpc.NewNode(*cfg, meta, connectCenter)
	if err != nil {
		return nil, err
	}

	defaultrpc.SetDefaultNodeInst(nodeInst)

	return nodeInst, nil
}

func BeginNodeMonitor()  {
	initCaller()
	initNotifier()
}