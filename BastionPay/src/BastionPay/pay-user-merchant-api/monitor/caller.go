package monitor

import (
	"BastionPay/bas-monitor/defaultrpc"
)

func initCaller()  {
	nodeInst := defaultrpc.DefaultNodeInst()
	if nodeInst == nil {
		return
	}

	//nodeInst.GetApiGroup().RegisterCaller("listsrv", func(req *common.Request, res *common.Response) {
	//	res.Data.SetResult(gatewayInstance.GetSrvInfo())
	//})
}
