package common

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Redis struct {
	pool *redis.Pool
}

func NewRedis(host string, port string, password string, db string) *Redis {
	r := Redis{}

	r.pool = &redis.Pool{
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

	return &r
}

func (r *Redis) GetConn() redis.Conn {
	conn := r.pool.Get()
	return conn
}

func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.GetConn()
	defer conn.Close()
	return conn.Do(commandName, args...)
}
