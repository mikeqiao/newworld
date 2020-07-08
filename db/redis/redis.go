package redis

import (
	"fmt"

	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/log"
)

var R *CRedis

type CRedis struct {
	Pool *redis.Pool
	Life uint32
}

func Init() {
	R = new(CRedis)
	R.InitDB()
}

func (r *CRedis) InitDB() {
	r.Pool = Newfactory("")
	r.Life = config.Conf.Redis.Life * 3600 * 24
}

func (r *CRedis) OnClose() {
	r.Pool.Close()
}

func Newfactory(name string) *redis.Pool {

	host := config.Conf.Redis.Host
	port := config.Conf.Redis.Port
	password := config.Conf.Redis.Password
	count := config.Conf.Redis.Count
	pool := &redis.Pool{
		IdleTimeout: 180 * time.Second,
		MaxIdle:     int(count),
		MaxActive:   1024,
		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%s", host, port)
			c, err := redis.Dial("tcp", address,
				redis.DialPassword(password),
			)
			if err != nil {
				log.Fatal("err:%v, pw:%v, addr:%v", err, password, address)
				return nil, err
			}

			return c, nil
		},
	}
	log.Debug("connnect redis success")
	return pool
}
