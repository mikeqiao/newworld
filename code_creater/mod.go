// Author: mike.qiao
// File:mod
// Date:2022/3/15 18:05

package code_creater

import (
	"errors"
	"github.com/mikeqiao/newworld/net"
)

type Mod_Root struct {
	DefaultModName ModName
}

func (m *Mod_Root) Route(data *net.CallData) (err error) {
	switch data.Mod {
	case "ModTestOne":
		err = m.DefaultModName.Route_Room(data)
	default:
		err = errors.New("nil function")
	}
	return
}

//自动生成
type ModName interface {
	Route_Room(*net.CallData) (err error)
}
