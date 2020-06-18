package data

import (
	"sync"

	"github.com/mikeqiao/newworld/db/redis"
)

type UpdateMod struct {
	table  string
	do     bool
	mutex  sync.RWMutex
	update map[string]interface{}
}

func (u *UpdateMod) Init(table string) {
	u.table = table
	u.update = make(map[string]interface{})
}

func (u *UpdateMod) Update() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if u.do {
		redis.R.Hash_SetDataMap(u.table, u.update)
		u.do = false
		u.update = make(map[string]interface{})
	}
}

func (u *UpdateMod) AddData(key string, value interface{}) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.update[key] = value
	u.do = true
}

func (u *UpdateMod) GetAllData() map[string]string {
	data, _ := redis.R.Hash_GetAllData(u.table)
	return data
}
