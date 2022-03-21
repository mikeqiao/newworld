package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mikeqiao/newworld/log"
)

func (r *CRedis) Hash_GetAllData(table string) (map[string]string, error) {
	c := r.Pool.Get()
	defer func() {
		err := c.Close()
		if nil != err {
			log.Error("Hash_GetAllData err :%v", err)
		}
	}()
	return redis.StringMap(c.Do("hgetall", table))
}

func (r *CRedis) Hash_SetDataMap(table string, data map[string]interface{}) error {
	c := r.Pool.Get()
	defer func() {
		err := c.Close()
		if nil != err {
			log.Error("Hash_GetAllData err :%v", err)
		}
	}()
	args := make([]interface{}, 1+len(data)*2)
	args[0] = table
	i := 1
	for k, v := range data {
		args[i] = k
		args[i+1] = v
		i += 2
	}
	_, err := c.Do("hmset", args...)
	return err
}

func (r *CRedis) Hash_DelDataMap(table string, data map[string]interface{}) error {
	c := r.Pool.Get()
	defer func() {
		err := c.Close()
		if nil != err {
			log.Error("Hash_GetAllData err :%v", err)
		}
	}()
	args := make([]interface{}, 1+len(data))
	args[0] = table
	i := 1
	for k := range data {
		args[i] = k
		i += 1
	}
	_, err := c.Do("hdel", args...)
	return err
}

func (r *CRedis) Hash_GetData(table string, key interface{}) (string, error) {
	c := r.Pool.Get()
	defer func() {
		err := c.Close()
		if nil != err {
			log.Error("Hash_GetAllData err :%v", err)
		}
	}()
	return redis.String(c.Do("hget", table, key))
}

func (r *CRedis) Hash_SetData(table string, name, value interface{}) error {
	c := r.Pool.Get()
	defer func() {
		err := c.Close()
		if nil != err {
			log.Error("Hash_GetAllData err :%v", err)
		}
	}()
	_, err := c.Do("hset", table, name, value)
	return err
}
