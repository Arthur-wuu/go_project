package src

import (
	"github.com/bluele/gcache"
	"BastionPay/bas-vipsys/config"
	"errors"
	"math/rand"
	"time"
)


var GCache Cache

type Cache struct{
	VipAuthCache gcache.Cache
	VipListCache gcache.Cache
	VipDisableCache gcache.Cache
}

func (this * Cache)Init(){
	GCache.VipDisableCache = gcache.New(config.GConfig.Cache.VipDisableMaxKey).LRU().Expiration(time.Duration(config.GConfig.Cache.VipDisableTimeout)*time.Second ).Build()
}

func (this * Cache) GetVipAuth(level int) (interface{}, error) {
	if this.VipAuthCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.VipAuthCache.Get(level)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	return value, nil
}

func (this * Cache) SetVipAuthFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.VipAuthCache = gcache.New(config.GConfig.Cache.VipAuthMaxKey).LRU().LoaderExpireFunc(f).Build()
}

func (this * Cache) GetVipList(key string) (interface{}, error) {
	if this.VipListCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.VipListCache.Get(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	return value, nil
}

func (this * Cache) SetVipListFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.VipListCache = gcache.New(config.GConfig.Cache.VipListMaxKey).LRU().LoaderExpireFunc(f).Build()
}

func (this * Cache) GetVipDisable(level int) (interface{}, error) {
	if this.VipDisableCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.VipDisableCache.Get(level)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	return value, nil
}

func (this * Cache) SetVipDisable(level, value int) {
	this.VipDisableCache.Set(level, value)
}



func (this *Activity) GetByUuidFromCache(uuid string) (*Activity, error) {
	data,err := db.GCache.GetActivity(uuid)
	if err == gorm.ErrRecordNotFound {
		return nil,nil
	}
	if err != nil {
		return nil,err
	}
	if data == nil {
		return nil, nil
	}

	acty, ok := data.(*Activity)
	if !ok {
		return nil, errors.Annotate(err, "type err")
	}
	return acty,nil
}



//在cache里面注册个方法
func (this * Cache) SetActivityFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.ActivityCache = gcache.New(config.GConfig.Cache.ActivityMaxKey).LRU().LoaderExpireFunc(f).Build()
}

//init
db.GCache.SetActivityFunc(new(models.Activity).InnerGetByUuid)

//实现
func (this *Activity) InnerGetByUuid(input interface{}) (interface{}, *time.Duration, error) {
	expire := time.Second * time.Duration(config.GConfig.Cache.ActivityTimeout+rand.Intn(600))
	userKey,ok := input.(string)
	if !ok {
		return nil,nil, errors.New("type err")
	}
	acty,err := new(Activity).GetByUuid(userKey, nil, []string{"sponsor_id", "valid", "off_at", "sponsor_account", "notify_url"})
	if err != nil {
		return nil, nil, errors.Annotate(err, "Activity GetByIdAndFields")
	}
	if acty == nil {
		return nil, nil, gorm.ErrRecordNotFound  // nil,nil,nil可能将是永远不超时
	}
	return  acty, &expire, nil
}