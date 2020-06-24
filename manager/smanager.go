package manager

import (
	"sync"

	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
)

type NetServerManager struct {
	mutex sync.RWMutex
	wg    *sync.WaitGroup
	SList map[uint64]*net.TCPServer
}

func (n *NetServerManager) Init() {
	n.wg = new(sync.WaitGroup)
	n.SList = make(map[uint64]*net.TCPServer)
}

func (n *NetServerManager) NewNetServer(uid uint64, name, addr string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if s, ok := n.SList[uid]; ok && nil != s {
		log.Error("this uid:%v server already working", uid)
		return
	}
	newServer := new(net.TCPServer)
	newServer.UId = uid
	newServer.Name = name
	newServer.Addr = addr
	newServer.Processor = DefaultProcessor
	newServer.CreateAgent = CreateAgent
	newServer.Start(n.wg)
	n.SList[uid] = newServer
}

func (n *NetServerManager) CloseNetServer(uid uint64) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if s, ok := n.SList[uid]; !ok {
		log.Error("this uid:%v server not working", uid)
		return
	} else {
		if nil != s {
			s.Close()
		}
	}

	delete(n.SList, uid)

}

func (n *NetServerManager) Close() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	for _, v := range n.SList {
		if nil != v {
			v.Close()
		}
	}
	n.wg.Wait()
}
