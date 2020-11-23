package module

import (
	"reflect"
	"time"

	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
)

type CallInfo struct {
	FuncId  string
	Out     reflect.Type  //返回数据类型
	CF      interface{}   //执行function
	Cb      interface{}   //callback
	Args    interface{}   //参数
	Data    *net.UserData //附加信息
	chanRet chan *Return  //

	STime int64 //开始时间
	ETime int64 //完成时间
}

type CallFuncData struct {
	Uid    uint64
	STime  int64 //开始时间
	ETime  int64 //完成时间
	FuncID string
	Req    interface{}
	Res    interface{}
	Err    error
}

type Return struct {
	err error
	ret interface{}
	cb  interface{}
}

func (c *CallInfo) SetResult(res interface{}, err error) {

	c.ETime = time.Now().Unix()
	//处理 数据统计
	cbdata := new(CallFuncData)
	cbdata.Err = err
	cbdata.Uid = c.Data.UId
	cbdata.STime = c.STime
	cbdata.ETime = c.ETime
	cbdata.Req = c.Args
	cbdata.Res = res
	go func() {
		c.Data.Agent.Processor.Route("Statistics", nil, cbdata, c.Data)
	}()

	if c.chanRet == nil || nil == c.Cb {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Debug("%v", err)
		}
	}()
	r := &Return{
		err: err,
		ret: res,
		cb:  c.Cb,
	}
	c.chanRet <- r

}
