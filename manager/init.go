package manager

import (
	"github.com/mikeqiao/newworld/log"
	"sync"

	"github.com/mikeqiao/newworld/db/redis"

	p "github.com/mikeqiao/newworld/net/processor/protobuff"
)

var wg *sync.WaitGroup
var ServerManager *NetServerManager
var ClientManager *NetClientManager

func Init() {
	wg = new(sync.WaitGroup)
	ModManager = new(MManager)
	ModManager.Init()
	DefaultProcessor = new(p.Processor)
	DefaultProcessor.Init(ModManager.BaseMod)
	ConnManager = new(NetConnManager)
	ConnManager.Init()
	ServerManager = new(NetServerManager)
	ServerManager.Init()
	ClientManager = new(NetClientManager)
	ClientManager.Init()
	redis.Init()
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
	if nil != redis.R {
		redis.R.OnClose()
	}
	wg.Wait()
}

func RegisterMod(modName ...string) {
	err := ModManager.RegisterMod(modName...)
	if nil != err {
		log.Fatal("err:%v", err)
	}
}
