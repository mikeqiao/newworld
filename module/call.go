package module

import (
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
)

type CallInfo struct {
	ModId   uint64 //module uid
	FuncId  uint32
	CF      interface{}   //执行function
	Cb      interface{}   //callback
	Args    interface{}   //参数
	Data    *net.UserData //附加信息
	chanRet chan *Return  //
}

type Return struct {
	ret interface{}
	cb  interface{}
}

func (c *CallInfo) SetResult(res interface{}) {

	if c.chanRet == nil || nil == c.Cb {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Debug("%v", err)
		}
	}()
	log.Debug("callback 5")
	r := &Return{
		ret: res,
		cb:  c.Cb,
	}
	c.chanRet <- r
}
