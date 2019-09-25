package common

import (
	. "github.com/BastionPay/bas-base/log/zap"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

const (
	LimterPrefix = "limiter_"
)

type PathLimit struct {
	Path   string
	Method string
	Limit  int
	Time   int
}

type Limiter struct {
	redis  *Redis
	pls    []*PathLimit
	path   string
	method string
	ip     string
}

func NewLimiter(redis *Redis, pls []*PathLimit, ctx iris.Context) *Limiter {
	return &Limiter{redis, pls, ctx.Path(), ctx.Method(), ctx.RemoteAddr()}
}

func DoLimiter(redis *Redis, pls []*PathLimit, ctx iris.Context) (bool, error) {
	return NewLimiter(redis, pls, ctx).Verification()
}

func (l *Limiter) Verification() (bool, error) {
	var (
		bol = true
	)

	for _, pl := range l.pls {
		if strings.ToUpper(pl.Method) == l.method && pl.Path == l.path {
			count, err := l.countDown(pl)
			if err != nil {
				//				glog.Error(err.Error())
				ZapLog().With(zap.NamedError("err", err)).Error("countDown err")
				return false, err
			}

			if count > pl.Limit {
				bol = false
			}
		}
	}

	return bol, nil
}

func (l *Limiter) countDown(pl *PathLimit) (int, error) {
	var (
		err error
		key string
	)

	key = LimterPrefix + l.ip + "_" + l.method + l.path + "_" + strconv.Itoa(pl.Time)

	result, err := l.redis.Do("INCR", key)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("redis INCR err")
		//		glog.Error(err.Error())
		return 0, err
	}
	if result == nil {
		return 0, nil
	}

	res := int(result.(int64))
	if res == 1 {
		_, err = l.redis.Do("SET", key, res, "EX", pl.Time)
		if err != nil {
			return 0, err
		}
	}

	return res, err
}
