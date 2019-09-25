package service2

import (
	"BastionPay/bas-api/api"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	"BastionPay/bas-base/config"
	"BastionPay/bas-base/data"
	"BastionPay/bas-base/nethelper"
	"context"
	"encoding/json"
	l4g "github.com/alecthomas/log4go"
	"github.com/cenkalti/rpc2"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ServiceGateway struct {
	// rpc2
	*rpc2.Server

	// config
	cfgGateway config.ConfigGateway

	// srv nodes
	rwmu                       sync.RWMutex
	srvNodeNameMapSrvNodeGroup map[string]*SrvNodeGroup       // name+version mapto srvnodegroup
	clientMapSrvNodeGroup      map[*rpc2.Client]*SrvNodeGroup // name+version mapto srvnodegroup
	ApiInfo                    map[string]*data.ApiInfo

	// wait group
	wg sync.WaitGroup

	// center's apis
	registerData data.SrvRegisterData
	apiHandler   map[string]*NodeApi
}

// new a center
func NewServiceGateway(confPath string) (*ServiceGateway, error) {
	serviceGateway := &ServiceGateway{}

	serviceGateway.cfgGateway.Load(confPath)

	serviceGateway.srvNodeNameMapSrvNodeGroup = make(map[string]*SrvNodeGroup)
	serviceGateway.clientMapSrvNodeGroup = make(map[*rpc2.Client]*SrvNodeGroup)

	serviceGateway.registerData.Srv = serviceGateway.cfgGateway.GatewayName
	serviceGateway.registerData.Version = serviceGateway.cfgGateway.GatewayVersion

	serviceGateway.apiHandler = make(map[string]*NodeApi)

	func() {
		// api listsrv
		apiInfo := data.ApiInfo{Name: "listsrv", Level: data.APILevel_admin}

		serviceGateway.apiHandler[apiInfo.Name] = &NodeApi{ApiHandler: serviceGateway.listSrv, ApiInfo: apiInfo}
		serviceGateway.registerData.Functions = append(serviceGateway.registerData.Functions, apiInfo)
	}()

	// register
	var res string
	serviceGateway.register(nil, &serviceGateway.registerData, &res)

	// rpc2
	serviceGateway.Server = rpc2.NewServer()

	return serviceGateway, nil
}

// start the service center
func StartCenter(ctx context.Context, mi *ServiceGateway) {
	mi.startHttpServer(ctx)

	mi.startTcpServer(ctx)
}

// Stop the service center
func StopCenter(mi *ServiceGateway) {
	mi.wg.Wait()
}

func (mi *ServiceGateway) register(client *rpc2.Client, reg *data.SrvRegisterData, res *string) error {
	err := func() error {
		mi.rwmu.Lock()
		defer mi.rwmu.Unlock()

		versionSrvName := strings.ToLower(reg.Srv + "." + reg.Version)
		srvNodeGroup := mi.srvNodeNameMapSrvNodeGroup[versionSrvName]
		if srvNodeGroup == nil {
			srvNodeGroup = &SrvNodeGroup{}
			mi.srvNodeNameMapSrvNodeGroup[versionSrvName] = srvNodeGroup
		}

		mi.clientMapSrvNodeGroup[client] = srvNodeGroup

		err := srvNodeGroup.RegisterNode(client, reg)
		if err == nil {
			if mi.ApiInfo == nil {
				mi.ApiInfo = make(map[string]*data.ApiInfo)
			}

			for _, v := range reg.Functions {
				mi.ApiInfo[strings.ToLower(versionSrvName+"."+v.Name)] = &data.ApiInfo{v.Name, v.Level}
			}
		}

		l4g.Info("%d-%d", len(mi.srvNodeNameMapSrvNodeGroup), len(mi.clientMapSrvNodeGroup))

		*res = "ok"
		return err
	}()

	return err
}

func (mi *ServiceGateway) unRegister(client *rpc2.Client, reg *data.SrvRegisterData, res *string) error {
	mi.disconnectClient(client)
	*res = "ok"
	return nil
}

func (mi *ServiceGateway) innerNotify(client *rpc2.Client, req *data.SrvRequest, res *data.SrvResponse) error {
	mi.wg.Add(1)
	defer mi.wg.Done()

	l4g.Debug("notify %s", req.String())

	mi.rwmu.RLock()
	defer mi.rwmu.RUnlock()

	for _, srvNodeGroup := range mi.srvNodeNameMapSrvNodeGroup {
		srvNodeGroup.Notify(client, req)
	}

	return nil
}

func (mi *ServiceGateway) innerCall(client *rpc2.Client, req *data.SrvRequest, res *data.SrvResponse) error {
	mi.wg.Add(1)
	defer mi.wg.Done()

	mi.srvCall(req, res)

	// make sure no data if err
	if res.Err != apibackend.NoErr {
		res.Value.Message = ""
		res.Value.Signature = ""
	}
	return nil
}

// start http server
func (mi *ServiceGateway) startHttpServer(ctx context.Context) {
	// http
	l4g.Debug("Start http server on %s", mi.cfgGateway.Port)

	http.Handle("/"+apibackend.HttpRouterApi+"/", http.HandlerFunc(mi.handleApi))
	http.Handle("/"+apibackend.HttpRouterUser+"/", http.HandlerFunc(mi.handleUser))
	http.Handle("/"+apibackend.HttpRouterAdmin+"/", http.HandlerFunc(mi.handleAdmin))
	http.Handle("/health", http.HandlerFunc(mi.handleHealth))

	// test mode
	if mi.cfgGateway.TestMode != 0 {
		http.Handle("/"+apibackend.HttpRouterApiTest+"/", http.HandlerFunc(mi.handleApiTest))
	}

	go func() {
		l4g.Info("Http server routine running... ")
		//err := http.ListenAndServe(":"+mi.cfgGateway.Port, nil)

		server := &http.Server{
			Addr:         ":" + mi.cfgGateway.Port,
			Handler:      nil,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}

		err := server.ListenAndServe()
		if err != nil {
			l4g.Crashf("", err)
		}
	}()
}

func (mi *ServiceGateway) disconnectClient(client *rpc2.Client) {
	mi.rwmu.Lock()
	defer mi.rwmu.Unlock()

	srvNodeGroup, ok := mi.clientMapSrvNodeGroup[client]
	if srvNodeGroup == nil || !ok {
		return
	}

	srvNodeGroup.UnRegisterNode(client)
	if srvNodeGroup.GetSrvNodes() == 0 {
		// remove srv node group
		reg, _ := srvNodeGroup.GetSrvInfo()
		versionSrvName := strings.ToLower(reg.Srv + "." + reg.Version)
		if mi.ApiInfo != nil {
			for _, v := range reg.Functions {
				delete(mi.ApiInfo, strings.ToLower(versionSrvName+"."+v.Name))
			}
		}

		delete(mi.srvNodeNameMapSrvNodeGroup, versionSrvName)
	}

	delete(mi.clientMapSrvNodeGroup, client)
}

// start tcp server
func (mi *ServiceGateway) startTcpServer(ctx context.Context) {
	mi.Server.OnConnect(func(client *rpc2.Client) {
		l4g.Info("rpc2 client connect...")
	})

	mi.Server.OnDisconnect(func(client *rpc2.Client) {
		l4g.Info("rpc2 client disconnect...")

		mi.disconnectClient(client)
	})

	mi.Server.Handle(data.MethodCenterRegister, mi.register)
	mi.Server.Handle(data.MethodCenterUnRegister, mi.unRegister)
	mi.Server.Handle(data.MethodCenterInnerNotify, mi.innerNotify)
	mi.Server.Handle(data.MethodCenterInnerCall, mi.innerCall)

	l4g.Debug("Start Tcp server on %s", mi.cfgGateway.GatewayPort)

	listener, err := nethelper.CreateTcpServer(":" + mi.cfgGateway.GatewayPort)
	if err != nil {
		l4g.Crashf("", err)
	}
	go func() {
		mi.wg.Add(1)
		defer mi.wg.Done()

		l4g.Info("Tcp server routine running... ")

		go mi.Server.Accept(listener)
		<-ctx.Done()

		l4g.Info("Tcp server routine stoped... ")
	}()
}

func (mi *ServiceGateway) srvCall(req *data.SrvRequest, res *data.SrvResponse) {
	api := mi.getApiInfo(req)
	if api == nil {
		res.Err = apibackend.ErrNotFindSrv
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// call function
	var rpcSrv data.SrvRequest
	rpcSrv.Method = req.Method
	rpcSrv.Argv = req.Argv
	rpcSrv.Context.ApiLever = api.Level
	var rpcSrvRes data.SrvResponse
	if mi.callFunction(&rpcSrv, &rpcSrvRes); rpcSrvRes.Err != apibackend.NoErr {
		*res = rpcSrvRes
		return
	}

	*res = rpcSrvRes
}

// user call by api
func (mi *ServiceGateway) apiCall(req *data.SrvRequest, res *data.SrvResponse, sourceIp string) {
	// can not call auth service
	if req.Method.Srv == "auth" {
		res.Err = apibackend.ErrIllegallyCall
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// find api
	api := mi.getApiInfo(req)
	if api == nil {
		res.Err = apibackend.ErrNotFindSrv
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// decode and verify data
	var rpcAuth data.SrvRequest
	rpcAuth = *req
	rpcAuth.Context.ApiLever = api.Level
	rpcAuth.Context.DataFrom = apibackend.DataFromApi
	if mi.cfgGateway.EnableRemoteIp {
		rpcAuth.Context.SetSourceIp(sourceIp)
	}
	var rpcAuthRes data.SrvResponse
	if mi.authData(&rpcAuth, &rpcAuthRes); rpcAuthRes.Err != apibackend.NoErr {
		*res = rpcAuthRes
		return
	}

	// call real srv
	var rpcSrv data.SrvRequest
	rpcSrv = *req
	rpcSrv.Context.ApiLever = api.Level
	rpcSrv.Context.DataFrom = apibackend.DataFromApi
	rpcSrv.Argv.Message = rpcAuthRes.Value.Message
	var rpcSrvRes data.SrvResponse
	if mi.callFunction(&rpcSrv, &rpcSrvRes); rpcSrvRes.Err != apibackend.NoErr {
		*res = rpcSrvRes
		return
	}

	// encode and sign data
	var reqEncrypted data.SrvRequest
	reqEncrypted = *req
	reqEncrypted.Context.ApiLever = api.Level
	reqEncrypted.Context.DataFrom = apibackend.DataFromApi
	reqEncrypted.Argv.Message = rpcSrvRes.Value.Message
	reqEncrypted.Argv.UserPubKey = rpcSrvRes.Value.UserPubKey
	var reqEncryptedRes data.SrvResponse
	if mi.encryptData(&reqEncrypted, &reqEncryptedRes); reqEncryptedRes.Err != apibackend.NoErr {
		*res = reqEncryptedRes
		return
	}

	*res = reqEncryptedRes
}

// user call by user
func (mi *ServiceGateway) userCall(req *data.SrvRequest, res *data.SrvResponse) {
	// can not call auth service
	if req.Method.Srv == "auth" {
		res.Err = apibackend.ErrIllegallyCall
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// find api
	api := mi.getApiInfo(req)
	if api == nil {
		res.Err = apibackend.ErrNotFindSrv
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// decode and verify data
	var rpcAuth data.SrvRequest
	rpcAuth = *req
	rpcAuth.Context.ApiLever = api.Level
	rpcAuth.Context.DataFrom = apibackend.DataFromUser
	var rpcAuthRes data.SrvResponse
	if mi.authData(&rpcAuth, &rpcAuthRes); rpcAuthRes.Err != apibackend.NoErr {
		*res = rpcAuthRes
		return
	}

	// 解析来自用户的数据后台
	userParams := apibackend.UserMessage{}
	err := json.Unmarshal([]byte(rpcAuthRes.Value.Message), &userParams)
	if err != nil {
		res.Err = apibackend.ErrDataCorrupted
		l4g.Error("parse user params %s %d %s", req.String(), res.Err, err.Error())
		return
	}

	// call real srv
	var rpcSrv data.SrvRequest
	rpcSrv = *req
	rpcSrv.Context.ApiLever = api.Level
	rpcSrv.Context.DataFrom = apibackend.DataFromUser
	rpcSrv.Argv.SubUserKey = userParams.SubUserKey
	rpcSrv.Argv.Message = userParams.Message
	var rpcSrvRes data.SrvResponse
	if mi.callFunction(&rpcSrv, &rpcSrvRes); rpcSrvRes.Err != apibackend.NoErr {
		*res = rpcSrvRes
		return
	}

	// encode and sign data
	var reqEncrypted data.SrvRequest
	reqEncrypted = *req
	reqEncrypted.Context.ApiLever = api.Level
	reqEncrypted.Context.DataFrom = apibackend.DataFromUser
	reqEncrypted.Argv.Message = rpcSrvRes.Value.Message
	var reqEncryptedRes data.SrvResponse
	if mi.encryptData(&reqEncrypted, &reqEncryptedRes); reqEncryptedRes.Err != apibackend.NoErr {
		*res = reqEncryptedRes
		return
	}

	*res = reqEncryptedRes
}

// user call by admin
func (mi *ServiceGateway) adminCall(req *data.SrvRequest, res *data.SrvResponse) {
	// can not call auth service
	if req.Method.Srv == "auth" {
		res.Err = apibackend.ErrIllegallyCall
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// find api
	api := mi.getApiInfo(req)
	if api == nil {
		res.Err = apibackend.ErrNotFindSrv
		l4g.Error("%s %d", req.String(), res.Err)
		return
	}

	// decode and verify data
	var rpcAuth data.SrvRequest
	rpcAuth = *req
	rpcAuth.Context.ApiLever = api.Level
	rpcAuth.Context.DataFrom = apibackend.DataFromAdmin
	var rpcAuthRes data.SrvResponse
	if mi.authData(&rpcAuth, &rpcAuthRes); rpcAuthRes.Err != apibackend.NoErr {
		*res = rpcAuthRes
		return
	}

	// 解析来自管理员的数据后台
	adminParams := apibackend.AdminMessage{}
	err := json.Unmarshal([]byte(rpcAuthRes.Value.Message), &adminParams)
	if err != nil {
		res.Err = apibackend.ErrDataCorrupted
		l4g.Error("parse admin params %s %d %s", req.String(), res.Err, err.Error())
		return
	}

	// call real srv
	var rpcSrv data.SrvRequest
	rpcSrv = *req
	rpcSrv.Context.ApiLever = api.Level
	rpcSrv.Context.DataFrom = apibackend.DataFromAdmin
	rpcSrv.Argv.SubUserKey = adminParams.SubUserKey
	rpcSrv.Argv.Message = adminParams.Message
	var rpcSrvRes data.SrvResponse
	if mi.callFunction(&rpcSrv, &rpcSrvRes); rpcSrvRes.Err != apibackend.NoErr {
		*res = rpcSrvRes
		return
	}

	// encode and sign data
	var reqEncrypted data.SrvRequest
	reqEncrypted = *req
	reqEncrypted.Context.ApiLever = api.Level
	reqEncrypted.Context.DataFrom = apibackend.DataFromAdmin
	reqEncrypted.Argv.Message = rpcSrvRes.Value.Message
	var reqEncryptedRes data.SrvResponse
	if mi.encryptData(&reqEncrypted, &reqEncryptedRes); reqEncryptedRes.Err != apibackend.NoErr {
		*res = reqEncryptedRes
		return
	}

	*res = reqEncryptedRes
}

// auth data
func (mi *ServiceGateway) authData(req *data.SrvRequest, res *data.SrvResponse) {
	reqAuth := *req

	reqAuth.Method.Srv = "auth"
	reqAuth.Method.Function = "AuthData"
	reqAuthRes := data.SrvResponse{}

	mi.callFunction(&reqAuth, &reqAuthRes)

	*res = reqAuthRes
}

// package data
func (mi *ServiceGateway) encryptData(req *data.SrvRequest, res *data.SrvResponse) {
	reqEnc := *req
	reqEnc.Method.Srv = "auth"
	reqEnc.Method.Function = "EncryptData"

	reqEncRes := data.SrvResponse{}

	mi.callFunction(&reqEnc, &reqEncRes)

	*res = reqEncRes
}

//  call a srv node
func (mi *ServiceGateway) callFunction(req *data.SrvRequest, res *data.SrvResponse) {
	centerVersionSrvName := strings.ToLower(mi.registerData.Srv + "." + mi.registerData.Version)
	versionSrvName := strings.ToLower(req.Method.Srv + "." + req.Method.Version)
	l4g.Debug("call %s", req.String())

	mi.rwmu.RLock()
	defer mi.rwmu.RUnlock()

	if centerVersionSrvName == versionSrvName {
		h := mi.apiHandler[strings.ToLower(req.Method.Function)]
		if h != nil {
			h.ApiHandler(req, res)
		} else {
			res.Err = apibackend.ErrNotFindFunction
		}
		if res.Err != apibackend.NoErr {
			l4g.Error("call failed: %d", res.Err)
		}
		return
	} else {
		srvNodeGroup := mi.srvNodeNameMapSrvNodeGroup[versionSrvName]
		if srvNodeGroup == nil {
			res.Err = apibackend.ErrNotFindSrv
			l4g.Error("%s %d", req.String(), res.Err)
			return
		}

		srvNodeGroup.Call(req, res)
	}
}

func (mi *ServiceGateway) getApiInfo(req *data.SrvRequest) *data.ApiInfo {
	mi.rwmu.RLock()
	defer mi.rwmu.RUnlock()

	name := strings.ToLower(req.Method.Srv + "." + req.Method.Version + "." + req.Method.Function)
	return mi.ApiInfo[name]
}

func (mi *ServiceGateway) handleApi(w http.ResponseWriter, req *http.Request) {
	l4g.Debug("Http server Accept a api client: %s, %s", req.RemoteAddr, mi.getRemoteAddr(req))
	defer req.Body.Close()

	//w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	mi.wg.Add(1)
	defer mi.wg.Done()

	userResponse := api.UserResponseData{}
	func() {
		//fmt.Println("path=", req.URL.Path)
		reqData := data.SrvRequest{}

		reqData.Method.FromPath(req.URL.Path)

		// get argv
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		//body := string(b)
		//fmt.Println("body=", string(b))

		// make data
		userData := api.UserData{}
		err = json.Unmarshal(b, &userData)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		reqData.Argv.FromApiData(&userData)

		resData := data.SrvResponse{}
		mi.apiCall(&reqData, &resData, mi.getRemoteAddr(req))

		resData.ToApiResponse(&userResponse)
		userResponse.Value.UserKey = userData.UserKey
	}()

	if userResponse.Err != apibackend.NoErr && userResponse.ErrMsg == "" {
		userResponse.ErrMsg = apibackend.GetErrMsg(userResponse.Err)
	}

	if userResponse.Err != apibackend.NoErr {
		l4g.Error("handleAPi request err: %d-%s", userResponse.Err, userResponse.ErrMsg)
	}

	// write back http
	connectionType := req.Header.Get("Connection")
	w.Header().Set("Connection", connectionType)
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(userResponse)
	w.Write(b)
	return
}

func (mi *ServiceGateway) handleUser(w http.ResponseWriter, req *http.Request) {
	l4g.Debug("Http server Accept a user client: %s", req.RemoteAddr)
	defer req.Body.Close()

	//w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	mi.wg.Add(1)
	defer mi.wg.Done()

	userResponse := api.UserResponseData{}
	func() {
		//fmt.Println("path=", req.URL.Path)
		reqData := data.SrvRequest{}

		reqData.Method.FromPath(req.URL.Path)

		// get argv
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		//body := string(b)
		//fmt.Println("body=", string(b))

		// make data
		userData := api.UserData{}
		err = json.Unmarshal(b, &userData)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		reqData.Argv.FromApiData(&userData)

		resData := data.SrvResponse{}
		mi.userCall(&reqData, &resData)

		resData.ToApiResponse(&userResponse)
		userResponse.Value.UserKey = userData.UserKey
	}()

	if userResponse.Err != apibackend.NoErr && userResponse.ErrMsg == "" {
		userResponse.ErrMsg = apibackend.GetErrMsg(userResponse.Err)
	}

	if userResponse.Err != apibackend.NoErr {
		l4g.Error("handleUser request err: %d-%s", userResponse.Err, userResponse.ErrMsg)
	}

	// write back http
	connectionType := req.Header.Get("Connection")
	w.Header().Set("Connection", connectionType)
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(userResponse)
	w.Write(b)
	return
}

func (mi *ServiceGateway) handleAdmin(w http.ResponseWriter, req *http.Request) {
	l4g.Debug("Http server Accept a admin client: %s", req.RemoteAddr)
	defer req.Body.Close()

	//w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	mi.wg.Add(1)
	defer mi.wg.Done()

	userResponse := api.UserResponseData{}
	func() {
		//fmt.Println("path=", req.URL.Path)
		reqData := data.SrvRequest{}

		reqData.Method.FromPath(req.URL.Path)

		// get argv
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		//body := string(b)
		//fmt.Println("body=", string(b))

		// make data
		userData := api.UserData{}
		err = json.Unmarshal(b, &userData)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		reqData.Argv.FromApiData(&userData)

		resData := data.SrvResponse{}
		mi.adminCall(&reqData, &resData)

		resData.ToApiResponse(&userResponse)
		userResponse.Value.UserKey = userData.UserKey
	}()

	if userResponse.Err != apibackend.NoErr && userResponse.ErrMsg == "" {
		userResponse.ErrMsg = apibackend.GetErrMsg(userResponse.Err)
	}

	if userResponse.Err != apibackend.NoErr {
		l4g.Error("handleAdmin request err: %d-%s", userResponse.Err, userResponse.ErrMsg)
	}

	// write back http
	connectionType := req.Header.Get("Connection")
	w.Header().Set("Connection", connectionType)
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(userResponse)
	w.Write(b)
	return
}

func (mi *ServiceGateway) handleApiTest(w http.ResponseWriter, req *http.Request) {
	l4g.Debug("Http server test Accept a api test client: %s", req.RemoteAddr)
	defer req.Body.Close()

	//w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	mi.wg.Add(1)
	defer mi.wg.Done()

	userResponse := api.UserResponseData{}
	func() {
		//fmt.Println("path=", req.URL.Path)
		reqData := data.SrvRequest{}

		reqData.Method.FromPath(req.URL.Path)

		// get argv
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		userData := api.UserData{}
		err = json.Unmarshal(b, &userData)
		if err != nil {
			l4g.Error("http handler: %s", err.Error())
			userResponse.Err = apibackend.ErrDataCorrupted
			return
		}

		reqData.Argv.FromApiData(&userData)

		resData := data.SrvResponse{}
		mi.srvCall(&reqData, &resData)

		resData.ToApiResponse(&userResponse)
		userResponse.Value.UserKey = userData.UserKey
	}()

	if userResponse.Err != apibackend.NoErr && userResponse.ErrMsg == "" {
		userResponse.ErrMsg = apibackend.GetErrMsg(userResponse.Err)
	}

	if userResponse.Err != apibackend.NoErr {
		l4g.Error("handleApiTest request err: %d-%s", userResponse.Err, userResponse.ErrMsg)
	}

	// write back http
	connectionType := req.Header.Get("Connection")
	w.Header().Set("Connection", connectionType)
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(userResponse)
	w.Write(b)
	return
}

func (mi *ServiceGateway) handleHealth(w http.ResponseWriter, req *http.Request) {
	l4g.Debug("Http server Accept health client: %s", req.RemoteAddr)
	defer req.Body.Close()

	nodes := mi.getSrvInfo()

	b, err := json.Marshal(nodes)
	if err != nil {
		return
	}

	w.Header().Set("Connection", "close")
	w.Write(b)
}

// RPC -- listsrv
func (mi *ServiceGateway) listSrv(req *data.SrvRequest, res *data.SrvResponse) {
	nodes := mi.getSrvInfo()

	b, err := json.Marshal(nodes)
	if err != nil {
		res.Err = apibackend.ErrDataCorrupted
		res.Value.Message = ""
		res.Value.Signature = ""
		return
	}
	res.Value.Message = string(b)

	// make sure no data if err
	if res.Err != apibackend.NoErr {
		res.Value.Message = ""
		res.Value.Signature = ""
	}
}

func (mi *ServiceGateway) getSrvInfo() *backend.ServiceInfoList {
	mi.rwmu.RLock()
	defer mi.rwmu.RUnlock()

	serviceInfoList := &backend.ServiceInfoList{}
	for _, v := range mi.srvNodeNameMapSrvNodeGroup {
		n, c := v.GetSrvInfo()

		node := backend.ServiceInfo{
			Version: n.Version,
			Srv:     n.Srv,
			Count:   c,
		}
		serviceInfoList.Data = append(serviceInfoList.Data, node)
	}

	return serviceInfoList
}

func (mi *ServiceGateway) getRemoteAddr(req *http.Request) string {
	if req == nil {
		return ""
	}
	//	l4g.Debug("http_x_forwarded_for=", req.Header.Get("http_x_forwarded_for"))
	remote := req.Header.Get("X-Forwarded-For")
	if len(remote) == 0 {
		remote = req.RemoteAddr
	}
	l4g.Debug("X-Forwarded-For=%s remote_addr=%s", req.Header.Get("X-Forwarded-For"), req.RemoteAddr)
	return remote
}
