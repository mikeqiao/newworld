package manager

import (
	"errors"
	"sync"

	"github.com/mikeqiao/newworld/net/proto"

	"github.com/mikeqiao/newworld/log"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
)

var ModManager *MManager

type MManager struct {
	createid uint64
	modList  map[uint64]*mod.Mod
	mutex    sync.RWMutex
	wg       *sync.WaitGroup
}

func NewMod(uid uint64, name string) *mod.Mod {
	mid := uid
	if 0 == mid {
		mid = ModManager.GetNewID()
	}
	m := new(mod.Mod)
	m.Uid = mid
	m.Name = name
	m.Init()
	return m
}

func (m *MManager) Init() {
	m.wg = new(sync.WaitGroup)
	m.modList = make(map[uint64]*mod.Mod)
	m.createid = 0
}

func (m *MManager) AddGroup(count int) {
	m.wg.Add(count)
}

func (m *MManager) Done() {
	m.wg.Done()
}

func (m *MManager) GetNewID() uint64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.createid += 1

	return m.createid
}

func (m *MManager) Registe(mod *mod.Mod) error {

	if nil == mod {
		log.Error("This mod is nil")
		return errors.New("This mod is nil")
	}
	m.mutex.Lock()
	m.modList[mod.Uid] = mod
	m.mutex.Unlock()
	if nil != DefaultProcessor {
		DefaultProcessor.Register(mod)
		log.Debug("register mod:%v", mod.Name)
	}
	return nil
}

func (m *MManager) Route(fid string, cb interface{}, in interface{}, data *net.UserData) {
	DefaultProcessor.Route(fid, cb, in, data)
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
			v.Start(m.wg)
		} else {
			log.Error("this mod  is nil:%v", k)
		}
	}
}

func (m *MManager) StartMod(mod *mod.Mod) {

	if nil != mod {
		mod.Start(m.wg)
	} else {
		log.Error("this mod  is nil")
	}
}

func (m *MManager) Close() {
	m.mutex.Lock()

	for k, v := range m.modList {
		if nil != v {
			v.Close()
		} else {
			log.Error("this mod  is nil:%v", k)
		}
	}
	m.mutex.Unlock()
	m.wg.Wait()
}

func (m *MManager) GetAllModState() []*proto.ModInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var info []*proto.ModInfo = make([]*proto.ModInfo, len(m.modList))
	for _, v := range m.modList {
		if nil != v {
			a := v.GetModState()
			info = append(info, a)
		}
	}
	return info

}
