package db

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"errors"
)

var GRedis Redis

type Redis struct {
	pool *redis.Pool
}

func (this *Redis) Init(host string, port string, password string, db string) error {
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

	_,err := this.Do("GET", "test")
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) GetConn() redis.Conn {
	conn := r.pool.Get()
	return conn
}

func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.GetConn()
	if conn == nil {
		return nil, errors.New("redis pool nil conn")
	}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (this *Redis) Scan(args ...interface{}) (int, []string, error) {
	repy, err := redis.Values(this.Do("SCAN", args...))
	if err != redis.ErrNil {
		return 0, nil, nil
	}
	if repy == nil {
		return 0,nil, nil
	}
	var cur int
	var keys []string
	_, err = redis.Scan(repy, &cur, &keys)
	if err != nil {
		return 0,nil,err
	}
	return cur , keys, nil
}

