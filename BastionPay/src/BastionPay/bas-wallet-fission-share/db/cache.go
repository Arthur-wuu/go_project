package db

import (
	"BastionPay/bas-wallet-fission-share/config"
	"errors"
	"github.com/bluele/gcache"
	"time"
)

var GCache Cache

type Cache struct {
	LevelAuthCache gcache.Cache
	UserLevelCache gcache.Cache
	LevelRuleCache gcache.Cache
}

func (this *Cache) Init() {
}

//**************************************************/
func (this *Cache) GetLevelAuth(level int) (interface{}, error) {
	if this.LevelAuthCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.LevelAuthCache.Get(level)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *Cache) SetLevelAuthFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.LevelAuthCache = gcache.New(config.GConfig.Cache.LevelAuthMaxKey).LRU().LoaderExpireFunc(f).Build()
}

func (this *Cache) RemoveLevelAuth(level int) {
	if this.LevelAuthCache == nil {
		return
	}
	this.LevelAuthCache.Remove(level)
}

//**************************************************/
func (this *Cache) GetUserLevel(user string) (interface{}, error) {
	if this.LevelAuthCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.UserLevelCache.Get(user)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *Cache) RemoveUserLevel(user string) {
	if this.LevelAuthCache == nil {
		return
	}
	this.UserLevelCache.Remove(user)
}

func (this *Cache) SetUserLevelFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.UserLevelCache = gcache.New(config.GConfig.Cache.UserLevelMaxKey).LRU().LoaderExpireFunc(f).Build()
}

//****************************************************/

func (this *Cache) GetLevelRule(lvId int) (interface{}, error) {
	if this.LevelRuleCache == nil {
		return nil, errors.New("not init")
	}
	value, err := this.LevelRuleCache.Get(lvId)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (this *Cache) RemoveLevelRule(lvId int) {
	if this.LevelRuleCache == nil {
		return
	}
	this.LevelRuleCache.Remove(lvId)
}

func (this *Cache) SetLevelRuleFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.LevelRuleCache = gcache.New(config.GConfig.Cache.LevelRuleMaxKey).LRU().LoaderExpireFunc(f).Build()
}

//****************************************************/
