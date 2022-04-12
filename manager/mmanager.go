package manager

import (
	"errors"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/module"
)

var ModManager *MManager

type MManager struct {
	BaseMod *BaseMod
	ModRoot *DefaultModRoot
}

func (m *MManager) Init() {
	m.BaseMod = new(BaseMod)
	m.BaseMod.Init()
	m.ModRoot = new(DefaultModRoot)
	m.ModRoot.Init()
}

func (m *MManager) SetModRoot(d *DefaultModRoot) {
	m.ModRoot = d
}

func (m *MManager) Run() {

}

func (m *MManager) Close() {
}

func (m *MManager) RegisterMod(modName []string, room *module.GORoom) error {
	if nil == m.ModRoot {
		err := errors.New("nil Module Root")
		log.Fatal("err:%v", err)
	}
	return m.ModRoot.Register(modName, room)
}

func (m *MManager) GetMod(name string) (*module.ModCluster, error) {
	if nil == m.ModRoot {
		err := errors.New("nil Module Root")
		log.Error("err:%v", err)
		return nil, err
	}
	return m.ModRoot.GetModCluster(name)
}

func GetNewRoom(mod module.Module) error {
	if nil == ModManager {
		log.Error("nil Module Root")
		return errors.New("nil Module Root")
	}
	cluster, err := ModManager.GetMod(mod.GetName())
	if nil != err {
		return err
	}
	return cluster.AddNewRoom(mod)
}

func CloseRoom(mod module.Module) error {
	if nil == ModManager {
		log.Error("nil Module Root")
		return errors.New("nil Module Root")
	}
	cluster, err := ModManager.GetMod(mod.GetName())
	if nil != err {
		return err
	}
	return cluster.CloseRoom(mod.GetKey())
}
