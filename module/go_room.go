// Author: mike.qiao
// File:go_room
// Date:2022/3/14 17:34

package module

import (
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/data"
	"github.com/mikeqiao/newworld/log"
	"sync"
)

type GORoom struct {
	RoomId       uint64
	ModName      string
	RoomState    common.ModState //room 的状态
	BindData     data.MemoryData //绑定的相关数据模块
	CallChan     chan *CallInfo  //请求的队列
	roomCloseSig chan bool       //room 模块关闭信号
}

func (r *GORoom) Init(modName string, roomId uint64, callLen uint32, bindData data.MemoryData) {
	r.ModName = modName
	r.RoomId = r.RoomId
	r.CallChan = make(chan *CallInfo, callLen)
	r.roomCloseSig = make(chan bool, 1)
	r.BindData = bindData
	r.RoomState = common.Mod_Init
}

func (r *GORoom) Start(wg *sync.WaitGroup) {
	go r.Run(wg)
}

func (r *GORoom) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	log.Debug("mod:%v, Start", r.ModName)
	r.RoomState = common.Mod_Running
	for {
		select {
		case <-r.roomCloseSig:
			log.Debug("mod:%v, room:%v, close", r.ModName, r.RoomId)
			r.RoomState = common.Mod_Stopping
			close(r.CallChan)
		case ri, ok := <-r.CallChan:
			if !ok {
				goto Loop
			}
			r.ExecFunc(ri)
		}
	}
Loop:
	r.RoomState = common.Mod_Closed
	//m.working = false
	wg.Done()
}

func (r *GORoom) ExecFunc(call *CallInfo) {
	//call.CF()
}

func (r *GORoom) Working() bool {
	if common.Mod_Running == r.RoomState {
		return true
	}
	return false
}

func (r *GORoom) Call(function Handler) {

}
