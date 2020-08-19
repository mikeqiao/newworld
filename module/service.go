package module

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	conf "github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	"github.com/mikeqiao/newworld/net/proto"
)

type RpcService struct {
	UId          uint64
	State        uint64 //服务者状态(服务数量)
	Max          uint32 //服务 最多等待数量
	Working      bool
	ChanCall     chan *CallInfo
	ChanCallBack chan *Return //
	WCallBack    map[string]*CallInfo
	NFunc        map[string]uint32 //用于统计
	ErrorList    []string          //用于统计
	mutex        sync.RWMutex
}

func (s *RpcService) Init() {
	s.ChanCall = make(chan *CallInfo, s.Max)
	s.ChanCallBack = make(chan *Return, s.Max)
	s.WCallBack = make(map[string]*CallInfo)
	s.NFunc = make(map[string]uint32)
}

func (s *RpcService) Start() {
	s.Working = true
}

func (s *RpcService) Call(name string, f, cb, in interface{}, out reflect.Type, udata *net.UserData) {
	if !s.Working {
		return
	}
	ci := &CallInfo{
		FuncId:  name,
		Out:     out,
		CF:      f,
		Args:    in,
		Data:    udata,
		chanRet: s.ChanCallBack,
		Cb:      cb,
	}
	var err error
	select {
	case s.ChanCall <- ci:
	default:
		err = errors.New("call channel full")
	}
	if err != nil && nil != cb {
		log.Error("err:%v", err)
		s.ChanCallBack <- &Return{err: err, ret: nil, cb: cb}
	} else {
		if v, ok := s.NFunc[name]; ok {
			v += 1
		} else {
			s.NFunc[name] = 1
		}

	}
}

func (s *RpcService) Exec(ci *CallInfo) {
	if nil == ci {
		log.Error("nil call")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			s.ErrorList = append(s.ErrorList, ci.FuncId)
			if conf.Conf.LenStackBuf > 0 {
				buf := make([]byte, int32(conf.Conf.LenStackBuf))
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
			s.ret(ci, &Return{err: fmt.Errorf("%v", r)})
		}
	}()
	f, ok := ci.CF.(func(*CallInfo))
	if ok {
		f(ci)
	} else if nil != ci.Cb {
		s.ret(ci, &Return{err: fmt.Errorf("err func format")})
		log.Error("err func format")
	}
	if v, ok := s.NFunc[ci.FuncId]; ok {
		v -= 1
	} else {
		log.Error("err func count")
	}
}

func (s *RpcService) CallBack(cbid string, in interface{}, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	cb, ok := s.WCallBack[cbid]
	if ok {
		delete(s.WCallBack, cbid)
	}

	cb.SetResult(in, err)

}

func (s *RpcService) GetCallBack(cbid string) *CallInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	cb, ok := s.WCallBack[cbid]
	if ok {
		delete(s.WCallBack, cbid)
		return cb
	}

	return nil

}

func (s *RpcService) ExecCallBack(ri *Return) {
	defer func() {
		if r := recover(); r != nil {
			if conf.Conf.LenStackBuf > 0 {
				buf := make([]byte, conf.Conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()
	if nil == ri.cb {
		return
	}
	f, ok := ri.cb.(func(interface{}, error))
	if ok {
		f(ri.ret, ri.err)
	} else {
		log.Error("err cb format")
	}
	return
}

func (s *RpcService) Close() {
	s.Working = false
	err := errors.New("Module colsed")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.WCallBack {
		if nil != v {
			r := &Return{
				err: err,
				ret: nil,
				cb:  v.Cb,
			}
			s.ExecCallBack(r)
		}
	}

}

func (s *RpcService) ret(ci *CallInfo, ri *Return) {
	if s.ChanCallBack == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Debug("%v", err)
		}
	}()

	ri.cb = ci.Cb
	s.ChanCallBack <- ri
}

func (s *RpcService) AddWaitCallBack(c *CallInfo) string {
	key := fmt.Sprintf("%p", c)
	s.mutex.Lock()
	s.WCallBack[key] = c
	s.mutex.Unlock()
	return key
}

func (s *RpcService) GetServiceInfo(info *proto.ModInfo) {
	info.Call = uint32(len(s.ChanCall))
	info.CallBack = uint32(len(s.ChanCallBack))
	info.WCallBack = uint32(len(s.WCallBack))
	var t []*proto.FuncState = make([]*proto.FuncState, len(s.NFunc))
	for k, v := range s.NFunc {
		t = append(t, &proto.FuncState{
			Name:  k,
			State: v,
		})
	}
	info.Info = t[:]
	info.Err = s.ErrorList[:]
}

//msg->route->mod->find func->callinfo ->add to server chan
// run-> get callinfo-> do func
//func 需要注册到 消息解析器中
