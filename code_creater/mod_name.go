// Author: mike.qiao
// File:mod_name
// Date:2022/3/15 18:08

package code_creater

import (
	"errors"
	"github.com/mikeqiao/newworld/code_creater/msg_proto"
	"github.com/mikeqiao/newworld/data"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
	processor "github.com/mikeqiao/newworld/net/processor/protobuff"
	"sync"
)

type ModName struct {
	name     string
	lock     sync.RWMutex
	roomList map[uint64]*mod.GORoom
}

func (m *ModName) GetRoom(key uint64) *mod.GORoom {
	m.lock.RLock()
	defer m.lock.Unlock()
	v, ok := m.roomList[key]
	if ok && v.Working() {
		return v
	}
	return nil
}

func (m *ModName) ModName_Handler(funcName string, data net.UserData, req []byte) (err error) {
	switch funcName {
	case "FuncTestOne":
		msg := new(msg_proto.FuncReq)
		err := processor.UnMarshalMsg(msg, req)
		if nil != err {
			return err
		}
		fn := NameHandler
	default:
		err = errors.New("nil function")
	}

	room := m.GetRoom(data.UId)
	if nil == room {
		return
	}
	room.Call(fn)
	return
}

func NameHandler(uData net.UserData, msg *msg_proto.FuncReq, Data data.MemoryData) (err error) {

	return
}
