package handler

import (
	"BastionPay/bas-account-srv/db"
	"BastionPay/bas-base/data"
	"io/ioutil"
	//service "github.com/BastionPay/bas-service/base/service"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	"BastionPay/bas-api/utils"
	"BastionPay/bas-base/config"
	service "BastionPay/bas-base/service2"
	sdkmail "BastionPay/bas-tools/sdk.notify.mail"
	"encoding/json"
	l4g "github.com/alecthomas/log4go"
	"github.com/satori/go.uuid"
)

///////////////////////////////////////////////////////////////////////
// 账号管理
type Account struct {
	node *service.ServiceNode

	privateKey                  []byte
	serverPublicKey             []byte
	auditeTemplateName          string
	userFrozen1ForAdminTempName string
	userFrozen1ForUserTempName  string
}

// 默认实例
var defaultAccount = &Account{}

func AccountInstance() *Account {
	return defaultAccount
}

// 初始化
func (s *Account) Init(dir string, node *service.ServiceNode, auditeTemplateName, userFrozen1ForAdminTempName, userFrozen1ForUserTempName string) {
	var err error
	s.privateKey, err = ioutil.ReadFile(dir + "/" + config.BastionPayPrivateKey)
	if err != nil {
		l4g.Crashf("", err)
	}
	s.serverPublicKey, err = ioutil.ReadFile(dir + "/" + config.BastionPayPublicKey)
	if err != nil {
		l4g.Crashf("", err)
	}
	s.auditeTemplateName = auditeTemplateName
	s.userFrozen1ForAdminTempName = userFrozen1ForAdminTempName
	s.userFrozen1ForUserTempName = userFrozen1ForUserTempName
	s.node = node
	l4g.Info("auditeTemplateName[%s]userFrozen1ForAdminTempName[%s]userFrozen1ForUserTempName[%s]", s.auditeTemplateName, s.userFrozen1ForAdminTempName, s.userFrozen1ForUserTempName)
}

func (s *Account) GetApiGroup() map[string]service.NodeApi {
	nam := make(map[string]service.NodeApi)

	func() {
		service.RegisterApi(&nam,
			"register", data.APILevel_genesis, s.Register)
	}()

	func() {
		service.RegisterApi(&nam,
			"updateprofile", data.APILevel_admin, s.UpdateProfile)
	}()

	func() {
		service.RegisterApi(&nam,
			"adminupdateprofile", data.APILevel_genesis, s.UpdateProfileAdmin)
	}()

	func() {
		service.RegisterApi(&nam,
			"readprofile", data.APILevel_admin, s.ReadProfile)
	}()

	func() {
		service.RegisterApi(&nam,
			"listusers", data.APILevel_admin, s.ListUsers)

	}()

	func() {
		service.RegisterApi(&nam,
			"updatefrozen", data.APILevel_admin, s.UpdateFrozen)

	}()
	func() {
		service.RegisterApi(&nam,
			"userfrozen", data.APILevel_client, s.UserFrozen)

	}()

	func() {
		service.RegisterApi(&nam,
			"updateaudite", data.APILevel_admin, s.UpdateAuditeStatus)
	}()

	func() {
		service.RegisterApi(&nam,
			"getaudite", data.APILevel_admin, s.GetAuditeStatus)
	}()

	func() {
		service.RegisterApi(&nam,
			"getuserstatus", data.APILevel_admin, s.GetUserAllStatus)
	}()

	func() {
		service.RegisterApi(&nam,
			"updatetransfer", data.APILevel_admin, s.UpdateTransferStatus)
	}()

	return nam
}

func (s *Account) HandleNotify(req *data.SrvRequest) {
	l4g.Info("HandleNotify-reloadUserLevel: do nothing")
}

// 创建账号
func (s *Account) Register(req *data.SrvRequest, res *data.SrvResponse) {
	// from req
	reqUserRegister := backend.ReqUserRegister{}
	err := json.Unmarshal([]byte(req.Argv.Message), &reqUserRegister)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		return
	}

	// userkey
	uuid, err := uuid.NewV4()
	if err != nil {
		l4g.Error("error create user key: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}
	userKey := uuid.String()

	// db
	err = db.Register(&reqUserRegister, userKey, backend.AUDITE_Status_Blank, 0)
	if err != nil {
		l4g.Error("error create user: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// to ack
	ackUserCreate := backend.AckUserRegister{}
	ackUserCreate.UserKey = userKey

	dataAck, err := json.Marshal(ackUserCreate)
	if err != nil {
		db.Delete(userKey)
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("create a new user: %s", res.Value.Message)
}

// 获取用户列表
// 登入
func (s *Account) ListUsers(req *data.SrvRequest, res *data.SrvResponse) {
	// from req
	reqUserList := struct {
		TotalLines   int                    `json:"total_lines" doc:"总数,0：表示首次查询"`
		PageIndex    int                    `json:"page_index" doc:"页索引,1开始"`
		MaxDispLines int                    `json:"max_disp_lines" doc:"页最大数，100以下"`
		Condition    map[string]interface{} `json:"condition" doc:"条件查询"`
	}{}
	err := json.Unmarshal([]byte(req.Argv.Message), &reqUserList)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		return
	}

	var (
		pageNum   = reqUserList.MaxDispLines
		totalLine = reqUserList.TotalLines
		pageIndex = reqUserList.PageIndex

		beginIndex = 0
	)

	if reqUserList.Condition != nil {
		v, ok := reqUserList.Condition["transfer_status"]
		if ok {
			delete(reqUserList.Condition, "transfer_status")
			reqUserList.Condition["can_transfer"] = v
		}
	}

	if totalLine == 0 {
		//totalLine, err = db.ListUserCount()
		totalLine, err = db.ListUserCountByBasic(reqUserList.Condition)
		if err != nil {
			l4g.Error("error json message: %s", err.Error())
			res.Err = apibackend.ErrAccountSrvListUsersCount
			return
		}
	}

	if pageNum < 1 || pageNum > 100 {
		pageNum = 50
	}

	beginIndex = pageNum * (pageIndex - 1)

	//ackUserList, err := db.ListUsers(beginIndex, pageNum)
	ackUserList, err := db.ListUsersByBasic(beginIndex, pageNum, reqUserList.Condition)
	if err != nil {
		l4g.Error("error ListUsers: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvListUsers
		return
	}

	ackUserList.PageIndex = pageIndex
	ackUserList.MaxDispLines = pageNum
	ackUserList.TotalLines = totalLine

	// to ack
	dataAck, err := json.Marshal(ackUserList)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("list users: %s", res.Value.Message)
}

// 获取key
func (s *Account) ReadProfile(req *data.SrvRequest, res *data.SrvResponse) {
	// from req
	//reqReadProfile := v1.ReqUserReadProfile{}
	//err := json.Unmarshal([]byte(req.Argv.Message), &reqReadProfile)
	//if err != nil {
	//	l4g.Error("error json message: %s", err.Error())
	//	res.Err = data.ErrDataCorrupted
	//	return
	//}

	// load profile
	ackReadProfile, err := db.ReadProfile(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("error ReadProfile[%s]: %s", req.Argv.SubUserKey, err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		return
	}

	if ackReadProfile.PublicKey != "" && ackReadProfile.CallbackUrl != "" && ackReadProfile.SourceIP != "" {
		ackReadProfile.ServerPublicKey = string(s.serverPublicKey)
	}

	// to ack
	dataAck, err := json.Marshal(ackReadProfile)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("read a user profile: %s", res.Value.Message)
}

func (s *Account) UpdateProfile(req *data.SrvRequest, res *data.SrvResponse) {
	ul, err := db.ReadUserLevel(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("dbGetAuditeStatus subuserkey[%s] err[%s]", req.Argv.SubUserKey, err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		res.ErrMsg = "db.GetAuditeStatus " + err.Error()
		return
	}

	if ul.AuditeStatus == backend.AUDITE_Status_Pass || ul.AuditeStatus == backend.AUDITE_Status_Doing {
		l4g.Error("error NoAdmin AuditeStatus no support update")
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus no support update"
		return
	}

	s.updateProfileInner(req, res, ul, backend.AUDITE_Status_Doing)

	go s.notifyAuditeDoing(req, res)
}

func (s *Account) updateProfileInner(req *data.SrvRequest, res *data.SrvResponse, ul *db.UserLevel, newAudite int) {
	reqUpdateProfile := backend.ReqUserUpdateProfile{}
	err := json.Unmarshal([]byte(req.Argv.Message), &reqUpdateProfile)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		return
	}

	err = utils.RsaVerifyPubKey([]byte(reqUpdateProfile.PublicKey))
	if err != nil {
		l4g.Error("pub key parse: %s", err.Error())
		res.Err = apibackend.ErrAccountPubKeyParse
		return
	}

	// load old key
	oldUserReadProfile, err := db.ReadProfile(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("error ReadProfile: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		return
	}

	if reqUpdateProfile.PublicKey == "" {
		reqUpdateProfile.PublicKey = oldUserReadProfile.PublicKey
	}
	if reqUpdateProfile.SourceIP == "" {
		reqUpdateProfile.SourceIP = oldUserReadProfile.SourceIP
	}
	if reqUpdateProfile.CallbackUrl == "" {
		reqUpdateProfile.CallbackUrl = oldUserReadProfile.CallbackUrl
	}

	// update key
	if err := db.UpdateProfile(req.Argv.SubUserKey, &reqUpdateProfile, newAudite); err != nil {
		l4g.Error("error update profile: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvUpdateProfile
		return
	}

	// to ack
	ackUpdateProfile := backend.AckUserUpdateProfile{
		ServerPublicKey: string(s.serverPublicKey),
	}
	dataAck, err := json.Marshal(ackUpdateProfile)
	if err != nil {
		// 写回去
		oldUserUpdateProfile := backend.ReqUserUpdateProfile{}
		oldUserUpdateProfile.PublicKey = oldUserReadProfile.PublicKey
		oldUserUpdateProfile.SourceIP = oldUserReadProfile.SourceIP
		oldUserUpdateProfile.CallbackUrl = oldUserReadProfile.CallbackUrl
		db.UpdateProfile(req.Argv.SubUserKey, &oldUserUpdateProfile, ul.AuditeStatus)
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("update a user profile: %s", res.Value.Message)

	// notify
	func() {
		notifyReq := data.SrvRequest{}
		notifyReq.Method.Version = "v1"
		notifyReq.Method.Srv = "account"
		notifyReq.Method.Function = "updateprofile"
		notifyReq.Argv.UserKey = ""
		notifyReq.Argv.SubUserKey = req.Argv.SubUserKey

		notifyRes := data.SrvResponse{}
		s.node.InnerNotify(&notifyReq, &notifyRes)

		l4g.Info("notify a user profile: %s", req.Argv.SubUserKey)
	}()
}

// 更新key
func (s *Account) UpdateProfileAdmin(req *data.SrvRequest, res *data.SrvResponse) {
	ul, err := db.ReadUserLevel(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("dbGetAuditeStatus subuserkey[%s] err[%s]", req.Argv.SubUserKey, err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		res.ErrMsg = "db.GetAuditeStatus " + err.Error()
		return
	}
	s.updateProfileInner(req, res, ul, ul.AuditeStatus)
	// from req

}

func (s *Account) UpdateAuditeStatus(req *data.SrvRequest, res *data.SrvResponse) {
	reqUserAuditeStatus := new(backend.ReqUserAuditeStatus)
	err := json.Unmarshal([]byte(req.Argv.Message), reqUserAuditeStatus)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		res.ErrMsg = err.Error()
		return
	}
	newAuditeStatus := reqUserAuditeStatus.AuditeStatus
	if (newAuditeStatus == backend.AUDITE_Status_Deny) && IsBlankAuditeInfo(reqUserAuditeStatus.AuditeInfo) {
		l4g.Error("param nohave audite_info")
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "param loss audite_info"
		return
	}
	oldstatus, _, err := db.GetAuditeStatus(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("error dbGetAuditeStatus: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		res.ErrMsg = "db.GetAuditeStatus " + err.Error()
		return
	}

	if (newAuditeStatus == backend.AUDITE_Status_Blank) && (oldstatus != backend.AUDITE_Status_Blank) {
		l4g.Error("AuditeStatus failed oldstatus[%d] newstatus[%d]", oldstatus, newAuditeStatus)
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus InValid"
		return
	}

	if (newAuditeStatus == backend.AUDITE_Status_Pass || newAuditeStatus == backend.AUDITE_Status_Deny) && (oldstatus != backend.AUDITE_Status_Doing) {
		l4g.Error("AuditeStatus failed oldstatus[%d] newstatus[%d]", oldstatus, newAuditeStatus)
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus InValid"
		return
	}
	if (newAuditeStatus == backend.AUDITE_Status_Doing) && !(oldstatus == backend.AUDITE_Status_Deny || oldstatus == backend.AUDITE_Status_Blank) {
		l4g.Error("AuditeStatus failed oldstatus[%d] newstatus[%d]", oldstatus, newAuditeStatus)
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus InValid"
		return
	}

	err = db.UpdateAuditeStatus(req.Argv.SubUserKey, reqUserAuditeStatus.AuditeStatus, reqUserAuditeStatus.AuditeInfo)
	if err != nil {
		l4g.Error("error UpdateAuditeStatus: %s", err.Error())
		res.Err = apibackend.ErrInternal
		res.ErrMsg = "db.UpdateAuditeStatus " + err.Error()
		return
	}
	res.Err = apibackend.NoErr

	if newAuditeStatus == backend.AUDITE_Status_Doing {
		go s.notifyAuditeDoing(req, res)
	}
	l4g.Info("UpdateAuditeStatus Success, subuserkey[%s] oldstatus[%d]newAuditeStatus[%d]", req.Argv.SubUserKey, oldstatus, newAuditeStatus)
	func() {
		notifyReq := data.SrvRequest{}
		notifyReq.Method.Version = "v1"
		notifyReq.Method.Srv = "account"
		notifyReq.Method.Function = "updateaudite"
		notifyReq.Argv.UserKey = ""
		notifyReq.Argv.SubUserKey = req.Argv.SubUserKey

		notifyRes := data.SrvResponse{}
		s.node.InnerNotify(&notifyReq, &notifyRes)

		l4g.Info("notify a user audite: %s", req.Argv.SubUserKey)
	}()
}

func (s *Account) UpdateTransferStatus(req *data.SrvRequest, res *data.SrvResponse) {
	reqUserStatus := new(backend.ReqUserTransferStatus)
	err := json.Unmarshal([]byte(req.Argv.Message), reqUserStatus)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		res.ErrMsg = err.Error()
		return
	}
	err = db.UpdateTransferStatus(req.Argv.SubUserKey, reqUserStatus.TransferStatus)
	if err != nil {
		l4g.Error("error UpdateCanTransfer: %s", err.Error())
		res.Err = apibackend.ErrInternal
		res.ErrMsg = "db.UpdateCanTransfer " + err.Error()
		return
	}
	l4g.Info("userkey[%s] UpdateTransferStatus[%d]", req.Argv.SubUserKey, reqUserStatus.TransferStatus)
	res.Err = apibackend.NoErr
}

func (s *Account) GetAuditeStatus(req *data.SrvRequest, res *data.SrvResponse) {
	status, info, err := db.GetAuditeStatus(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("GetAuditeStatus userkey[%s] subuserkey[%s] dataFrom[%d]error[%s]", req.Argv.UserKey, req.Argv.SubUserKey, req.Context.DataFrom, err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		return
	}
	//	status := uint(1)
	resAudite := new(backend.ResUserAuditeStatus)
	resAudite.AuditeStatus = status

	if len(info) != 0 {
		resAudite.AuditeInfo = new(string)
		*resAudite.AuditeInfo = info
	}

	content, err := json.Marshal(resAudite)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}
	res.Err = apibackend.NoErr
	res.Value.Message = string(content)
}

//包括 audite, can_transfer
func (s *Account) GetUserAllStatus(req *data.SrvRequest, res *data.SrvResponse) {
	info, err := db.GetUserAllStatus(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("GetAuditeStatus userkey[%s] subuserkey[%s] dataFrom[%d]error[%s]", req.Argv.UserKey, req.Argv.SubUserKey, req.Context.DataFrom, err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		return
	}

	content, err := json.Marshal(info)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}
	res.Err = apibackend.NoErr
	res.Value.Message = string(content)
}

// 设置冻结
func (s *Account) UpdateFrozen(req *data.SrvRequest, res *data.SrvResponse) {
	// from req
	reqFrozeUser := backend.ReqFrozenUser{}
	err := json.Unmarshal([]byte(req.Argv.Message), &reqFrozeUser)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		return
	}

	oldstatus, _, err := db.GetAuditeStatus(reqFrozeUser.UserKey)
	if err != nil {
		l4g.Error("error dbGetAuditeStatus: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		res.ErrMsg = "db.GetAuditeStatus " + err.Error()
		return
	}
	if oldstatus != backend.AUDITE_Status_Pass {
		l4g.Error("error AuditeStatus no support UpdateFrozen")
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus no support UpdateFrozen"
		return
	}

	// set frozen
	err = db.UpdateFrozen(reqFrozeUser.UserKey, reqFrozeUser.IsFrozen)
	if err != nil {
		l4g.Error("error UpdateFrozen: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvSetFrozen
		return
	}

	// get frozen
	ackFrozenUser := backend.AckFrozenUser{}
	ackFrozenUser.UserKey = reqFrozeUser.UserKey
	ackFrozenUser.IsFrozen, err = db.ReadFrozen(ackFrozenUser.UserKey)
	if err != nil {
		l4g.Error("error UpdateFrozen: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvSetFrozen
		return
	}

	// to ack
	dataAck, err := json.Marshal(ackFrozenUser)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("update a user frozen: %s", res.Value.Message)

	// notify
	func() {
		notifyReq := data.SrvRequest{}
		notifyReq.Method.Version = "v1"
		notifyReq.Method.Srv = "account"
		notifyReq.Method.Function = "updatefrozen"
		notifyReq.Argv.UserKey = ""
		notifyReq.Argv.SubUserKey = ackFrozenUser.UserKey

		notifyRes := data.SrvResponse{}
		s.node.InnerNotify(&notifyReq, &notifyRes)

		l4g.Info("notify a user frozen: %s", req.Argv.SubUserKey)
	}()
}

// 设置冻结
func (s *Account) UserFrozen(req *data.SrvRequest, res *data.SrvResponse) {
	// from req
	reqFrozeUser := backend.ReqFrozenUser{}
	err := json.Unmarshal([]byte(req.Argv.Message), &reqFrozeUser)
	if err != nil {
		l4g.Error("error json message: %s", err.Error())
		res.Err = apibackend.ErrDataCorrupted
		return
	}
	if reqFrozeUser.IsFrozen == 0 {
		return
	}

	status, err := db.GetUserAllStatus(reqFrozeUser.UserKey)
	if err != nil {
		l4g.Error("error GetUserAllStatus: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvNoUser
		res.ErrMsg = "db.GetUserAllStatus " + err.Error()
		return
	}
	if status.AuditeStatus != backend.AUDITE_Status_Pass {
		l4g.Error("error AuditeStatus no support UpdateFrozen")
		res.Err = apibackend.ErrAccountSrvAudite
		res.ErrMsg = "AuditeStatus no support UpdateFrozen"
		return
	}

	if status.IsFrozen == 1 {
		return
	}

	// set frozen
	err = db.UpdateFrozen(reqFrozeUser.UserKey, reqFrozeUser.IsFrozen)
	if err != nil {
		l4g.Error("error UserFrozen: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvSetFrozen
		return
	}

	// get frozen
	ackFrozenUser := backend.AckFrozenUser{}
	ackFrozenUser.UserKey = reqFrozeUser.UserKey
	ackFrozenUser.IsFrozen, err = db.ReadFrozen(ackFrozenUser.UserKey)
	if err != nil {
		l4g.Error("error UserFrozen: %s", err.Error())
		res.Err = apibackend.ErrAccountSrvSetFrozen
		return
	}

	// to ack
	dataAck, err := json.Marshal(ackFrozenUser)
	if err != nil {
		l4g.Error("error Marshal: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value.Message = string(dataAck)
	l4g.Info("userfrozen itself[%s]: %s", reqFrozeUser.UserKey, res.Value.Message)
	go s.notifyUserFrozen(reqFrozeUser.UserKey)
	// notify
	func() {
		notifyReq := data.SrvRequest{}
		notifyReq.Method.Version = "v1"
		notifyReq.Method.Srv = "account"
		notifyReq.Method.Function = "userfrozen"
		notifyReq.Argv.UserKey = ""
		notifyReq.Argv.SubUserKey = ackFrozenUser.UserKey

		notifyRes := data.SrvResponse{}
		s.node.InnerNotify(&notifyReq, &notifyRes)

		l4g.Info("notify a user frozen: %s", req.Argv.SubUserKey)
	}()
}

func (s *Account) notifyAuditeDoing(req *data.SrvRequest, res *data.SrvResponse) {
	if res.Err != apibackend.NoErr {
		return
	}
	name, err := db.ReadUserName(req.Argv.SubUserKey)
	if err != nil {
		l4g.Error("db.ReadUserName[%s] err[%s]", req.Argv.SubUserKey, err.Error())
		return
	}
	err = sdkmail.GNotifySdk.SendSmsByGroupName(s.auditeTemplateName, "zh-CN", nil, map[string]interface{}{"key1": name})
	if err != nil {
		l4g.Error("SendSmsByGroupName[%s][%s] err[%s]", req.Argv.SubUserKey, name, err.Error())
		return
	}
	l4g.Info("SendSmsByGroupName[%s][%s] success", req.Argv.SubUserKey, name)
}

func (this *Account) notifyUserFrozen(userKey string) {
	if len(userKey) == 0 {
		l4g.Error("userkey[%s] err[is nil]", userKey)
		return
	}

	user, err := db.ReadUserInfo(userKey)
	if err != nil {
		l4g.Error("(%s) get user info failed: %s", userKey, err.Error())
		return
	}
	l4g.Info("user[%s][%s] start notify UserFrozen", user.UserName, userKey)
	if len(this.userFrozen1ForAdminTempName) == 0 {
		l4g.Error("user[%s][%s] notify err[admin template is nil]", user.UserName, userKey)
	} else {
		//发送给管理员的
		err = sdkmail.GNotifySdk.SendSmsByGroupName(this.userFrozen1ForAdminTempName, "zh-CN", nil, map[string]interface{}{"key1": user.UserName})
		if err != nil {
			l4g.Error("admin SendSmsByGroupName[%s][%s] err[%s]", user.UserName, userKey, err.Error())
			//发送给管理员 失败，可能是 手机号错了，不能返回，要继续
		}
	}
	if len(this.userFrozen1ForUserTempName) == 0 {
		l4g.Error("user[%s][%s] notify err[user template is nil]", user.UserName, userKey)
		return
	}

	//发送给用户的
	if (len(user.CountryCode) != 0) && (len(user.UserMobile) != 0) {
		err = sdkmail.GNotifySdk.SendSmsByGroupName(this.userFrozen1ForUserTempName, "en-US", []string{user.CountryCode + user.UserMobile}, nil)
		if err != nil {
			l4g.Error("user SendSmsByGroupName[%s][%s][%s] err[%s]", user.UserName, userKey, user.CountryCode+user.UserMobile, err.Error())
		}
	}
	if len(user.UserEmail) != 0 {
		err = sdkmail.GNotifySdk.SendMailByGroupName(this.userFrozen1ForUserTempName, "en-US", []string{user.UserEmail}, nil)
		if err != nil {
			l4g.Error("user SendMailByGroupName[%s][%s][%s] err[%s]", user.UserName, userKey, user.UserEmail, err.Error())
		}
	}
}
