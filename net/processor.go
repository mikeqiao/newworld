package net

type UserData struct {
	ModId   uint64
	FuncId  uint32
	MsgType uint32
	MsgId   uint32
	UId     uint64    //用户id
	UIdList []uint64  //目标群id
	Agent   *TcpAgent //网络链接
	Cid     string    //回调id
	Route   string    //客户端调用的route
}

type Processor interface {

	//解包数据
	Unmarshal(data []byte) error
	//打包数据
	Marshal(u *UserData, msg interface{}) (*UserData, [][]byte, error)
}
