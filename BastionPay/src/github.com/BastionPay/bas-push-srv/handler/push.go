package handler

import (
	"BastionPay/bas-push-srv/db"
	"github.com/BastionPay/bas-base/data"
	//"github.com/BastionPay/bas-service/base/service"
	"encoding/json"
	"github.com/BastionPay/bas-api/api"
	"github.com/BastionPay/bas-api/apibackend"
	"github.com/BastionPay/bas-api/apibackend/v1/backend"
	"github.com/BastionPay/bas-base/config"
	"github.com/BastionPay/bas-base/nethelper"
	service "github.com/BastionPay/bas-base/service2"
	l4g "github.com/alecthomas/log4go"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

type Push struct {
	privateKey []byte

	intervalSecond int

	rwmu             sync.RWMutex
	usersCallbackUrl map[string]*backend.AckUserReadProfile
}

var defaultPush = &Push{}

func PushInstance() *Push {
	return defaultPush
}

func (push *Push) Init(cfgPath, dir string) {
	var (
		err      error
		interval string
	)
	push.privateKey, err = ioutil.ReadFile(dir + "/" + config.BastionPayPrivateKey)
	if err != nil {
		l4g.Crashf("", err)
	}

	push.intervalSecond = 30
	err = config.LoadJsonNode(cfgPath, "interval", &interval)
	if err == nil && interval != "" {
		push.intervalSecond, err = strconv.Atoi(interval)
		if err != nil {
			l4g.Crashf("", err)
		}
	}
	l4g.Info("interval:%d", push.intervalSecond)

	push.usersCallbackUrl = make(map[string]*backend.AckUserReadProfile)
}

func (push *Push) getUserProfile(userKey string) (*backend.AckUserReadProfile, error) {
	rp := func() *backend.AckUserReadProfile {
		push.rwmu.RLock()
		defer push.rwmu.RUnlock()

		return push.usersCallbackUrl[userKey]
	}()
	if rp != nil {
		return rp, nil
	}

	return func() (*backend.AckUserReadProfile, error) {
		push.rwmu.Lock()
		defer push.rwmu.Unlock()

		rp := push.usersCallbackUrl[userKey]
		if rp != nil {
			return rp, nil
		}
		rp, err := db.ReadProfile(userKey)
		if err != nil {
			return nil, err
		}
		push.usersCallbackUrl[userKey] = rp
		return rp, nil
	}()
}

func (push *Push) reloadUserCallbackUrl(userKey string) (*backend.AckUserReadProfile, error) {
	push.rwmu.Lock()
	defer push.rwmu.Unlock()

	rp, err := db.ReadProfile(userKey)
	if err != nil {
		return nil, err
	}
	push.usersCallbackUrl[userKey] = rp
	return rp, nil
}

func (push *Push) GetApiGroup() map[string]service.NodeApi {
	nam := make(map[string]service.NodeApi)

	func() {
		service.RegisterApi(&nam,
			"pushdata", data.APILevel_client, push.PushData)
	}()

	return nam
}

func (push *Push) HandleNotify(req *data.SrvRequest) {
	if req.Method.Srv == "account" && req.Method.Function == "updateprofile" {
		//reqUpdateProfile := v1.ReqUserUpdateProfile{}
		//err := json.Unmarshal([]byte(req.Argv.Message), &reqUpdateProfile)
		//if err != nil {
		//	l4g.Error("HandleNotify-Unmarshal: %s", err.Error())
		//	return
		//}

		// reload profile
		rp, err := push.reloadUserCallbackUrl(req.Argv.SubUserKey)
		if err != nil {
			l4g.Error("HandleNotify-reloadUserCallbackUrl: %s", err.Error())
			return
		}

		l4g.Info("HandleNotify-reloadUserCallbackUrl: ", rp)
	}
}

// 推送数据
func (push *Push) PushData(req *data.SrvRequest, res *data.SrvResponse) {
	rp, err := push.getUserProfile(req.Argv.UserKey)
	if err != nil {
		l4g.Error("(%s) no user callback: %s", req.Argv.UserKey, err.Error())
		res.Err = apibackend.ErrPushSrvPushData
		return
	}

	l4g.Info("push %s to %s-%s", req.Argv.Message, req.Argv.UserKey, rp.CallbackUrl)

	func() {
		// encrypt
		timeStamp := ""
		if push.intervalSecond > 0 {
			timeStamp = strconv.FormatInt(time.Now().Unix(), 10)
		}

		srvData, err := data.EncryptionAndSignData([]byte(req.Argv.Message), timeStamp, req.Argv.UserKey, []byte(rp.PublicKey), push.privateKey)
		if err != nil {
			l4g.Error("EncryptionAndSignData: %s", err.Error())
			res.Err = apibackend.ErrInternal
			return
		}

		pushData := api.UserResponseData{}
		srvData.ToApiData(&pushData.Value)

		// call url
		b, err := json.Marshal(pushData)
		if err != nil {
			l4g.Error("error json message: %s", err.Error())
			res.Err = apibackend.ErrDataCorrupted
			return
		}

		l4g.Info("push data: %s", string(b))

		httpCode, ret, err := nethelper.CallToHttpServer(rp.CallbackUrl, "", string(b))
		if err != nil {
			l4g.Error("push http: %s", err.Error())
			res.Err = apibackend.ErrPushSrvPushData
			return
		}
		res.Value.Message = ret

		l4g.Info("push status:%d-%s", httpCode, ret)
	}()

	l4g.Info("push fin to %s", req.Argv.UserKey)
}
