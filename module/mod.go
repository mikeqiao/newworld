package module

import (
	"reflect"
	"sync"

	"github.com/mikeqiao/newworld/log"
)

type Mod struct {
	Uid      uint64
	Name     string
	Server   *RpcService
	FuncList map[string]*SFunc
	closeSig chan bool //模块关闭信号
	Closed   bool
	Working  bool
}

type SFunc struct {
	Ftyp    uint32       // 1 本地服务  2 远程服务
	InName  string       //请求数据类型名字
	OutName string       //返回数据类型名字
	In      reflect.Type //请求数据类型
	Out     reflect.Type //返回数据类型
	F       interface{}  //服务
}

func (m *Mod) Init() {
	m.Server = new(RpcService)
	m.Server.Init()
	m.FuncList = make(map[string]*SFunc)

}

func (m *Mod) Start(wg *sync.WaitGroup) {
	go m.Run(wg)
}

func (m *Mod) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	m.Server.Start()
	for {
		select {
		case <-m.closeSig:
			m.Server.Close()
			goto Loop
		case ri := <-m.Server.ChanCallBack:
			m.Server.ExecCallBack(ri)
		case ci := <-m.Server.ChanCall:
			m.Server.Exec(ci)
		}
	}
Loop:
	m.Closed = true
	//m.working = false
	wg.Done()
}

func (m *Mod) Close() {
	m.closeSig <- true
}

func (m *Mod) GetAllFunc() map[string]*SFunc {
	return m.FuncList

}

func (m *Mod) Register(fname string, f, req, res interface{}) {
	if _, ok := m.FuncList[fname]; ok {
		log.Error("func already registed, name:%v", fname)
		return
	}

	if nil != f {
		sf := new(SFunc)
		sf.F = f
		sf.In = reflect.TypeOf(req)
		sf.Out = reflect.TypeOf(res)
		sf.InName = sf.In.Elem().Name()
		sf.OutName = sf.Out.Elem().Name()
		sf.Ftyp = 1
		m.FuncList[fname] = sf
	}
}

func (m *Mod) RegisterRemote(fname, req, res string, f interface{}) {
	if _, ok := m.FuncList[fname]; ok {
		log.Error("func already registed, name:%v", fname)
		return
	}

	if nil != f {
		sf := new(SFunc)
		sf.F = f
		sf.InName = req
		sf.OutName = res
		sf.Ftyp = 2
		m.FuncList[fname] = sf
	}
}
