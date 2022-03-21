package net

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

type CallData struct {
	Mod      string    //请求的 mod 模块 key
	Function string    //请求的 function 模块 key
	Uid      uint64    //请求者id
	Agent    *TcpAgent //网络链接
	Req      []byte    //请求的信息
	Context  []byte    //传递的上下文
}

func (u *CallData) GetReqMsg(message proto.Message) error {
	if nil == u.Agent || nil == u.Agent.Processor {
		return errors.New("nil Processor")
	}
	return u.Agent.Processor.UnMarshalMsg(message, u.Req)
}
func (u *CallData) GetReqContext(message proto.Message) error {
	if nil == u.Agent || nil == u.Agent.Processor {
		return errors.New("nil Processor")
	}
	return u.Agent.Processor.UnMarshalMsg(message, u.Context)
}

type Processor interface {
	//解包数据
	Unmarshal(a *TcpAgent, data []byte) error
	UnMarshalMsg(msg proto.Message, data []byte) error
	//打包数据
	Marshal(u *CallData, msg interface{}) ([]byte, error)
	Route(u *CallData) error
	Register(module ModuleRoot)
}

type ModuleRoot interface {
	Route(u *CallData) error
}
