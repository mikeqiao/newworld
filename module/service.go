package module

type RpcService struct {
	UId      uint64
	State    uint64 //服务者状态(服务数量)
	Max      uint32 //服务 最多等待数量
	ChanCall chan *CallInfo
}

func (s *RpcService) Init() {
	s.ChanCall = make(chan *CallInfo, s.Max)
}

func (s *RpcService) Call() {

}

//msg->route->mod->find func->callinfo ->add to server chan
// run-> get callinfo-> do func
//func 需要注册到 消息解析器中
