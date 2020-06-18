package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mikeqiao/newworld/config"
)

var R *CRedis

type CRedis struct {
	Pool *redis.Pool
	Life uint32
}

func init() {
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
				log.Panicf("err:%v", err)
				return nil, err
			}

			return c, nil
		},
	}
	log.Println("connnect redis success")
	return pool
}
