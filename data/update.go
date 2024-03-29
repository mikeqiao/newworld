package data

import (
	"github.com/mikeqiao/newworld/log"
	"sync"

	"github.com/mikeqiao/newworld/db/redis"
)

type UpdateMod struct {
	table  string
	do     bool
	mutex  sync.RWMutex
	del    map[string]interface{}
	update map[string]interface{}
}

func (u *UpdateMod) Init(table string) {
	u.table = table
	u.del = make(map[string]interface{})
	u.update = make(map[string]interface{})
}

func (u *UpdateMod) Update() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if u.do {
		if len(u.update) > 0 {
			err := redis.R.Hash_SetDataMap(u.table, u.update)
			if nil != err {
				u.update = make(map[string]interface{})
			} else {
				log.Error("table:%v, err:%v", u.table, err)
			}
		}
		if len(u.del) > 0 {
			err := redis.R.Hash_DelDataMap(u.table, u.del)
			if nil != err {
				u.do = false
				u.del = make(map[string]interface{})
			} else {
				log.Error("table:%v, err:%v", u.table, err)
			}
		}
	}
}

func (u *UpdateMod) AddData(key string, value interface{}) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.update[key] = value
	u.do = true
	if _, ok := u.del[key]; ok {
		delete(u.del, key)
	}
}

func (u *UpdateMod) DelData(key string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.del[key] = 1
	u.do = true
	if _, ok := u.update[key]; ok {
		delete(u.update, key)
	}
}

func (u *UpdateMod) GetAllData() map[string]string {
	data, err := redis.R.Hash_GetAllData(u.table)
	if nil != err {
		log.Error("table:%v, err:%v", u.table, err)
		return nil
	}
	return data
}
