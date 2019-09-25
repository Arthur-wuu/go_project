package db

import (
	//"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/go-redis/redis"
)

var GRedis DbRedis

type DbRedis struct {
	client *redis.Client
}

func (this *DbRedis) Init(host string, port string, password string, db int) error {
	this.client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db, // use default DB
	})

	_, err := this.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *DbRedis) Close() error {
	return r.client.Close()
}

func (r *DbRedis) GetConn() *redis.Client {
	conn := r.client
	return conn
}

func (r *DbRedis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if r.client == nil {
		return nil, errors.New("redis pool nil conn")
	}
	return r.client.Do(commandName, args).Result()
}
