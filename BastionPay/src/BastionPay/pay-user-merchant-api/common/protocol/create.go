package protocol

import (
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"time"
)

func InitConnPool(min, max int, factory func() (interface{}, interface{}, error), close func(v interface{}) error) (Pool, error) {
	poolConfig := &PoolConfig{
		InitialCap: min,
		MaxCap:     max,
		Factory:    factory,
		Close:      close,
		//链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err := NewChannelPool(poolConfig)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("NewChannelPool err")
		//		glog.Error(err.Error())
		return nil, err
	}
	return p, nil
}
