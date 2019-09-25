package nethelper

import (
	"net"
	l4g "github.com/alecthomas/log4go"
)

// Create a tcp server
// @parameter: port string, like: ":8080"
// @return: net.TCPListener, error
func CreateTcpServer(port string) (*net.TCPListener, error){
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		l4g.Error("%s", err.Error())
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		l4g.Error("%s", err.Error())
		return nil, err
	}

	return listener, nil
}
