// Author: mike.qiao
// File:mod_name
// Date:2022/3/15 18:08

package code_creater

import (
	"errors"
	"fmt"
	"github.com/mikeqiao/newworld/code_creater/msg_proto"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/net"
	"sync"
)

var DefaultModName *ModNameRoomManager

type ModNameRoomManager struct {
	name     string
	lock     sync.RWMutex
	roomList map[uint64]*ModNameRoom
}

func (m *ModNameRoomManager) GetRoom(key uint64) *ModNameRoom {
	m.lock.RLock()
	defer m.lock.Unlock()
	v, ok := m.roomList[key]
	if ok && v.Working() {
		return v
	}
	return nil
}
func (m *ModNameRoomManager) Route(data *net.CallData) (err error) {
	room := m.GetRoom(data.Uid)
	if nil == room {
		return errors.New(fmt.Sprintf("no this mod room:%v", data.Uid))
	}
	return room.Call(data)
}

type ModNameRoom_Type interface {
	NameHandler_One(*msg_proto.FuncReq) *msg_proto.FuncRes
	NameHandler_Two(*msg_proto.FuncReq) *msg_proto.FuncRes
	ModName_Handler(*net.CallData) (err error)
}
type ModNameRoom struct {
	RoomId       uint64
	ModName      string
	RoomState    common.ModState    //room 的状态
	CallChan     chan *net.CallData //请求的队列
	roomCloseSig chan bool          //room 模块关闭信号
	Data         *msg_proto.FuncReq
}

func (m *ModNameRoom) ModName_Handler(data *net.CallData) (err error) {
	if nil == data {
		return errors.New("nil CallData")
	}
	switch data.Function {
	case "FuncTestOne":
		msg := new(msg_proto.FuncReq)
		err := data.GetReqMsg(msg)
		if nil != err {
			return err
		}
		res := m.NameHandler_One(data)
		err = data.CallBack(res)
		return err
	default:
		err = errors.New("nil function")
	}

	return
}

func (r *ModNameRoom) Working() bool {
	if common.Mod_Running == r.RoomState {
		return true
	}
	return false
}

func (c *ModNameRoom) Call(data *net.CallData) error {
	if nil == data {
		return errors.New("nil CallData")
	}
	if len(c.CallChan) >= cap(c.CallChan) {
		return errors.New("full CallChan")
	}
	c.CallChan <- data
	return nil
}

//手动添加
func (m *ModNameRoom) NameHandler_One(data *net.CallData) (res *msg_proto.FuncRes) {

	return
}
func (m *ModNameRoom) NameHandler_Two(data *net.CallData) (res *msg_proto.FuncRes) {
	return
}
