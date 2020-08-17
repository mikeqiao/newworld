package manager

import (
	"sync"

	conf "github.com/mikeqiao/newworld/config"
	"github.com/mikeqiao/newworld/http/httpserver"
	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net"
)

type NetServerManager struct {
	mutex sync.RWMutex
	wg    *sync.WaitGroup
	SList map[uint64]*net.TCPServer

	HttpS map[string]*httpserver.Server
}

func (n *NetServerManager) Init() {
	n.wg = new(sync.WaitGroup)
	n.SList = make(map[uint64]*net.TCPServer)
	n.HttpS = make(map[string]*httpserver.Server)
}

func (n *NetServerManager) AddHttpServer(addr string, server *httpserver.Server) {
	if _, ok := n.HttpS[addr]; ok {
		log.Error("already have this address:%v httpserver", addr)
		return
	}
	if nil == server {
		log.Error("this address:%v httpserver is nil", addr)
		return
	}

	n.HttpS[addr] = server
}

func (n *NetServerManager) NewNetServer(uid uint64, ctype uint32, name, addr string) {
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
	newServer.Ctype = ctype
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
	for _, v := range n.SList {
		if nil != v {
			v.Close()
		}
	}
	n.mutex.Unlock()
	n.wg.Wait()
}

func (n *NetServerManager) Run() {
	for _, v := range conf.Conf.Servers {
		n.NewNetServer(v.Uid, v.CType, v.Name, v.ListenAddr)
	}
}

func NewHtpServer(addr string) *httpserver.Server {
	s := new(httpserver.Server)
	s.Init(addr)
	ServerManager.AddHttpServer(addr, s)
	return s
}
