package db

import (
	"github.com/bluele/gcache"
	//"errors"
	"time"
)

var GCache Cache

type Cache struct {
	QuoteCache gcache.Cache
	HuilvCache gcache.Cache
	//KxianIdCache gcache.Cache
	//KxianUsdCache gcache.Cache

}

//行情
func (this *Cache) SetQuoteCacheFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.QuoteCache = gcache.New(500).LRU().LoaderExpireFunc(f).Build()
}

////k线 use id
//func (this * Cache) SetKxianIdCacheFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.KxianIdCache = gcache.New(500).LRU().LoaderExpireFunc(f).Build()
//}
//
////k线 use usd
//func (this * Cache) SetKxianUsdCacheFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
//	this.KxianUsdCache = gcache.New(500).LRU().LoaderExpireFunc(f).Build()
//}

//汇率
func (this *Cache) SetHuilvCacheFunc(f func(key interface{}) (interface{}, *time.Duration, error)) {
	this.HuilvCache = gcache.New(500).LRU().LoaderExpireFunc(f).Build()
}
