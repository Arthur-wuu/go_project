package db

import (
	"github.com/bluele/gcache"
)

var GCache Cache

type Cache struct {
	Dingtalk gcache.Cache
}

func (this *Cache) SetDingtalkCache() {
	this.Dingtalk = gcache.New(3).Simple().Build()
}
