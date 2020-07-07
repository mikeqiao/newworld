package manager

import (
	"sync"

	"github.com/mikeqiao/newworld/db/redis"

	p "github.com/mikeqiao/newworld/net/processor/protobuff"
	"github.com/mikeqiao/newworld/net/proto"
)

var wg *sync.WaitGroup
var ServerManager *NetServerManager
var ClientManager *NetClientManager

func Init() {
	wg = new(sync.WaitGroup)
	DefaultProcessor = new(p.Processor)
	DefaultProcessor.Init()
	ConnManager = new(NetConnManager)
	ConnManager.Init()
	ModManager = new(MManager)
	ModManager.Init()
	ServerManager = new(NetServerManager)
	ServerManager.Init()
	ClientManager = new(NetClientManager)
	ClientManager.Init()
	redis.Init()
	Register()
}

func Register() {
	DefaultProcessor.SetHandler("ServerConnectOK", proto.ServerConnect{}, HandleServerConnectOK)
	DefaultProcessor.SetHandler("ServerLoginRQ", proto.ServerLogInReq{}, HandleServerLoginRQ)
	DefaultProcessor.SetHandler("ServerLoginRS", proto.ServerLogInRes{}, HandleServerLoginRS)
	DefaultProcessor.SetHandler("ServerTick", proto.ServerTick{}, HandleServerTick)
	DefaultProcessor.SetHandler("ServerTickBack", proto.ServerTick{}, HandleServerTickBack)
}

func Run() {

	ModManager.Run()
	ConnManager.Run(wg)
	ClientManager.Run()
	ServerManager.Run()
}

func Close() {
	ServerManager.Close()
	ClientManager.Close()
	ModManager.Close()
	ConnManager.Close()
	redis.R.OnClose()
	wg.Wait()
}
