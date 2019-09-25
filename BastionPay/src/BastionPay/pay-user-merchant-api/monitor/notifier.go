package monitor

import (
	"BastionPay/bas-monitor/defaultrpc"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
)

func connectCenter(status int)  {
	ZapLog().Info("bas_monitor", zap.Int("connectCenter", status))
}

func initNotifier()  {
	nodeInst := defaultrpc.DefaultNodeInst()
	if nodeInst == nil {
		return
	}
}