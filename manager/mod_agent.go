// Author: mike.qiao
// File:mod_agent
// Date:2022/4/11 14:13

package manager

import (
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
)

type AgentMod struct {
	Name    string
	Uid     uint64
	CallLen uint32
	agent   *net.TcpAgent
}

func (a *AgentMod) Init() {
	a.Name = common.Mod_Agent
	a.CallLen = 1024
}
func (a *AgentMod) SetData(agent *net.TcpAgent, id uint64) {
	a.Uid = id
	a.agent = agent
}
func (a *AgentMod) GetCallLen() uint32 {
	return a.CallLen
}

func (a *AgentMod) GetName() string {
	return a.Name
}

func (a *AgentMod) GetKey() uint64 {
	return a.Uid
}

func (a *AgentMod) GetHandler(key string) (module.Handler, error) {
	return a.RequestFunc, nil
}

func (a *AgentMod) RequestFunc(uData *net.CallData) (err error) {
	if nil != a.agent {
		a.agent.SendMsg(uData, nil)
	}
	return nil
}
