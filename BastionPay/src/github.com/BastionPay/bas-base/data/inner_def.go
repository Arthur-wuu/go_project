package data

import (
	"BastionPay/bas-api/apibackend"
	"fmt"
)

// /////////////////////////////////////////////////////
// internal api gateway and service RPC data define
// /////////////////////////////////////////////////////

const (
	MethodCenterRegister    = "ServiceCenter.Register"    // register to center
	MethodCenterUnRegister  = "ServiceCenter.UnRegister"  // unregister to center
	MethodCenterInnerNotify = "ServiceCenter.InnerNotify" // notify data to center to nodes
	MethodCenterInnerCall   = "ServiceCenter.InnerCall"   // call a api to center

	MethodNodeCall   = "ServiceNode.Call"   // center call a srv node function
	MethodNodeNotify = "ServiceNode.Notify" // center notify to a srv node function
)

const (
	// normal client
	UserClass_Client = 0

	// hot
	UserClass_Hot = 1

	// admin
	UserClass_Admin = 2

	// client
	APILevel_client = 0

	// common administrator
	APILevel_admin = 100

	// genesis administrator
	APILevel_genesis = 200
)

// API info
type ApiInfo struct {
	Name  string `json:"name"`  // api name
	Level int    `json:"level"` // api level, refer APILevel_*
}

// srv register data
type SrvRegisterData struct {
	Version   string    `json:"version"`   // srv version
	Srv       string    `json:"srv"`       // srv name
	Functions []ApiInfo `json:"functions"` // srv functions
}

// srv context
type SrvContext struct {
	ApiLever int     `json:"apilevel"`         // api info level
	DataFrom int     `json:"datafrom"`         // data from router
	SourceIp *string `json:"src_ip,omitempty"` //source ip
	// future...
}

func (this *SrvContext) SetSourceIp(ip string) {
	if this.SourceIp == nil {
		this.SourceIp = new(string)
	}
	*this.SourceIp = ip
}

func (this *SrvContext) GetSourceIp() string {
	if this.SourceIp == nil {
		return ""
	}
	return *this.SourceIp
}

// srv data
type SrvData struct {
	// user unique key
	UserKey string `json:"user_key"`
	// sub user key
	SubUserKey string `json:"sub_user_key"`
	// user request message
	Message string `json:"message"`
	// timestamp = Unix timestamp string
	TimeStamp string `json:"time_stamp"`
	// signature = origin data -> sha512 -> rsa sign -> base64
	Signature string `json:"signature"`
	//用户公钥信息，在用户更新公钥时使用
	UserPubKey string `json:"user_pubkey"`
}

// input/output method
type SrvMethod struct {
	Version  string `json:"version"`  // srv version
	Srv      string `json:"srv"`      // srv name
	Function string `json:"function"` // srv function
}

// srv request
type SrvRequest struct {
	Context SrvContext `json:"context"` // api info
	Method  SrvMethod  `json:"method"`  // request method
	Argv    SrvData    `json:"argv"`    // request argument
}

// srv response/push
type SrvResponse struct {
	Err    int     `json:"err"`    // error code
	ErrMsg string  `json:"errmsg"` // error message
	Value  SrvData `json:"value"`  // response data
}

//////////////////////////////////////////////////////////////////////
func (sr *SrvRequest) GetAccessUserKey() string {
	if sr.Context.DataFrom == apibackend.DataFromUser {
		return sr.Argv.SubUserKey
	} else if sr.Context.DataFrom == apibackend.DataFromAdmin {
		return sr.Argv.SubUserKey
	} else {
		return sr.Argv.UserKey
	}
}

func (urd SrvRequest) String() string {
	return fmt.Sprintf("%s %s-%s-%s", urd.Argv.UserKey, urd.Method.Srv, urd.Method.Version, urd.Method.Function)
}
