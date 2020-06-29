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
}

type SFunc struct {
	In  reflect.Type //请求数据类型
	Out reflect.Type //返回数据类型
	F   interface{}  //服务
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

	//m.working = false
	wg.Done()
}

func (m *Mod) Close() {
	m.closeSig <- true
}

func (m *Mod) Route() {

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
		m.FuncList[fname] = sf
	}
}
