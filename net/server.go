package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/newworld/log"
)

type TCPServer struct {
	UId             uint64       //监听服务唯一ID
	Name            string       // 监听服务名称
	Addr            string       // 监听的地址端口
	ln              net.Listener // 监听
	PendingWriteNum uint32       //配置 chan 容量
	//agent
	Processor   Processor
	CreateAgent func(*TCPConn, Processor, uint32, int64, int64) *TcpAgent // 代理
}

func (t *TCPServer) init() {
	ln, err := net.Listen("tcp", t.Addr)
	if err != nil {
		log.Error("%v", err)
	}
	if t.CreateAgent == nil {
		log.Error("CreateAgent must not be nil")
	}
	t.ln = ln
}

func (t *TCPServer) run(wg *sync.WaitGroup) {
	wg.Add(1)
	var tempDelay time.Duration
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
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
			goto Loop
		}
		tempDelay = 0
		//创建了新的链接  创建 agent 加入 conn管理
		tcpConn := newTCPConn(conn)
		t.CreateAgent(tcpConn, t.Processor, t.PendingWriteNum, 0, 0)
	}
Loop:
	wg.Done()
}

func (t *TCPServer) Start(wg *sync.WaitGroup) {
	t.init()
	go t.run(wg)
}

func (t *TCPServer) Close() {
	_ = t.ln.Close()
}
