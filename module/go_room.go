// Author: mike.qiao
// File:go_room
// Date:2022/3/14 17:34

package module

import (
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	"runtime"
	"sync"
)

type GORoom struct {
	RoomId       uint64
	ModName      string
	RoomState    common.ModState //room 的状态
	Mod          Module
	CallChan     chan *net.CallData //请求的队列
	roomCloseSig chan bool          //room 模块关闭信号
}

func (r *GORoom) Init(modName string, roomId uint64, callLen uint32, mod Module) {
	r.ModName = modName
	r.RoomId = roomId
	r.CallChan = make(chan *net.CallData, callLen)
	r.roomCloseSig = make(chan bool, 1)
	r.Mod = mod
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

func (r *GORoom) ExecFunc(call *net.CallData) {
	if nil != r.Mod {
		h, err := r.Mod.GetHandler(call.Function)
		if nil != err || nil == h {
			log.Error("GetHandler err:%v", err)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				if config.Conf.LenStackBuf > 0 {
					buf := make([]byte, int32(config.Conf.LenStackBuf))
					l := runtime.Stack(buf, false)
					log.Error("%v: %s", r, buf[:l])
				} else {
					log.Error("%v", r)
				}
			}
		}()
		err = h(call)
		if nil != err {
			log.Error("run handler err:%v", err)
			return
		}
	}
}

func (r *GORoom) Working() bool {
	if common.Mod_Running == r.RoomState {
		return true
	}
	return false
}

func (c *GORoom) Call(data *net.CallData) {
	if nil == data {
		return
	}
	c.CallChan <- data
}

func (c *GORoom) Close() {
	c.roomCloseSig <- true
}
