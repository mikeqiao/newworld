package module

import (
	"reflect"

	"github.com/mikeqiao/newworld/log"
)

type Mod struct {
	Uid      uint64
	Name     string
	Server   *RpcService
	FuncList map[string]*SFunc
	CallBack map[string]interface{}
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
	m.CallBack = make(map[string]interface{})
}

func (m *Mod) Start() {

}

func (m *Mod) Run() {

}

func (m *Mod) Close() {

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
