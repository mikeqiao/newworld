package module

import (
	"reflect"
)

type Mod struct {
	Uid      uint64
	Name     string
	Server   *RpcService
	FuncList map[uint32]*SFunc
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
	m.FuncList = make(map[uint32]*SFunc)
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

func (m *Mod) GetAllFunc() map[uint32]*SFunc {
	return m.FuncList

}
