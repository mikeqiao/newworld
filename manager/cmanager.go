package manager

import (
	"sync"
	"time"

	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
	p "github.com/mikeqiao/newworld/net/processor/protobuff"
)

var DefaultProcessor *p.Processor
var ConnManager *NetConnManager
var CreatID uint64

func CreateAgent(conn *net.TCPConn, tp net.Processor, uid uint64, c chan bool) *net.TcpAgent {
	a := new(net.TcpAgent)

	CreatID += 1
	a.Init(conn, tp, uid, CreatID, c)
	ConnManager.AddAgent(a)

	return a
}

type NetConnManager struct {
	isClose bool
	mutex   sync.RWMutex
	wg      *sync.WaitGroup
	CList   map[uint64]*net.TcpAgent
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

		agent.Start(n.wg)
		n.wg.Add(1)
		agent.Run()
		agent.Close()
		n.wg.Done()
		//从map 删除
	}()
}

func (n *NetConnManager) Run(wg *sync.WaitGroup) {
	go n.Updata(wg)
}

func (n *NetConnManager) Updata(wg *sync.WaitGroup) {
	wg.Add(1)
	t1 := time.NewTimer(time.Second * 1)
	for {
		select {
		case <-t1.C:
			if n.isClose == true {
				goto Loop
			}
			n.mutex.Lock()

			for k, v := range n.CList {
				//v.agent.OnClose()

				if v.IsClose() == true {
					log.Debug("userConn close, id:%v", k)
					delete(n.CList, k)
				}
			}
			n.mutex.Unlock()

			t1.Reset(time.Second * 1)
		}

	}

Loop:
	wg.Done()
}

func (n *NetConnManager) Close() {
	n.mutex.Lock()

	for k, v := range n.CList {
		//v.agent.OnClose()

		if nil != v {
			v.Close()
			delete(n.CList, k)
		}
	}
	n.mutex.Unlock()
	n.wg.Wait()
}
