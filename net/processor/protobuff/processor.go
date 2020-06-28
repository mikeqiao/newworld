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
func (p *Processor) Unmarshal(data []byte) error {
	msg := new(bmsg.CallMsgInfo)
	err := proto.Unmarshal(data, msg)
	if nil != err {
		log.Error("Unmarshal NetMsg err:%v", err)
		return err
	} else {
		switch msg.MsgType {
		case common.Msg_Req:
			//找到callid 赋值调用
			cid := msg.CallID
			if v, ok := p.FuncList[cid]; ok && nil != v {
				cmsg := reflect.New(v.in.Elem()).Interface()
				err = proto.Unmarshal(msg.Info, cmsg.(proto.Message))
				if nil != err {
					if nil != v.server {
						v.server.Call()
					} else {
						log.Error("this service:%v not working", cid)
					}
				} else {
					log.Error("Unmarshal call msg err:%v", err)
				}
			} else {
				log.Error("no this service:%v", cid)
			}
		case common.Msg_Res:
			//找到callback 赋值调用

		case common.Msg_Push:
			//找到callid 赋值调用 不用解析 in
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
