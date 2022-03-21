// Author: mike.qiao
// File:base_mod
// Date:2022/3/21 14:49

package manager

import (
	"errors"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	"github.com/mikeqiao/newworld/net/base_proto"
	"time"
)

type BaseMod struct {
	ModeName string
}

func (b *BaseMod) Init() {
	b.ModeName = common.Mod_Base
}

func (b *BaseMod) GetModName() string {
	return b.ModeName
}
func (b *BaseMod) Route(u *net.CallData) error {
	if nil == u {
		return errors.New("nil Call net.UserData")
	}
	switch u.Function {
	case common.BaseMod_AgentTick:
		return BaseModAgentTick(u)
	case common.BaseMod_AgentCheck:
		return BaseModAgentCheck(u)
	}
	return nil
}

func BaseModAgentTick(call *net.CallData) error {
	m := new(base_proto.ServerTick)
	err := call.GetReqMsg(m)
	if nil != err {
		return err
	}
	// 消息的发送者
	if nil != call && nil != call.Agent {
		log.Debug("server:%v, tick:%v", call.Agent.GetRemoteAddr(), m.GetTime())
		call.Agent.SetTick(time.Now().Unix())
	}
	return nil
}

func BaseModAgentCheck(call *net.CallData) error {
	m := new(base_proto.ServerLogInCheckReq)
	err := call.GetReqMsg(m)
	if nil != err {
		return err
	}
	// 消息的发送者
	if nil != call && nil != call.Agent {
		log.Debug("server:%v, tick:%v", call.Agent.GetRemoteAddr(), m.Sid)
		call.Agent.Check(m.Sid)
		ConnManager.AddAgent(call.Agent)
	}
	return nil
}
