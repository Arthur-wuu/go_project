package common

import (
	"github.com/ulule/limiter"
	ltredis "github.com/ulule/limiter/drivers/store/redis"
)

//业务上限制，非接口限制，

type BusLimiter struct{
	redis    *DbRedis
	rates    []*limiter.Rate
	lms      []*limiter.Limiter
	prefix   string
}

func NewBusLimiter(redis *DbRedis, prefix string, rates []*limiter.Rate) *BusLimiter {
	bus := &BusLimiter{
		redis: redis,
		rates : rates,
		prefix:prefix,
	}
	return bus
}

func (this *BusLimiter) Init() error {
	this.lms = make([]*limiter.Limiter, 0, len(this.rates))
	store, err := ltredis.NewStoreWithOptions(this.redis.GetConn(), limiter.StoreOptions{
		Prefix:   this.prefix,
		MaxRetry: 3,
	})
	if err != nil {
		return err
	}
	for i:=0; i< len(this.rates); i++ {
		lm := limiter.New(store, *this.rates[i])
		this.lms = append(this.lms, lm)
	}
	return nil
}

func (this *BusLimiter)  Check(key string) (bool,string, error){
	for i:=0; i < len(this.lms);i++ {
		//fmt.Println(len(this.lms), config.GConfig.BussinessLimits.PhoneSms[i], config.GConfig.BussinessLimits.IpSms[i])
		ctx, err := this.lms[i].Get(nil, this.rates[i].Formatted + key)
		if err != nil {
			return true, this.rates[i].Formatted, err
		}
		if ctx.Reached {
			return true,this.rates[i].Formatted,nil
		}
	}
	return false,"",nil
}