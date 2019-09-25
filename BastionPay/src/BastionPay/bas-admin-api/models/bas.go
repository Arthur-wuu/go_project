package models

import (
	"encoding/json"
	"fmt"
	"github.com/BastionPay/bas-admin-api/bastionpay"
	"github.com/BastionPay/bas-admin-api/config"
	"github.com/BastionPay/bas-api/apibackend/v1/backend"
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/bluele/gcache"
	"github.com/bugsnag/bugsnag-go/errors"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const (
	CONST_Cache_Audite_Expire = 10 * 60 //可以配置替换
)

var GlobalBasModel BasModel

type BasModel struct {
	mAuditeCache gcache.Cache
	mConf        *config.Config

	mRunFlag bool
}

func (this *BasModel) Init(c *config.Config) {
	if this.mRunFlag {
		return
	}
	this.mConf = c
	this.mRunFlag = true
	if this.mConf.Cache.AuditeMaxKeyNum > 0 {
		//先不用回调，暂时不清楚回调是否阻塞其余接口
		if this.mConf.Cache.AuditeTimeout < 1 {
			this.mConf.Cache.AuditeTimeout = CONST_Cache_Audite_Expire
		}
		this.mAuditeCache = gcache.New(this.mConf.Cache.AuditeMaxKeyNum).LRU().Build()
		ZapLog().With(zap.Int("maxKey", this.mConf.Cache.AuditeMaxKeyNum), zap.Int("timeout", this.mConf.Cache.AuditeTimeout)).
			Info(" with AuditeCache")
	}
	ZapLog().Info("Init BasModel ok")
}

//必须加超时缓存的，10分钟都完全够了，缓存超时时间加上随机数避免雪崩
func (this *BasModel) GetUserAuditeStatus(userKey string) (uint, error) {
	if !this.mRunFlag {
		return 0, errors.New("not init", 0)
	}
	if this.mAuditeCache != nil {
		status, err := this.getAuditeFromCache(userKey)
		if err == nil {
			return status, nil
		}
		ZapLog().With(zap.Error(err), zap.String("userKey", userKey)).Info("AuditeCache Get err")
	}
	auditeStatus, err := this.loadAuditeStatus(userKey)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.String("userKey", userKey)).Error("loadAuditeStatus err")
		return 0, err
	}
	if this.mAuditeCache != nil {
		if err := this.setAuditeFromCache(userKey, auditeStatus); err != nil {
			ZapLog().With(zap.Error(err), zap.String("userKey", userKey)).Error("setAuditeFromCache err")
		}
	}
	return auditeStatus, nil
}

//无缓存
func (this *BasModel) GetUserAllAccountStatus(userKey string) (*backend.ResUserAccountStatus, error) {
	if !this.mRunFlag {
		return nil, errors.New("not init", 0)
	}
	auditeStatus, err := this.loadUserAccountStatus(userKey)
	if err != nil {
		ZapLog().With(zap.Error(err), zap.String("userKey", userKey)).Error("loadAuditeStatus err")
		return nil, err
	}
	return auditeStatus, nil
}

func (this *BasModel) loadUserAccountStatus(userKey string) (*backend.ResUserAccountStatus, error) {
	path := "/v1/account/getuserstatus"
	res, msgBytes, err := bastionpay.CallApi(userKey, "", path)
	if err != nil {
		return nil, err
	}
	if res.Err != 0 {
		return nil, fmt.Errorf("CallApi fail code %d %s", res.Err, res.ErrMsg)
	}
	info := new(backend.ResUserAccountStatus)
	if err := json.Unmarshal(msgBytes, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (this *BasModel) loadAuditeStatus(userKey string) (uint, error) {
	path := "/v1/account/getaudite"
	res, msgBytes, err := bastionpay.CallApi(userKey, "", path)
	if err != nil {
		return 0, err
	}
	if res.Err != 0 {
		return 0, fmt.Errorf("CallApi fail code %d %s", res.Err, res.ErrMsg)
	}
	auditeInfo := new(backend.ResUserAuditeStatus)
	if err := json.Unmarshal(msgBytes, auditeInfo); err != nil {
		return 0, err
	}
	return auditeInfo.AuditeStatus, nil
}

func (this *BasModel) getAuditeFromCache(userKey string) (uint, error) {
	if this.mAuditeCache == nil {
		return 0, errors.New("No Init Cache", 0)
	}
	value, err := this.mAuditeCache.Get(userKey)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, errors.New("cache with nil value", 0)
	}

	auditeStatus, ok := value.(uint)
	if !ok {
		return 0, errors.New("Cache Get Wrong Type", 0)
	}
	return auditeStatus, nil
}

func (this *BasModel) setAuditeFromCache(k, v interface{}) error {
	if this.mAuditeCache == nil {
		return errors.New("No Init Cache", 0)
	}
	timeout := this.mConf.Cache.AuditeTimeout
	expire := timeout + this.randInt(1, timeout/3)
	return this.mAuditeCache.SetWithExpire(k, v, time.Second*time.Duration(expire))
}

func (this *BasModel) randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
