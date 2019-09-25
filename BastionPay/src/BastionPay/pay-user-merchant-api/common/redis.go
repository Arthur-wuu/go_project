package common


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
		Addr:     host+":"+port,
		Password: password,
		DB:       db,  // use default DB
	})

	_,err := this.client.Ping().Result()
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

//import (
//	"github.com/garyburd/redigo/redis"
//	"time"
//)
//
//type Redis struct {
//	pool *redis.Pool
//}
//
//func NewRedis(host string, port string, password string, db string, useTLS bool) *Redis {
//	r := Redis{}
//
//	r.pool = &redis.Pool{
//		MaxIdle:     1,
//		MaxActive:   100,
//		IdleTimeout: 180 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			dailOpt := redis.DialUseTLS(useTLS)
//			c, err := redis.Dial("tcp", host+":"+port, dailOpt)
//			if err != nil {
//				return nil, err
//			}
//			// 使用密码
//			if password != "" {
//				_, err := c.Do("AUTH", password)
//				if err != nil {
//					return nil, err
//				}
//			}
//			// 选择db
//			if db != "" {
//				_, err := c.Do("SELECT", db)
//				if err != nil {
//					return nil, err
//				}
//			}
//			return c, nil
//		},
//	}
//
//	return &r
//}
//
//func (r *Redis) GetConn() redis.Conn {
//	conn := r.pool.Get()
//	return conn
//}
//
//func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
//	conn := r.GetConn()
//	defer conn.Close()
//	return conn.Do(commandName, args...)
//}
