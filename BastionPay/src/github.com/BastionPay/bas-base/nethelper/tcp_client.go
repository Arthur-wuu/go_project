package nethelper

import (
	"net/rpc"
	l4g "github.com/alecthomas/log4go"
)

// Call a JRPC to Tcp server
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params interface{}
// @parameter: res interface{}
// @return: error
func CallJRPCToTcpServer(addr string, method string, params interface{}, res interface{}) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		l4g.Error("dial %s", err.Error())
		return err
	}
	defer client.Close()

	return CallJRPCToTcpServerOnClient(client, method, params, res)
}

// Call a JRPC to Tcp server on a client
// @parameter: client rpc.Client
// @parameter: addr string, like "127.0.0.1:8080"
// @parameter: method string
// @parameter: params interface{}
// @parameter: res interface{}
// @return: error
func CallJRPCToTcpServerOnClient(client *rpc.Client, method string, params interface{}, res interface{}) error {
	err := client.Call(method, params, res)
	if err != nil {
		l4g.Error("call %s", err.Error())
		return err
	}

	return nil
}