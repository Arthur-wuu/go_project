package handler

import (
	"BastionPay/bas-auth-srv/db"
	"BastionPay/bas-base/data"
	//"github.com/BastionPay/bas-service/base/service"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	"BastionPay/bas-base/config"
	service "BastionPay/bas-base/service2"
	l4g "github.com/alecthomas/log4go"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
	"BastionPay/bas-tools/sdk.notify.mail"
)

type Auth struct {
	node *service.ServiceNode

	privateKey []byte

	intervalSecond int

	rwmu       sync.RWMutex
	usersLevel map[string]*db.UserLevel
	mNoValidIpNotifyLimit     *limiter.Limiter //ip不匹配时，限制短信通知的频率
	mNoValidIpNotifyTemplateForAdmin  string
	mNoValidIpNotifyTemplateForUser  string
}

var defaultAuth = &Auth{}

func AuthInstance() *Auth {
	return defaultAuth
}

func (auth *Auth) Init(cfgPath, dir string, node *service.ServiceNode) {
	var (
		err      error
		interval string
	)
	auth.privateKey, err = ioutil.ReadFile(dir + "/" + config.BastionPayPrivateKey)
	if err != nil {
		l4g.Crashf("", err)
	}

	auth.intervalSecond = 30
	err = config.LoadJsonNode(cfgPath, "interval", &interval)
	if err == nil && interval != "" {
		auth.intervalSecond, err = strconv.Atoi(interval)
		if err != nil {
			l4g.Crashf("", err)
		}
	}
	l4g.Info("interval:%d", auth.intervalSecond)

	auth.node = node
	auth.usersLevel = make(map[string]*db.UserLevel)

	noValidIpRate := ""
	err = config.LoadJsonNode(cfgPath, "novalid_ip_notify_rate", &noValidIpRate)
	if err != nil || noValidIpRate == "" {
		l4g.Error("config no find novalid_ip_notify_rate err[%v]", err)
	}else{
		rate, err := limiter.NewRateFromFormatted(noValidIpRate)
		if err != nil {
			l4g.Error("limiter NewRateFromFormatted err[%s]", err.Error())
			return
		}
		store := memory.NewStore()
		auth.mNoValidIpNotifyLimit = limiter.New(store, rate)
	}

	err = config.LoadJsonNode(cfgPath, "novalid_ip_temp", &auth.mNoValidIpNotifyTemplateForAdmin)
	if err != nil || auth.mNoValidIpNotifyTemplateForAdmin == "" {
		l4g.Error("config no find novalid_ip_temp err[%v]", err)
	}
	err = config.LoadJsonNode(cfgPath, "novalid_ip_temp_foruser", &auth.mNoValidIpNotifyTemplateForUser)
	if err != nil || auth.mNoValidIpNotifyTemplateForUser == "" {
		l4g.Error("config no find novalid_ip_temp err[%v]", err)
	}

	bas_notifyAddr := ""
	err = config.LoadJsonNode(cfgPath, "bas_notify", &bas_notifyAddr)
	if err != nil || bas_notifyAddr == "" {
		l4g.Error("config no find bas_notify err[%v]", err)
	}else{
		if err := sdk_notify_mail.GNotifySdk.Init(bas_notifyAddr, "bas-auth-api");err != nil {
			l4g.Error("bas_notify Init err[%v]", err)
		}
	}
	l4g.Info("bas_notify[%s]novalid_ip_temp[%s]novalid_ip_notify_rate[%s]novalid_ip_temp_foruser[%s]", bas_notifyAddr, auth.mNoValidIpNotifyTemplateForAdmin, noValidIpRate, auth.mNoValidIpNotifyTemplateForUser)

}

func (auth *Auth) getUserLevel(userKey string) (*db.UserLevel, error) {
	ul := func() *db.UserLevel {
		auth.rwmu.RLock()
		defer auth.rwmu.RUnlock()

		return auth.usersLevel[userKey]
	}()
	if ul != nil {
		return ul, nil
	}

	return func() (*db.UserLevel, error) {
		auth.rwmu.Lock()
		defer auth.rwmu.Unlock()

		ul := auth.usersLevel[userKey]
		if ul != nil {
			return ul, nil
		}
		ul, err := db.ReadUserLevel(userKey)
		if err != nil {
			return nil, err
		}
		auth.usersLevel[userKey] = ul
		return ul, nil
	}()
}

func (auth *Auth) reloadUserLevel(userKey string) (*db.UserLevel, error) {
	auth.rwmu.Lock()
	defer auth.rwmu.Unlock()

	ul, err := db.ReadUserLevel(userKey)
	if err != nil {
		return nil, err
	}
	auth.usersLevel[userKey] = ul
	return ul, nil
}

func (auth *Auth) GetApiGroup() map[string]service.NodeApi {
	nam := make(map[string]service.NodeApi)

	func() {
		service.RegisterApi(&nam,
			"authdata", data.APILevel_client, auth.AuthData)
	}()

	func() {
		service.RegisterApi(&nam,
			"encryptdata", data.APILevel_client, auth.EncryptData)
	}()

	return nam
}

func (auth *Auth) HandleNotify(req *data.SrvRequest) {
	if req.Method.Srv == "account" {
		//reqUpdateProfile := v1.ReqUserUpdateProfile{}
		//err := json.Unmarshal([]byte(req.Argv.Message), &reqUpdateProfile)
		//if err != nil {
		//	l4g.Error("HandleNotify-Unmarshal: %s", err.Error())
		//	return
		//}
		if req.Method.Function == "updateprofile" || req.Method.Function == "updatefrozen" || req.Method.Function == "updateaudite" || req.Method.Function == "userfrozen" {
			// reload profile
			ndata, err := auth.reloadUserLevel(req.Argv.SubUserKey)
			if err != nil {
				l4g.Error("HandleNotify-reloadUserLevel: %s", err.Error())
				return
			}

			l4g.Info("HandleNotify-reloadUserLevel: ", req.Argv.SubUserKey, ndata)
		}
	}
}

func (auth *Auth) isRequestValid(req *data.SrvRequest) bool {
	if auth.intervalSecond <= 0 {
		return true
	}

	iTime, err := strconv.ParseInt(req.Argv.TimeStamp, 10, 64)
	if err != nil {
		return false
	}
	uTime := time.Unix(iTime, 0)
	uDuration := time.Now().Sub(uTime)

	if math.Abs(uDuration.Seconds()) > float64(auth.intervalSecond)*time.Second.Seconds() {
		return false
	}

	return true
}

// 验证数据
func (auth *Auth) AuthData(req *data.SrvRequest, res *data.SrvResponse) {
	if auth.isRequestValid(req) == false {
		l4g.Error("(%s) get user request invalid failed: %s", req.Argv.UserKey, req.Argv.TimeStamp)
		res.Err = apibackend.ErrRequestInvalid
		return
	}

	ul, err := auth.getUserLevel(req.Argv.UserKey)
	if err != nil {
		l4g.Error("(%s) get user level failed: %s", req.Argv.UserKey, err.Error())
		res.Err = apibackend.ErrAuthSrvNoUserKey
		return
	}

	if ul.PublicKey == "" {
		l4g.Error("(%s-%s) failed: no public key", req.Argv.UserKey, req.Method.Function)
		res.Err = apibackend.ErrAuthSrvNoPublicKey
		return
	}

	if req.Context.ApiLever > ul.Level {
		l4g.Error("(%s-%s) failed: no api level", req.Argv.UserKey, req.Method.Function)
		l4g.Error("%d %d ", req.Context.ApiLever, ul.Level)
		res.Err = apibackend.ErrAuthSrvNoApiLevel
		return
	}
	l4g.Debug("DataFrom[%d] AuditeStatus[%d]", req.Context.DataFrom, ul.AuditeStatus)
	if (req.Context.DataFrom == apibackend.DataFromApi) && !auth.IsValidAudite(ul) {
		l4g.Error("(%s-%s) failed: user Audite", req.Argv.UserKey, req.Method.Function)
		res.Err = apibackend.ErrAuthSrvIllegalAudite
		return
	}

	if req.Context.ApiLever > ul.Level || ul.IsFrozen != 0 {
		l4g.Error("(%s-%s) failed: user frozen", req.Argv.UserKey, req.Method.Function)
		res.Err = apibackend.ErrAuthSrvUserFrozen
		return
	}

	l4g.Debug("httpIp[%v] DbSourceIP[%s] userkey[%s]", req.Context.SourceIp, ul.SourceIP, req.Argv.UserKey)
	if req.Context.SourceIp != nil { //兼容老版gateway
		whiteOk,cleanCtxSrcIp := auth.IsWhiteList(ul.SourceIP, *req.Context.SourceIp)
		if (req.Context.DataFrom == apibackend.DataFromApi) && !whiteOk {
			l4g.Error("UserSourceIp[%s] reqSourceIp[%s] userkey[%s] err:not equal", ul.SourceIP, *req.Context.SourceIp, req.Argv.UserKey)
			res.Err = apibackend.ErrAuthSrvIllegalIp
			go auth.VaildIpNotify(req.Argv.UserKey ,cleanCtxSrcIp)
			return
		}
	}

	if req.Context.DataFrom == apibackend.DataFromUser || req.Context.DataFrom == apibackend.DataFromAdmin {
		if ul.UserClass != data.UserClass_Admin {
			l4g.Error("%s illegally call data type %d", req.Argv.UserKey, ul.UserClass)
			res.Err = apibackend.ErrAuthSrvIllegalDataType
			return
		}
	}

	originData, err := data.DecryptionAndVerifyData(&req.Argv, []byte(ul.PublicKey), auth.privateKey)
	if err != nil {
		l4g.Error("DecryptionAndVerifyData: %s", err.Error())
		res.Err = apibackend.ErrAuthSrvIllegalData
		return
	}

	res.Value.Message = string(originData)
}

// 打包数据
func (auth *Auth) EncryptData(req *data.SrvRequest, res *data.SrvResponse) {
	ul, err := auth.getUserLevel(req.Argv.UserKey)
	if err != nil {
		l4g.Error("(%s) get user level failed: %s", req.Argv.UserKey, err.Error())
		res.Err = apibackend.ErrAuthSrvNoUserKey
		return
	}

	if ul.PublicKey == "" {
		l4g.Error("(%s-%s) failed: no public key", req.Argv.UserKey, req.Method.Function)
		res.Err = apibackend.ErrAuthSrvNoPublicKey
		return
	}

	timeStamp := ""
	if auth.intervalSecond > 0 {
		timeStamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	srvData, err := data.EncryptionAndSignData([]byte(req.Argv.Message), timeStamp, req.Argv.UserKey, []byte(ul.PublicKey), auth.privateKey)
	if err != nil {
		l4g.Error("EncryptionAndSignData: %s", err.Error())
		res.Err = apibackend.ErrInternal
		return
	}

	// ok
	res.Value = *srvData
}

func (auth *Auth) IsValidAudite(ul *db.UserLevel) bool {
	if ul == nil {
		return false
	}
	if ul.AuditeStatus == backend.AUDITE_Status_Pass {
		return true
	}
	return false
}

func (auth *Auth) IsWhiteList(whiteLsit, srcIp string) (bool,string) {
	srcIp = strings.TrimLeft(srcIp, "http://")
	srcIp = strings.TrimLeft(srcIp, "https://")
	arr0 := strings.Split(srcIp, ",")
	if len(arr0) > 1 {
		srcIp = arr0[0]
	}
	srcIp = strings.TrimSpace(srcIp)
	arr := strings.Split(srcIp, ":")
	if len(arr) > 1 {
		srcIp = arr[0]
	}
	whiteLsit = strings.TrimSpace(whiteLsit)
	whiteLsitArr := strings.Split(whiteLsit, ",")
	if (len(whiteLsitArr) == 1) && (whiteLsitArr[0] == "*.*.*.*"){
		return true,srcIp
	}
	for i := 0; i < len(whiteLsitArr); i++ {
		if srcIp == whiteLsitArr[i] {
			return true,srcIp
		}
	}
	return false,srcIp
}

func (auth *Auth) VaildIpNotify(userKey ,srcIp string) {
	if len(userKey) == 0 {
		l4g.Error("userkey[%s] err[is nil]", userKey)
		return
	}
	flag, err := auth.OverNoVaildIpNotifyLimit(userKey)
	if err != nil {
		l4g.Error("OverNoVaildIpNotifyLimit userkey[%s]err[%s]", userKey, err.Error())
		return
	}
	if flag {
		return
	}
	user, err := auth.getUserLevel(userKey)
	if err != nil {
		l4g.Error("(%s) get user level failed: %s", userKey, err.Error())
		return
	}
	l4g.Info("user[%s][%s] start notify ip[%s] mobile[%s]email[%s]", user.UserName, userKey, srcIp,user.CountryCode+user.UserMobile, user.UserEmail)
	if len(auth.mNoValidIpNotifyTemplateForAdmin) == 0 {
		l4g.Error("user[%s][%s] notify err[admin template is nil]", user.UserName, userKey)
	}else {
		//发送给管理员的
		err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(auth.mNoValidIpNotifyTemplateForAdmin, "zh-CN", nil, map[string]interface{}{"key1":user.UserName, "key2":srcIp})
		if err != nil {
			l4g.Error("admin SendSmsByGroupName[%s][%s] err[%s]", user.UserName, userKey, err.Error())
			//发送给管理员 失败，可能是 手机号错了，不能返回，要继续
		}else{
			l4g.Info("admin SendSmsByGroupName[%s][%s] success", user.UserName, userKey)
		}
	}
	if len(auth.mNoValidIpNotifyTemplateForUser) == 0 {
		l4g.Error("user[%s][%s] notify err[user template is nil]", user.UserName, userKey)
		return
	}

	//发送给用户的
	if (len(user.CountryCode) != 0) && (len(user.UserMobile) != 0) {
		err = sdk_notify_mail.GNotifySdk.SendSmsByGroupName(auth.mNoValidIpNotifyTemplateForUser, "en-US", []string{user.CountryCode+user.UserMobile}, map[string]interface{}{"key1":user.UserName, "key2":srcIp})
		if err != nil {
			l4g.Error("user SendSmsByGroupName[%s][%s][%s] err[%s]", user.UserName, userKey,user.CountryCode+user.UserMobile, err.Error())
		}else{
			l4g.Info("user SendSmsByGroupName[%s][%s][%s] success", user.UserName, userKey, user.CountryCode+user.UserMobile)
		}
	}
	if len(user.UserEmail) != 0 {
		err = sdk_notify_mail.GNotifySdk.SendMailByGroupName(auth.mNoValidIpNotifyTemplateForUser, "en-US", []string{user.UserEmail}, map[string]interface{}{"key1":user.UserName, "key2":srcIp})
		if err != nil {
			l4g.Error("user SendMailByGroupName[%s][%s][%s] err[%s]", user.UserName, userKey,user.UserEmail, err.Error())
		}else{
			l4g.Info("user SendMailByGroupName[%s][%s][%s] success", user.UserName, userKey, user.UserEmail)
		}
	}
}

func (auth *Auth) OverNoVaildIpNotifyLimit(key string) (bool,error) {
	ctx, err := auth.mNoValidIpNotifyLimit.Get(nil, key)
	if err != nil {
		return true, err
	}
	if ctx.Reached {
		return true,nil
	}
	return false,nil
}