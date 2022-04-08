// Author: mike.qiao
// File:mod_root
// Date:2022/3/24 17:43

package manager

import (
	"errors"
	"fmt"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/module"
)

type DefaultModRoot struct {
	ModeName string
	//mod 模块应该是服务器启动时就注册的，暂时不考虑热更时候动态修改
	//暂时 是不需要加锁是管理控制的
	mList map[string]*module.ModCluster
}

func (b *DefaultModRoot) Init() {
	b.ModeName = common.ModRoot_Default
	b.mList = make(map[string]*module.ModCluster)
}

func (b *DefaultModRoot) GetName() string {
	return b.ModeName
}

func (b *DefaultModRoot) Register(mod ...module.Module) error {
	for _, v := range mod {
		if nil == v {
			return errors.New("nil module")
		}
		if _, ok := b.mList[v.GetName()]; ok {
			return errors.New("same key module already registered")
		}
		modCluster := new(module.ModCluster)
		modCluster.Init(v.GetName())
		b.mList[v.GetName()] = modCluster
	}
	return nil
}
func (b *DefaultModRoot) GetModCluster(modName string) (*module.ModCluster, error) {
	if v, ok := b.mList[modName]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("no this modCluster named:%v", modName))
}
