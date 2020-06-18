package module

type Mod struct {
	Uid      uint64
	Name     string
	server   *RpcService
	FuncList map[uint32]interface{}
	Call     map[uint64]interface{}
	CallBack map[uint64]interface{}
}

func (m *Mod) Init() {
	m.server = new(RpcService)
	m.server.Init()
	m.FuncList = make(map[uint32]interface{})
}

func (m *Mod) Start() {

}

func (m *Mod) Run() {

}

func (m *Mod) Close() {

}

func (m *Mod) Route() {

}
