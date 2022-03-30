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
	m.BaseMod.Init()
}

func (m *MManager) SetModRoot(d *DefaultModRoot) {
	m.ModRoot = d
}

func (m *MManager) Run() {

}

func (m *MManager) Close() {
}

func (m *MManager) RegisterMod(mod ...module.Module) error {
	if nil == m.ModRoot {
		err := errors.New("nil Module Root")
		log.Fatal("err:%v", err)
	}
	return m.ModRoot.Register(mod...)
}
