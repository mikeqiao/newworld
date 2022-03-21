package manager

import (
	"github.com/mikeqiao/newworld/net"
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
	"sync"
)

var DefaultProcessor *p.Processor
var ConnManager *NetConnManager

func CreateAgent(conn *net.TCPConn, tp net.Processor, writeNum uint32, unCheckLifetime, tickLifetime int64) *net.TcpAgent {
	if nil == ConnManager {
		return nil
	}
	a := new(net.TcpAgent)
	a.Init(conn, tp, writeNum, unCheckLifetime, tickLifetime, ConnManager.CloseUid)
	a.Start(ConnManager.wg)
	//	ConnManager.AddAgent(a)
	return a
}

type NetConnManager struct {
	isClose   bool
	mutex     sync.RWMutex
	wg        *sync.WaitGroup
	CList     map[uint64]*net.TcpAgent
	CloseUid  chan uint64
	CLoseSign chan bool
}

func (n *NetConnManager) Init() {
	n.CloseUid = make(chan uint64, 10)
	n.wg = new(sync.WaitGroup)
	n.CList = make(map[uint64]*net.TcpAgent)
}

func (n *NetConnManager) AddAgent(agent *net.TcpAgent) {
	if nil == agent {
		return
	}
	n.mutex.Lock()
	n.CList[agent.UId] = agent
	n.mutex.Unlock()

}

func (n *NetConnManager) Run(wg *sync.WaitGroup) {
	go n.Update(wg)
}

func (n *NetConnManager) Update(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		select {
		case <-n.CLoseSign:
			goto Loop
		case agentId, ok := <-n.CloseUid:
			if ok {
				n.mutex.Lock()
				delete(n.CList, agentId)
				n.mutex.Unlock()
			}
		}
	}
Loop:
	wg.Done()
}

func (n *NetConnManager) Close() {
	n.CLoseSign <- true
	n.mutex.Lock()
	for k, v := range n.CList {
		if nil != v {
			v.DoClose()
			delete(n.CList, k)
		}
	}
	n.mutex.Unlock()
	n.wg.Wait()
}
