package module

import (
	"github.com/mikeqiao/newworld/common"
	"github.com/mikeqiao/newworld/data"
	"github.com/mikeqiao/newworld/net"
	"sync"

	"github.com/mikeqiao/newworld/config"
)

type Handler func(net.UserData, []byte, data.MemoryData) (err error)

type Mod struct {
	Name        string                  //mod type
	Uid         uint64                  //实例化的 modId
	Version     uint32                  //当前 mod 的版本
	FuncList    map[string]*ServiceFunc //mod  提供的服务列表
	modCloseSig chan bool               //mod 模块关闭信号
	ModState    common.ModState         //mod 的状态
	roomList    map[uint64]*GORoom      //模块的房间数量，每个房间一个 go 协程
	MaxCallLen  uint32                  //请求队列长度
}

type ServiceFunc func([]byte, *net.UserData, data.MemoryData)

func (m *Mod) Init() {
	m.Version = config.Conf.Version
	m.FuncList = make(map[string]*ServiceFunc)
	m.roomList = make(map[uint64]*GORoom)
	m.modCloseSig = make(chan bool, 1)
	m.ModState = common.Mod_Init
}

func (m *Mod) Start(wg *sync.WaitGroup) {

}

func (m *Mod) Route(funcName string, uData *net.UserData, msg interface{}, processor net.Processor) {

}

func (m *Mod) Close() {

}
