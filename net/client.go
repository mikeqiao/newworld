package net

import (
	"net"
	"sync"
	"time"

	"github.com/mikeqiao/newworld/net/proto"

	"github.com/mikeqiao/newworld/log"
)

type TCPClient struct {
	UId             uint64
	Addr            string // 地址
	Name            string
	ConnectInterval time.Duration                                                  // 请求链接的间隔
	CreateAgent     func(*TCPConn, Processor, uint64, uint32, chan bool) *TcpAgent // 代理
	Closed          bool                                                           // 关闭标识符
	Working         bool
	Processor       Processor
	CloseChannel    chan bool
	//msg parser
	Agent *TcpAgent
}

func (this *TCPClient) init() {
	if this.ConnectInterval <= 0 {
		this.ConnectInterval = 3 * time.Second
		log.Debug("invalid ConnectInterval, reset to %v", this.ConnectInterval)
	}
	if this.CreateAgent == nil {
		log.Error("CreateAgent must not be nil")
		return
	}
	this.CloseChannel = make(chan bool, 1)
	this.Closed = false
}

func (this *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", this.Addr)
		if err == nil || this.Closed {
			return conn
		}
		log.Release("connect to %v error: %v", this.Addr, err)
		time.Sleep(this.ConnectInterval)
		continue
	}
}

func (this *TCPClient) connect() {

	conn := this.dial()
	if conn == nil {
		log.Error("conn is nil")
		return
	}
	if this.Closed {
		conn.Close()
		log.Debug("this is close")
		return
	}
	tcpConn := newTCPConn(conn)
	agent := this.CreateAgent(tcpConn, this.Processor, this.UId, 1, this.CloseChannel)
	this.Agent = agent
	agent.SetLocalUID(this.UId)
	log.Debug("client connect ok:%v", this.Addr)
	//通知上层，连接成功，开始登录流程
	tmsg := new(proto.ServerConnect)
	tmsg.Uid = this.UId
	this.Processor.Handle("ServerConnectOK", tmsg, &UserData{UId: this.UId, Agent: this.Agent})
}

func (this *TCPClient) Start() {
	this.init()
	this.connect()

}

func (this *TCPClient) ReStart() {
	this.connect()

}

func (this *TCPClient) Close() {
	this.Closed = true
	this.Agent.Close()
}

func (this *TCPClient) Run(wg *sync.WaitGroup) {
	if this.Working {
		return
	}
	wg.Add(1)
	defer wg.Done()
	for {
		this.Working = true
		select {
		case <-this.CloseChannel:
			if this.Closed {
				return
			} else {
				this.ReStart()
			}
		}
	}

}
