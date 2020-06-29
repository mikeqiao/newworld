package net

type UserData struct {
	MsgType    uint32    //消息类型
	UId        uint64    //用户id
	UIdList    []uint64  //目标群id
	Agent      *TcpAgent //网络链接
	CallId     string    //调用的id
	CallBackId string    //回调id
}

type Processor interface {

	//解包数据
	Unmarshal(a *TcpAgent, data []byte) error
	//打包数据
	Marshal(u *UserData, msg interface{}) (*UserData, [][]byte, error)
}
