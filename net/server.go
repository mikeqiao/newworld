package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/newworld/log"
)

type TCPServer struct {
	UId  uint64       //监听服务唯一ID
	Name string       // 监听服务名称
	Addr string       // 监听的地址端口
	ln   net.Listener // 监听
	//agent
	Processor   Processor
	CreateAgent func(*TCPConn, Processor, uint64) *TcpAgent // 代理
}

func (this *TCPServer) init() {
	ln, err := net.Listen("tcp", this.Addr)
	if err != nil {
		log.Error("%v", err)
	}
	if this.CreateAgent == nil {
		log.Error("CreateAgent must not be nil")
	}
	this.ln = ln
}

func (this *TCPServer) run(wg *sync.WaitGroup) {
	wg.Add(1)

	var tempDelay time.Duration
	for {
		conn, err := this.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Error("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			wg.Done()
			return
		}
		tempDelay = 0
		//创建了新的链接  创建 agent 加入 conn管理
		tcpConn := newTCPConn(conn)
		this.CreateAgent(tcpConn, this.Processor, this.UId)

	}
	wg.Done()
}

func (this *TCPServer) Start(wg *sync.WaitGroup) {
	this.init()
	go this.run(wg)
}

func (this *TCPServer) Close() {
	this.ln.Close()
}
