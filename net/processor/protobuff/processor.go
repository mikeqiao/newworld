package processor

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/log"
	mod "github.com/mikeqiao/newworld/module"
	"github.com/mikeqiao/newworld/net"
	bmsg "github.com/mikeqiao/newworld/net/proto"
)

type MsgHandler func(msg interface{}, data *net.UserData)

type Processor struct {
	FuncList   map[string]*FuncInfo
	MsgList    map[string]reflect.Type
	HandleList map[string]*HandleInfo
	mutex      sync.RWMutex
}

type HandleInfo struct {
	handle MsgHandler
	in     reflect.Type //请求数据类型
}

type FuncInfo struct {
	ftype   uint32          //1 本地服务 2 远程注册来的服务
	fid     string          //服务名称
	InName  string          //请求数据类型名字
	OutName string          //返回数据类型名字
	in      reflect.Type    //请求数据类型
	out     reflect.Type    //返回数据类型
	f       interface{}     //服务
	server  *mod.RpcService //服务的模块
}

func (p *Processor) Init() {
	p.FuncList = make(map[string]*FuncInfo)
	p.MsgList = make(map[string]reflect.Type)
	p.HandleList = make(map[string]*HandleInfo)
}

func (p *Processor) SetHandler(fid string, in interface{}, msgHandler MsgHandler) {
	if _, ok := p.HandleList[fid]; ok {
		log.Error("function %s already registered", fid)
	} else {
		info := new(HandleInfo)
		info.handle = msgHandler
		info.in = reflect.TypeOf(in)
		p.HandleList[fid] = info
	}
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
			p.mutex.RLock()
			defer p.mutex.RUnlock()
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

						v.server.Call(v.f, cb, cmsg, v.out, udata)
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
			v, ok := a.Mod.(*mod.RpcService)
			if ok && nil != v {
				cbid := msg.CallBackID
				cb := v.GetCallBack(cbid)
				if nil == cb {
					log.Error("no this callinfo callBackId:%v", cbid)
					return nil
				}
				cmsg := reflect.New(cb.Out.Elem()).Interface()
				err = proto.Unmarshal(msg.Info, cmsg.(proto.Message))
				if nil == err {
					cb.SetResult(cmsg, err)
				} else {
					log.Error("Unmarshal call msg err:%v", err)
				}
			} else {
				log.Error("no this service:%v", a)
			}
		case common.Msg_Push:
			//找到callid 赋值调用 不用解析 in
			callId := msg.CallID
			p.mutex.RLock()
			defer p.mutex.RUnlock()
			if v, ok := p.FuncList[callId]; ok && nil != v {
				if nil != v.server {
					udata := new(net.UserData)
					udata.UId = msg.UId
					udata.UIdList = msg.UIdList[:]
					udata.CallId = msg.CallID
					udata.CallBackId = msg.CallBackID
					udata.MsgType = msg.MsgType
					udata.Agent = a
					v.server.Call(v.f, nil, msg.Info, v.out, udata)
				} else {
					log.Error("this service:%v not working", callId)
				}
			}
		case common.Msg_Handle:
			callId := msg.CallID
			if v, ok := p.HandleList[callId]; ok && nil != v {
				cmsg := reflect.New(v.in.Elem()).Interface()
				if nil != v.handle {
					udata := new(net.UserData)
					udata.UId = msg.UId
					udata.UIdList = msg.UIdList[:]
					udata.CallId = msg.CallID
					udata.CallBackId = msg.CallBackID
					udata.MsgType = msg.MsgType
					udata.Agent = a
					v.handle(cmsg, udata)
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
	msg.CallID = u.CallId
	msg.CallBackID = u.CallBackId
	msg.UId = u.UId
	msg.UIdList = u.UIdList[:]
	data, err := proto.Marshal(in.(proto.Message))
	if nil != err {
		log.Error("err:%v", err)
		return nil, nil, err
	}
	msg.Info = data[:]
	udata, err := proto.Marshal(msg)
	if nil == err {
		return u, [][]byte{udata}, err
	}
	return nil, nil, err
}

func (p *Processor) Register(tmod *mod.Mod) {
	if nil == tmod {
		return
	}
	flist := tmod.GetAllFunc()
	if nil != flist {
		p.mutex.Lock()
		for k, v := range flist {
			if nil != v {
				f := new(FuncInfo)
				f.fid = k
				f.f = v.F
				f.in = v.In
				f.out = v.Out
				f.InName = v.InName
				f.OutName = v.OutName
				f.server = tmod.Server
				f.ftype = v.Ftyp
				if 2 == f.ftype {
					if ttype, ok := p.MsgList[f.InName]; ok {
						f.in = ttype
					} else {
						log.Error("no this name msg:%v", f.InName)
					}
					if ttype, ok := p.MsgList[f.OutName]; ok {
						f.out = ttype
					} else {
						log.Error("no this name msg:%v", f.OutName)
					}
				}

				p.FuncList[f.fid] = f

			}
		}
		p.mutex.Unlock()
	}
}

func (p *Processor) RegisterMsg(name string, mtype reflect.Type) {
	p.MsgList[name] = mtype
}

func (p *Processor) Route(funcName string, cb, in interface{}, udata *net.UserData) {
	if v, ok := p.FuncList[funcName]; ok && nil != v {

		v.server.Call(v.f, cb, in, v.out, udata)
	}
}

func (p *Processor) Handle(funcName string, in interface{}, udata *net.UserData) {
	if v, ok := p.HandleList[funcName]; ok && nil != v {
		v.handle(in, udata)
	}
}

func (p *Processor) GetLocalFunc() (flist []*bmsg.FuncInfo) {
	for _, v := range p.FuncList {
		log.Debug("info:%v, type:%v", v.fid, v.ftype)
		if 1 == v.ftype {
			nf := new(bmsg.FuncInfo)
			nf.Name = v.fid
			nf.In = v.InName
			nf.Out = v.OutName
			flist = append(flist, nf)
		}
	}
	return
}
