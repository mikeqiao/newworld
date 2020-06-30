package manager

import (
	"sync"
	"time"

	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
)

type NetClientManager struct {
	ConnectInterval time.Duration // 请求链接的间隔
	AutoReconnect   bool          // 是否重新连接
	mutex           sync.RWMutex
	wg              *sync.WaitGroup
	CList           map[uint64]*net.TCPClient
}

func (n *NetClientManager) Init() {
	if n.ConnectInterval <= 0 {
		n.ConnectInterval = 3 * time.Second
		log.Debug("invalid ConnectInterval, reset to %v", n.ConnectInterval)
	}
	n.AutoReconnect = true
	n.wg = new(sync.WaitGroup)
	n.CList = make(map[uint64]*net.TCPClient)
}

func (n *NetClientManager) NewClient(uid uint64, name, addr string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if s, ok := n.CList[uid]; ok && nil != s {
		log.Error("this uid:%v client already working", uid)
		return
	}
	newClient := new(net.TCPClient)
	newClient.UId = uid
	newClient.Name = name
	newClient.Addr = addr
	newClient.Processor = DefaultProcessor
	newClient.CreateAgent = CreateAgent
	newClient.Start()
	n.CList[uid] = newClient
}

func (n *NetClientManager) Run() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	for k, v := range n.CList {
		if nil != v {
			go v.Run(n.wg)
		} else {
			log.Error("this Client  is nil:%v", k)
		}
	}
}
