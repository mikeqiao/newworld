package processor

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	"github.com/mikeqiao/newworld/net/base_proto"
)

type Processor struct {
	//map[消息id]消息处理方法
	ModRoot net.ModuleRoot
	BaseMod net.ModuleRoot
}

func (p *Processor) Init(module net.ModuleRoot) {
	p.BaseMod = module
}

//解包数据
func (p *Processor) Unmarshal(a *net.TcpAgent, data []byte) error {
	msg := new(base_proto.CallMsgInfo)
	err := proto.Unmarshal(data, msg)
	if nil != err {
		log.Error("Unmarshal NetMsg err:%v", err)
		return err
	} else {
		uData := new(net.CallData)
		uData.Agent = a
		uData.Uid = msg.UserId
		uData.Mod = msg.Mod
		uData.Function = msg.Function
		uData.Req = msg.MsgInfo
		uData.CallType = msg.CallType
		uData.Context = msg.Context
		return p.Route(uData)
	}
}

func (p *Processor) UnMarshalMsg(msg proto.Message, data []byte) error {
	if nil == msg {
		return errors.New("nil message")
	}
	if nil == data {
		return errors.New("nil data")
	}
	return proto.Unmarshal(data, msg)
}

//打包数据
func (p *Processor) Marshal(u *net.CallData, in interface{}) ([]byte, error) {
	msg := new(base_proto.CallMsgInfo)
	msg.Mod = u.Mod
	msg.Function = u.Function
	msg.UserId = u.Uid
	msg.Context = u.Context
	msg.CallType = u.CallType
	if nil != in {
		data, err := proto.Marshal(in.(proto.Message))
		if nil != err {
			log.Error("err:%v", err)
			return nil, err
		}
		msg.MsgInfo = data[:]
	} else {
		msg.MsgInfo = u.Req
	}
	msgData, err := proto.Marshal(msg)
	if nil == err {
		return msgData, err
	}
	return nil, err
}

func (p *Processor) Register(module net.ModuleRoot) {
	p.ModRoot = module
}

func (p *Processor) Route(uData *net.CallData) error {
	switch uData.Mod {
	case common.Mod_Base:
		if nil == p.BaseMod {
			return errors.New("nil modRoot")
		} else {
			return p.BaseMod.Route(uData)
		}
	default:
		if nil == p.ModRoot {
			return errors.New("nil modRoot")
		}
		return p.ModRoot.Route(uData)
	}

}
