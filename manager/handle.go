package manager

import (
	"time"

	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/log"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
	"github.com/mikeqiao/newworld/net/proto"
)

func HandleServerConnectOK(msg interface{}, data *net.UserData) {
	//m := msg.(*proto.ServerConnect)
	if nil != data && nil != data.Agent {
		data.MsgType = common.Msg_Handle
		data.CallId = "ServerLoginRQ"
		reqMsg := new(proto.ServerLogInReq)
		reqMsg.Sid = config.Conf.SInfo.Uid
		reqMsg.Name = config.Conf.SInfo.Name
		fdata := data.Agent.Processor.GetLocalFunc()
		reqMsg.Flist = fdata[:]
		data.Agent.WriteMsg(data, reqMsg)
	}
}

func HandleDelConnect(msg interface{}, data *net.UserData) {

}

func HandleServerDelConnect(msg interface{}, data *net.UserData) {

}

func HandleServerTick(msg interface{}, data *net.UserData) {
	//	m := msg.(*proto.ServerTick)
	// 消息的发送者
	if nil != data && nil != data.Agent {
		data.Agent.SetTick(time.Now().Unix())
		u := new(net.UserData)
		u.MsgType = common.Msg_Handle
		u.CallId = "ServerTickBack"
		nowtime := time.Now().Unix()
		data.Agent.WriteMsg(u, &proto.ServerTick{
			Time: uint32(nowtime),
		})
	}
}

func HandleServerTickBack(msg interface{}, data *net.UserData) {
	//	m := msg.(*proto.ServerTick)
	// 消息的发送者
	if nil != data && nil != data.Agent {
		data.Agent.SetTick(time.Now().Unix())
	}
}

func HandleServerLoginRQ(msg interface{}, data *net.UserData) {

	m := msg.(*proto.ServerLogInReq)
	// 消息的发送者
	if nil == data || nil == data.Agent {
		log.Debug(" wrong agent ,msg: %v", msg)
		return
	}
	s := m.GetFlist()
	if 0 == len(s) {
		log.Debug(" wrong server func info msg: %+v", m)
		return
	}
	uid := m.GetSid()
	name := m.GetName()
	module := ModManager.GetMod(uid)
	if nil != module {
		ModManager.RemoveMod(uid)
	}
	agent := data.Agent
	agent.SetRemotUID(uid)
	agent.SetLogin()
	module = NewMod(uid, name)
	for _, v := range m.GetFlist() {
		if nil != v {
			f := func(c *mod.CallInfo) {
				if nil == c || nil == module {
					return
				}
				key := module.Server.AddWaitCallBack(c)
				if nil != agent {

					u := new(net.UserData)
					u.CallId = v.Name
					u.CallBackId = key
					u.MsgType = common.Msg_Req
					if nil != c.Data {
						if 0 != c.Data.MsgType {
							u.MsgType = c.Data.MsgType
						}
						u.UId = c.Data.UId
						u.UIdList = c.Data.UIdList
					}
					agent.WriteMsg(u, c.Args)

				} else if "" != key {
					log.Error("err agent")
				}
			}

			module.RegisterRemote(v.Name, v.In, v.Out, f)
		}
	}
	ModManager.Registe(module)

	data.CallId = "ServerLogInRes"
	data.MsgType = common.Msg_Handle
	reqMsg := new(proto.ServerLogInRes)
	reqMsg.Sid = config.Conf.SInfo.Uid
	reqMsg.Name = config.Conf.SInfo.Name
	fdata := data.Agent.Processor.GetLocalFunc()
	reqMsg.Flist = fdata[:]
	data.Agent.WriteMsg(data, reqMsg)
}

func HandleServerLoginRS(msg interface{}, data *net.UserData) {
	m := msg.(*proto.ServerLogInRes)
	// 消息的发送者
	if nil == data || nil == data.Agent {
		log.Debug(" wrong agent ,msg: %v", msg)
		return
	}

	uid := m.GetSid()
	name := m.GetName()
	module := ModManager.GetMod(uid)
	if nil != module {
		ModManager.RemoveMod(uid)
	}
	agent := data.Agent
	agent.SetRemotUID(uid)
	agent.SetLogin()
	module = NewMod(uid, name)
	for _, v := range m.GetFlist() {
		if nil != v {
			f := func(c *mod.CallInfo) {
				if nil == c || nil == module {
					return
				}
				key := module.Server.AddWaitCallBack(c)
				if nil != agent {

					u := new(net.UserData)
					u.CallId = v.Name
					u.CallBackId = key
					u.MsgType = common.Msg_Req
					if nil != c.Data {
						if 0 != c.Data.MsgType {
							u.MsgType = c.Data.MsgType
						}
						u.UId = c.Data.UId
						u.UIdList = c.Data.UIdList
					}
					agent.WriteMsg(u, c.Args)

				} else if "" != key {
					log.Error("err agent")
				}
			}

			module.RegisterRemote(v.Name, v.In, v.Out, f)
		}
	}
	ModManager.Registe(module)
}
