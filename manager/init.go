package manager

import (
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
	"github.com/mikeqiao/newworld/net/proto"
)

func Init() {
	DefaultProcessor = new(p.Processor)
	DefaultProcessor.Init()
	ConnManager = new(NetConnManager)
	ConnManager.Init()
	ModManager = new(MManager)
	ModManager.Init()
	Register()
}

func Register() {
	DefaultProcessor.SetHandler("ServerConnectOK", proto.ServerConnect{}, HandleServerConnectOK)
	DefaultProcessor.SetHandler("ServerLoginRQ", proto.ServerLogInReq{}, HandleServerLoginRQ)
	DefaultProcessor.SetHandler("ServerLoginRS", proto.ServerLogInRes{}, HandleServerLoginRS)
}
