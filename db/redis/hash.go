package redis

import (
	"github.com/gomodule/redigo/redis"
)

func (r *CRedis) Hash_GetAllData(table string) (map[string]string, error) {
	c := r.Pool.Get()
	value, err := redis.StringMap(c.Do("hgetall", table))
	if nil != err {
		//	log.Error("table:%v, error:%v", table, err)
	}

	c.Close()
	return value, err
}

func (r *CRedis) Hash_SetDataMap(table string, data map[string]interface{}) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data)*2)
	args[0] = table
	i := 1
	for k, v := range data {
		args[i] = k
		args[i+1] = v
		i += 2
	}
	_, err := c.Do("hmset", args...)
	if nil != err {
		//	log.Error("error table:%v, data:%v", table, data)
	}
	c.Close()
	return err
}
