package net

import (
	"time"

	"github.com/mikeqiao/newworld/log"
)

var CreatID = uint64(1)

type TcpAgent struct {
	LUId         uint64 //本端local uid
	RUId         uint64 //对方端remote uid
	conn         Conn
	Processor    Processor
	lifetime     int64 //链接未验证有效期时间（秒）
	starttime    int64 //链接开始的时间戳
	tick         int64 //上次心跳的时间戳
	islogin      bool  //是否登陆验证过
	isClose      bool  //是否关闭
	Closed       bool
	Ctype        uint32 //连接类型 1 server  2 client
	userData     interface{}
	Mod          interface{} //agent 所属于的mod 的指针
	CloseChannel chan bool
}

func (a *TcpAgent) Init(conn *TCPConn, tp Processor, luid, ruid uint64, c chan bool) {
	if 0 != luid {
		a.LUId = luid
	}
	a.RUId = ruid
	a.conn = conn
	a.Processor = tp
	a.lifetime = 10
	a.starttime = time.Now().Unix()
	a.tick = time.Now().Unix()
	a.islogin = false
	a.isClose = false
	a.Closed = false
	a.CloseChannel = c
}

func (a *TcpAgent) SetLocalUID(uid uint64) {
	a.LUId = uid
}

func (a *TcpAgent) SetRemotUID(uid uint64) {
	a.RUId = uid
}

func (a *TcpAgent) SetMod(mod interface{}) {
	a.Mod = mod
}

func (a *TcpAgent) Start() {

}

func (a *TcpAgent) Run() {

	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Error("agent:localUID:%v, RemoteUID:%v, ,read message err: %v", a.LUId, a.RUId, err)
			goto Loop
		}
		if a.Processor != nil && nil != data {
			err := a.Processor.Unmarshal(a, data)
			if err != nil {
				log.Error("unmarshal message error: %v", err)
			}
		}
	}
Loop:
	a.Close()
}

func (a *TcpAgent) Close() {
	if a.Closed {
		return
	}
	a.Closed = true
	if nil != a.CloseChannel {
		a.CloseChannel <- true
	}
	a.conn.Close()
}

func (a *TcpAgent) WriteMsg(u *UserData, msg interface{}) {

	if a.Processor != nil {
		ud, data, err := a.Processor.Marshal(u, msg)
		if err != nil {
			log.Error("marshal message ud:%v error: %v", ud, err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message ud:%v error: %v", ud, err)
		}
	}

}
