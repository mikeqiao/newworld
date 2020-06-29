package processor

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/log"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
	bmsg "github.com/mikeqiao/newworld/net/proto"
)

type MsgHandler func(msg interface{}, data *net.UserData)

type Processor struct {
	FuncList map[string]*FuncInfo
}

type FuncInfo struct {
	fid    string
	in     reflect.Type    //请求数据类型
	out    reflect.Type    //返回数据类型
	f      interface{}     //服务
	server *mod.RpcService //服务的模块
}

//解包数据
func (p *Processor) Unmarshal(a *net.TcpAgent, data []byte) error {
	msg := new(bmsg.CallMsgInfo)
	err := proto.Unmarshal(data, msg)
	if nil != err {
		log.Error("Unmarshal NetMsg err:%v", err)
		return err
	} else {
		switch msg.MsgType {
		case common.Msg_Req:
			//找到callid 赋值调用
			callId := msg.CallID
			if v, ok := p.FuncList[callId]; ok && nil != v {
				cmsg := reflect.New(v.in.Elem()).Interface()
				err = proto.Unmarshal(msg.Info, cmsg.(proto.Message))
				if nil != err {
					if nil != v.server {
						udata := new(net.UserData)
						udata.UId = msg.UId
						udata.UIdList = msg.UIdList[:]
						udata.CallId = msg.CallID
						udata.CallBackId = msg.CallBackID
						udata.MsgType = msg.MsgType
						udata.Agent = a
						if nil != a && 2 == a.Ctype {
							udata.UId = a.RUId
						}
						cb := func(in interface{}, e error) {
							udata.MsgType = common.Msg_Res
							a.WriteMsg(udata, in)
						}

						v.server.Call(v.f, cb, cmsg, udata)
					} else {
						log.Error("this service:%v not working", callId)
					}
				} else {
					log.Error("Unmarshal call msg err:%v", err)
				}
			} else {
				log.Error("no this service:%v", callId)
			}
		case common.Msg_Res:
			//找到callid 赋值调用
			callId := msg.CallID
			if v, ok := p.FuncList[callId]; ok && nil != v {
				cmsg := reflect.New(v.out.Elem()).Interface()
				err = proto.Unmarshal(msg.Info, cmsg.(proto.Message))
				if nil != err {
					if nil != v.server {
						cbid := msg.CallBackID
						v.server.CallBack(cbid, cmsg, err)
					} else {
						log.Error("this service:%v not working", callId)
					}
				} else {
					log.Error("Unmarshal call msg err:%v", err)
				}
			} else {
				log.Error("no this service:%v", callId)
			}
		case common.Msg_Push:
			//找到callid 赋值调用 不用解析 in
			callId := msg.CallID
			if v, ok := p.FuncList[callId]; ok && nil != v {
				if nil != v.server {
					udata := new(net.UserData)
					udata.UId = msg.UId
					udata.UIdList = msg.UIdList[:]
					udata.CallId = msg.CallID
					udata.CallBackId = msg.CallBackID
					udata.MsgType = msg.MsgType
					udata.Agent = a
					v.server.Call(v.f, nil, msg.Info, udata)
				} else {
					log.Error("this service:%v not working", callId)
				}
			}
		default:
			log.Error("err msgType:%v", msg.MsgType)
		}
	}
	return nil
}

//打包数据
func (p *Processor) Marshal(u *net.UserData, in interface{}) (*net.UserData, [][]byte, error) {
	msg := new(bmsg.CallMsgInfo)
	msg.MsgType = u.MsgType
	return nil, nil, nil
}

func (p *Processor) Register(tmod *mod.Mod) {
	if nil == tmod {
		return
	}
	flist := tmod.GetAllFunc()
	if nil != flist {
		for k, v := range flist {
			if nil != v {
				f := new(FuncInfo)
				f.fid = k
				f.f = v.F
				f.in = v.In
				f.in = v.Out
				f.server = tmod.Server
			}
		}
	}
}
