package data

import (
	"github.com/mikeqiao/newworld/newworld/db/redis"
)

type UpdateMod struct {
	table  string
	do     bool
	update map[string]interface{}
}

func (u *UpdateMod) Init() {
	u.update = make(map[string]interface{})
}

func (u *UpdateMod) Update() {
	if u.do {
		redis.R.Hash_SetDataMap(u.table, u.update)
		u.do = false
		u.update = make(map[string]interface{})
	}
}

func (u *UpdateMod) AddData(key string, value interface{}) {
	u.update[key] = value
	u.do = true
}
