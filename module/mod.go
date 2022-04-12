package module

import (
	"errors"
	"fmt"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	"sync"
)

type Handler func(uData *net.CallData) (err error)

type Module interface {
	GetHandler(string) (Handler, error)
	GetName() string
	GetKey() uint64
	GetCallLen() uint32
	Init()
}

type ModCluster struct {
	lock    sync.RWMutex
	ModName string
	ModRoom map[uint64]*GORoom
	wg      *sync.WaitGroup
}

func (m *ModCluster) Init(name string) {
	m.ModName = name
	m.wg = new(sync.WaitGroup)
	m.ModRoom = make(map[uint64]*GORoom)
}

//添加一个房间协程
func (m *ModCluster) AddNewRoom(mod Module) error {
	if nil == mod {
		return errors.New("nil mod Data")
	}
	key := mod.GetKey()
	m.lock.Lock()
	r, ok := m.ModRoom[key]
	if ok {
		log.Debug("room :%v already existed", key)
	} else {
		r = new(GORoom)
		r.Init(mod.GetName(), key, mod.GetCallLen(), mod)
		m.ModRoom[key] = r
	}
	m.lock.Unlock()
	if nil != r && !r.Working() {
		r.Start(m.wg)
	}
	return nil
}
func (m *ModCluster) AddRoom(room *GORoom) error {
	if nil == room {
		return errors.New("nil room Data")
	}
	key := room.RoomId
	m.lock.Lock()
	_, ok := m.ModRoom[key]
	if ok {
		log.Debug("room :%v already existed", key)
	} else {
		m.ModRoom[key] = room
	}
	m.lock.Unlock()
	if nil != room && !room.Working() {
		room.Start(m.wg)
	}
	return nil
}

//直接关闭房间并且删除
func (m *ModCluster) DeleteRoom(key uint64) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	r, ok := m.ModRoom[key]
	if !ok {
		return errors.New(fmt.Sprintf("no this key:%v room", key))
	}
	if nil != r {
		r.Close()
	}
	delete(m.ModRoom, key)
	return nil
}

//只关闭房间，数据保留一段时间
func (m *ModCluster) CloseRoom(key uint64) error {
	m.lock.RLock()
	r, ok := m.ModRoom[key]
	m.lock.RUnlock()
	if !ok {
		return errors.New(fmt.Sprintf("no this key:%v room", key))
	}
	if nil != r {
		r.Close()
	}
	return nil
}

func (m *ModCluster) Close() {
	m.lock.Lock()
	for _, v := range m.ModRoom {
		if nil != v {
			v.Close()
		}
	}
	m.lock.Unlock()
	m.wg.Wait()
}
func (m *ModCluster) Route(u *net.CallData) error {
	if 0 == u.RoomId {
		u.RoomId = m.GetRecommendRoom()
	}
	m.lock.RLock()
	r, ok := m.ModRoom[u.RoomId]
	m.lock.RUnlock()
	if !ok {
		return errors.New(fmt.Sprintf("no this key:%v room", u.Uid))
	}
	if nil != r {
		r.Call(u)
	}
	return nil
}

func (m *ModCluster) GetRecommendRoom() uint64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.ModRoom {
		if nil != v {
			return k
		}
	}
	return 0
}
