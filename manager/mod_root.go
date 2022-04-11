// Author: mike.qiao
// File:mod_root
// Date:2022/3/24 17:43

package manager

import (
	"errors"
	"fmt"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
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

func (b *DefaultModRoot) Register(modName []string, room *module.GORoom) error {
	for _, v := range modName {
		if "" == v {
			return errors.New("nil module")
		}
		if _, ok := b.mList[v]; ok {
			return errors.New("same key module already registered")
		}
		modCluster := new(module.ModCluster)
		modCluster.Init(v)
		b.mList[v] = modCluster
		if nil != room {
			err := modCluster.AddRoom(room)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func (b *DefaultModRoot) GetModCluster(modName string) (*module.ModCluster, error) {
	if v, ok := b.mList[modName]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("no this modCluster named:%v", modName))
}

func (b *DefaultModRoot) Route(u *net.CallData) error {
	if nil == u {
		return errors.New("nil Call net.UserData")
	}
	if v, ok := b.mList[u.Mod]; ok {
		return errors.New("same key module already registered")
	} else {
		return v.Route(u)
	}
}

func (b *DefaultModRoot) RegisterAgentMod(modName []string) error {
	for _, v := range modName {
		if "" == v {
			return errors.New("nil module")
		}
		if _, ok := b.mList[v]; ok {
			return errors.New("same key module already registered")
		}
		modCluster := new(module.ModCluster)
		modCluster.Init(v)
		b.mList[v] = modCluster
	}
	return nil
}
