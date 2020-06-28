package manager

import (
	"errors"
	"sync"

	"github.com/mikeqiao/newworld/log"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
)

var ModManager *MManager

type MManager struct {
	createid uint64
	modList  map[uint64]*mod.Mod
	mutex    sync.RWMutex
}

func NewMod(name string) *mod.Mod {
	mid := ModManager.GetNewID()
	m := new(mod.Mod)
	m.Uid = mid
	m.Name = name
	m.Init()
	return m
}

func (m *MManager) Init() {
	m.modList = make(map[uint64]*mod.Mod)
	m.createid = 0
}

func (m *MManager) GetNewID() uint64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.createid += 1

	return m.createid
}

func (m *MManager) Registe(mod *mod.Mod) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if nil == m {
		return errors.New("This mod is nil")
	}

	m.modList[mod.Uid] = mod
	if nil != DefaultProcessor {
		DefaultProcessor.Register(mod)
	}
	return nil
}

func (m *MManager) Route(mid uint64, fid uint32, cb interface{}, in interface{}, data *net.UserData) {

}

func (m *MManager) RemoveMod(mid uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tmod, ok := m.modList[mid]
	if ok && nil != tmod {
		tmod.Close()
	}
	delete(m.modList, mid)
}

func (m *MManager) GetMod(mid uint64) *mod.Mod {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	tmod, ok := m.modList[mid]
	if ok && nil != tmod {
		return tmod
	}
	return nil
}

func (m *MManager) Run() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.modList {
		if nil != v {
			v.Start()
		} else {
			log.Error("this mod  is nil:%v", k)
		}
	}
}

func (m *MManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.modList {
		if nil != v {
			v.Close()
		} else {
			log.Error("this mod  is nil:%v", k)
		}
	}
}
