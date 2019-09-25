package service2

import (
	"sync"
	"BastionPay/bas-base/data"
	"BastionPay/bas-base/config"
	"context"
	"strings"
	"log"
	"time"
	"errors"
	"net"
	l4g "github.com/alecthomas/log4go"
	"github.com/cenkalti/rpc2"
	"BastionPay/bas-api/apibackend"
)

// node api interface
type NodeApiHandler func(req *data.SrvRequest, res *data.SrvResponse)
type NodeApi struct{
	ApiInfo 	data.ApiInfo
	ApiHandler 	NodeApiHandler
}

func RegisterApi(nap *map[string]NodeApi, name string, level int, handler NodeApiHandler) error {
	if _, ok := (*nap)[name]; ok {
		return errors.New("function exist")
	}

	apiInfo := data.ApiInfo{Name:name, Level:level}
	(*nap)[name] = NodeApi{ApiHandler:handler, ApiInfo:apiInfo}

	return nil
}
type NodeApiGroup interface {
	GetApiGroup()(map[string]NodeApi)
 	HandleNotify(req *data.SrvRequest)
}

// service node
type ServiceNode struct{
	// register data
	registerData data.SrvRegisterData

	// callback
	apiHandler map[string]*NodeApi

	// center addr
	serviceCenterAddr string

	// wait group
	wg sync.WaitGroup

	// connection to center
	rwmu sync.RWMutex
	client *rpc2.Client

	// handler
	nodeApiGroup NodeApiGroup
}

// New a service node
func NewServiceNode(confPath string) (*ServiceNode, error){
	cfgNode := config.ConfigNode{}
	cfgNode.Load(confPath)

	serviceNode := &ServiceNode{}

	serviceNode.apiHandler = make(map[string]*NodeApi)

	// node info
	serviceNode.registerData.Srv = cfgNode.SrvName
	serviceNode.registerData.Version = cfgNode.SrvVersion

	// center info
	serviceNode.serviceCenterAddr = cfgNode.CenterAddr

	return serviceNode, nil
}

// register api group
func RegisterNodeApi(ni *ServiceNode, nodeApiGroup NodeApiGroup) {
	ni.nodeApiGroup = nodeApiGroup
	if(nodeApiGroup == nil){
		return
	}

	nam := ni.nodeApiGroup.GetApiGroup()

	for k, v := range nam{
		if ni.apiHandler[k] != nil {
			log.Fatal("#Error api repeat:", k)
		}
		ni.apiHandler[k] = &NodeApi{ApiInfo:v.ApiInfo, ApiHandler:v.ApiHandler}
		ni.registerData.Functions = append(ni.registerData.Functions, v.ApiInfo)
	}
}

// Start the service node
func StartNode(ctx context.Context, ni *ServiceNode) {
	ni.startToCenter(ctx)
}

// Stop the service node
func StopNode(ni *ServiceNode)  {
	ni.wg.Wait()
}

// RPC -- call
func (ni *ServiceNode) call(client *rpc2.Client, req *data.SrvRequest, res *data.SrvResponse) error {
	l4g.Debug("got call:", req)

	h := ni.apiHandler[strings.ToLower(req.Method.Function)]
	if h != nil {
		h.ApiHandler(req, res)
	}else{
		res.Err = apibackend.ErrNotFindFunction
	}
	if res.Err != apibackend.NoErr {
		l4g.Error("call failed: %d", res.Err)
	}
	return nil
}

// RPC -- call
func (ni *ServiceNode) notify(client *rpc2.Client, req *data.SrvRequest, res *data.SrvResponse) error {
	l4g.Debug("got notify:", req.Method)

	if ni.nodeApiGroup != nil {
		ni.nodeApiGroup.HandleNotify(req)
	}

	return nil
}

// inner call a request to router
func (ni *ServiceNode) InnerCall(req *data.SrvRequest, res *data.SrvResponse) error {
	ni.rwmu.RLock()
	defer ni.rwmu.RUnlock()

	var err error
	if ni.client != nil {
		err = ni.client.Call(data.MethodCenterInnerCall, req, res)
	}else{
		err = errors.New("client is nil")
	}
	return err
}

// inner notify a request to router
func (ni *ServiceNode) InnerNotify(req *data.SrvRequest, res *data.SrvResponse) error {
	ni.rwmu.RLock()
	defer ni.rwmu.RUnlock()

	var err error
	if ni.client != nil {
		err = ni.client.Notify(data.MethodCenterInnerNotify, req)
	}else{
		err = errors.New("client is nil")
	}
	return err
}

func (ni *ServiceNode)connectToCenter() (*rpc2.Client, error){
	conn, err := net.Dial("tcp", ni.serviceCenterAddr)
	if err != nil {
		return nil, err
	}

	clt := rpc2.NewClient(conn)
	return clt, nil
}

func (ni *ServiceNode)registToCenter() error{
	var err error
	var res string
	if ni.client != nil {
		err = ni.client.Call(data.MethodCenterRegister, ni.registerData, &res)
	}else{
		err = errors.New("client is nil")
	}
	return err
}

func (ni *ServiceNode)unRegistToCenter() error{
	var err error
	var res string
	if ni.client != nil {
		ni.client.Call(data.MethodCenterUnRegister, ni.registerData, &res)
	}else{
		err = errors.New("client is nil")
	}

	return err
}

func (ni *ServiceNode)startToCenter(ctx context.Context) {
	go func() {
		ni.wg.Add(1)
		defer ni.wg.Done()

		go func() {
			for {
				// connect and regist
				func(){
					ni.rwmu.Lock()
					defer ni.rwmu.Unlock()

					var err error
					if ni.client == nil {
						l4g.Info("client try to connect...")
						ni.client, err = ni.connectToCenter()
						if ni.client != nil && err == nil {
							l4g.Info("client connect to center...")
							ni.client.Handle(data.MethodNodeCall, ni.call)
							ni.client.Handle(data.MethodNodeNotify, ni.notify)

							go ni.client.Run()

							ni.registToCenter()
						}
					}

					if err != nil {
						if(ni.client != nil){
							ni.client.Close()
							ni.client = nil
						}
						l4g.Error("connect failed, %s", err.Error())
					}
				}()

				// listen
				func() {
					ni.rwmu.RLock()
					defer ni.rwmu.RUnlock()

					if ni.client == nil {
						return
					}

					l4g.Info("client run...")
					<- ni.client.DisconnectNotify()
					l4g.Error("client disconnect...")
				}()

				func() {
					ni.rwmu.Lock()
					defer ni.rwmu.Unlock()
					if(ni.client != nil){
						ni.client.Close()
						ni.client = nil
					}
					l4g.Info("reset client...")
				}()

				l4g.Info("wait 5 second to connect...")
				time.Sleep(time.Second*5)
			}
		}()

		<-ctx.Done()
		ni.unRegistToCenter()
		l4g.Info("UnRegist to center ok %s.%s", ni.registerData.Version, ni.registerData.Srv)
	}()
}
