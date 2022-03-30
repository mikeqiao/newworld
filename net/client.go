package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/newworld/net/base_proto"

	"github.com/mikeqiao/newworld/log"
)

type TCPClient struct {
	UId             uint64
	Addr            string // 地址
	Name            string
	ConnectInterval time.Duration                                                             // 请求链接的间隔
	CreateAgent     func(*TCPConn, Processor, uint32, int64, int64, string, uint64) *TcpAgent // 代理
	Closed          bool                                                                      // 关闭标识符
	Working         bool
	Processor       Processor
	CloseChannel    chan uint64
	PendingWriteNum uint32 //配置 chan 容量
	//msg parser
	Agent *TcpAgent
}

func (t *TCPClient) init() {
	if t.ConnectInterval <= 0 {
		t.ConnectInterval = 3 * time.Second
		log.Debug("invalid ConnectInterval, reset to %v", t.ConnectInterval)
	}
	if t.CreateAgent == nil {
		log.Error("CreateAgent must not be nil")
		return
	}
	t.CloseChannel = make(chan uint64, 1)
	t.Closed = false
}

func (t *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", t.Addr)
		if err == nil || t.Closed {
			return conn
		}
		log.Release("connect to %v error: %v", t.Addr, err)
		time.Sleep(t.ConnectInterval)
		continue
	}
}

func (t *TCPClient) connect() {
	conn := t.dial()
	if conn == nil {
		log.Error("conn is nil")
		return
	}
	if t.Closed {
		_ = conn.Close()
		log.Debug("this is close")
		return
	}
	tcpConn := newTCPConn(conn)
	agent := t.CreateAgent(tcpConn, t.Processor, t.PendingWriteNum, 0, 0, t.Name, t.UId)
	t.Agent = agent
	agent.SetUID(t.UId)
	log.Debug("client connect ok:%v", t.Addr)
	//通知上层，连接成功，开始登录流程
	tMsg := new(base_proto.ServerConnect)
	tMsg.Uid = t.UId
	//	t.Processor.Handle("ServerConnectOK", tMsg, &CallData{Uid: t.UId, Agent: t.Agent})
}

func (t *TCPClient) Start(wg *sync.WaitGroup) {
	t.init()
	t.connect()
	t.Run(wg)
}

func (t *TCPClient) ReStart() {
	t.connect()
}

func (t *TCPClient) Close() {
	t.Closed = true
	t.Agent.DoClose()
}

func (t *TCPClient) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	if t.Working {
		return
	}
	t.Working = true
	for {
		select {
		case <-t.CloseChannel:
			if t.Closed {
				return
			} else {
				t.ReStart()
			}
		}
	}
}
