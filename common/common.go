package common

const (
	Msg_Type = uint32(iota)
	Msg_Req
	Msg_Res
	Msg_Push
)

type ModState int32

const (
	Mod_state    ModState = iota
	Mod_Init              //初始化成功
	Mod_Running           //正在运行工作
	Mod_Stopping          //正在停止
	Mod_Stopped           //已经停止工作
	Mod_Closed            //关闭销毁
)
