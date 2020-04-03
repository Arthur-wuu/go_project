package db

import (
	"BastionPay/merchant-api/config"
	"github.com/bluele/gcache"
	"time"
)

var GCache Cache

type Cache struct {
	AccountCache gcache.Cache
	//VipListCache gcache.Cache
	//VipDisableCache gcache.Cache
}

func (this *Cache) Init() {
	GCache.AccountCache = gcache.New(config.GConfig.Cache.VipAuthMaxKey).LRU().Expiration(time.Duration(config.GConfig.Cache.VipAuthTimeout) * time.Second).Build()
}

func (this *Cache) GetAccountCache(name string) (interface{}, error) {
	//if this.AccountCache == nil {
	//	return nil, errors.New("not init")
	//}
	value, err := this.AccountCache.Get(name)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *Cache) SetAccountCache(account, statu string) {
	this.AccountCache.SetWithExpire(account, statu, time.Second*time.Duration(config.GConfig.Cache.VipAuthMaxKey))
}

//func (this * Cache) GetVipList(key string) (interface{}, error) {
//	if this.VipListCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.VipListCache.Get(key)
//	if err != nil {
//		return nil, err
//	}
//	if value == nil {
//		return nil, nil
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetVipListFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.VipListCache = gcache.New(config.GConfig.Cache.VipListMaxKey).LRU().LoaderExpireFunc(f).Build()
//}
//
//func (this * Cache) RemoveVipList(key interface{}) {
//	this.VipListCache.Remove(key)
//}
//
//func (this * Cache) GetVipDisable(level int) (interface{}, error) {
//	if this.VipDisableCache == nil {
//		return nil, errors.New("not init")
//	}
//	value, err := this.VipDisableCache.Get(level)
//	if err != nil {
//		return nil, err
//	}
//	if value == nil {
//		return nil, nil
//	}
//
//	return value, nil
//}
//
//func (this * Cache) SetVipDisable(level, value int) {
//	this.VipDisableCache.Set(level, value)
//}
