package manager

import (
	"sync"

	"github.com/mikeqiao/newworld/net"
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
)

var DefaultProcessor *p.Processor
var ConnManager *NetConnManager
var CreatID uint64

func CreateAgent(conn *net.TCPConn, tp net.Processor, uid uint64) *net.TcpAgent {
	a := new(net.TcpAgent)

	CreatID += 1
	a.Init(conn, tp, uid, CreatID)
	ConnManager.AddAgent(a)

	return a
}

type NetConnManager struct {
	mutex sync.RWMutex
	wg    *sync.WaitGroup
	CList map[uint64]*net.TcpAgent
}

func (n *NetConnManager) Init() {
	n.wg = new(sync.WaitGroup)
	n.CList = make(map[uint64]*net.TcpAgent)
}

func (n *NetConnManager) AddAgent(agent *net.TcpAgent) {
	if nil == agent {
		return
	}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.CList[agent.LUId] = agent
	go func() {
		n.wg.Add(1)
		agent.Start()
		agent.Run()
		agent.Close()
		n.wg.Done()
		//从map 删除
	}()
}
