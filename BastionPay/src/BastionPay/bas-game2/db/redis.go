package db

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"time"
)

var GRedis DbRedis

type DbRedis struct {
	pool *redis.Pool
}

func (this *DbRedis) Init(host string, port string, password string, db string) error {
	this.pool = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   100,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
			if err != nil {
				return nil, err
			}
			// 使用密码
			if password != "" {
				_, err := c.Do("AUTH", password)
				if err != nil {
					return nil, err
				}
			}
			// 选择db
			if db != "" {
				_, err := c.Do("SELECT", db)
				if err != nil {
					return nil, err
				}
			}
			return c, nil
		},
	}

	_, err := this.Do("GET", "test")
	if err != nil {
		return err
	}
	return nil
}

func (r *DbRedis) Close() error {
	return r.pool.Close()
}

func (r *DbRedis) GetConn() redis.Conn {
	conn := r.pool.Get()
	return conn
}

func (r *DbRedis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.GetConn()
	if conn == nil {
		return nil, errors.New("redis pool nil conn")
	}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (this *DbRedis) Scan(args ...interface{}) (int, []string, error) {
	repy, err := redis.Values(this.Do("SCAN", args...))
	if err == redis.ErrNil {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}
	if repy == nil {
		return 0, nil, nil
	}
	var cur int
	var keys []string
	_, err = redis.Scan(repy, &cur, &keys)
	if err != nil {
		return 0, nil, err
	}
	return cur, keys, nil
}

func (this *DbRedis) Llen(key string) (int64, error) {
	repy, err := this.Do("llen", key)
	if err != nil {
		return 0, err
	}
	len, ok := repy.(int64)
	if !ok {
		return 0, errors.New("redis llen type err")
	}
	return len, nil
}

func (this *DbRedis) Setnx(k, v string) error {
	_, err := this.Do("setnx", k, v)
	return err
}

func (this *DbRedis) Setex(k, v string, expire int) error {
	_, err := this.Do("setex", k, expire, v)
	return err
}

func (this *DbRedis) Exist(k string) (bool, error) {
	repy, err := this.Do("EXISTS", k)
	if err != nil {
		return false, err
	}
	exist, ok := repy.(int64)
	if !ok {
		return false, errors.New("redis exist type err " + reflect.TypeOf(repy).String())
	}
	if exist > 0 {
		return true, nil
	}
	return false, err
}

func (this *DbRedis) Get(k string) ([]byte, error) {
	repy, err := this.Do("get", k)
	if err != nil {
		return nil, err
	}
	if repy == nil {
		return nil, nil
	}
	data, ok := repy.([]byte)
	if !ok {
		return nil, errors.New("redis get type err")
	}
	return data, err
}
