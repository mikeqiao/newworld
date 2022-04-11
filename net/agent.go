package net

import (
	"github.com/golang/protobuf/proto"
	"github.com/mikeqiao/newworld/common"
	"sync"
	"time"

	"github.com/mikeqiao/newworld/log"
	"github.com/mikeqiao/newworld/net/base_proto"
)

const (
	//默认长度
	Default_WriteNum = uint32(10240)

	Agent_State     = int32(iota)
	Agent_Init      //初始化完成
	Agent_UnChecked //未验证合法化
	Agent_Running   //正在运行工作
	Agent_Closing   //关闭中
	Agent_Closed    //已经关闭
)

type TcpAgent struct {
	//基础配置属性
	Processor       Processor //序列化，反序列化 和 route msg
	PendingWriteNum uint32    //配置 chan 容量
	UnCheckLifetime int64     //链接未验证有效期时间（秒）
	TickLifetime    int64     //没心跳断开时间（秒）
	//基础属性
	Name      string      //所在服务的名字
	UId       uint64      //远端 uid (remote IP,port 的唯一标识符)
	LocalUId  uint64      //本端 uid (local IP,port 的唯一标识符)
	conn      Conn        //网络连接
	closeSign chan bool   //关闭通知 chan
	WriteChan chan []byte //待处理发送消息 chan

	//生效控制类属性
	startTime  int64    //链接开始的时间戳
	tick       int64    //上次心跳的时间戳
	agentState int32    //agent 状态
	RemoteMod  []string //连接对方提供的 mod 模块
	//附加属性
	userData interface{} //和agent 绑定的相关 数据
}

func (a *TcpAgent) Init(conn *TCPConn, tp Processor, writeNum uint32, unCheckLifetime, tickLifetime int64, c chan uint64, name string, localUid uint64) {
	a.conn = conn
	a.Processor = tp
	a.PendingWriteNum = Default_WriteNum
	if writeNum > 0 {
		a.PendingWriteNum = writeNum
	}
	a.WriteChan = make(chan []byte, a.PendingWriteNum)
	a.closeSign = make(chan bool, 1)

	//状态属性
	a.UnCheckLifetime = 5
	if unCheckLifetime > 0 {
		a.UnCheckLifetime = unCheckLifetime
	}
	a.TickLifetime = 20
	if tickLifetime > 0 {
		a.TickLifetime = tickLifetime
	}
	a.agentState = Agent_Init
	a.Name = name
	a.LocalUId = localUid
}

func (a *TcpAgent) GetRemoteAddr() string {
	if nil == a.conn {
		return ""
	}
	return a.conn.RemoteAddr().String()
}

func (a *TcpAgent) GetLocalAddr() string {
	if nil == a.conn {
		return ""
	}
	return a.conn.LocalAddr().String()
}

func (a *TcpAgent) SetUID(uid uint64) {
	a.UId = uid
}

func (a *TcpAgent) Check(uid uint64, modList []string) {
	a.UId = uid
	a.RemoteMod = modList
	a.agentState = Agent_Running
}

func (a *TcpAgent) SetTick(time int64) {
	a.tick = time
}

func (a *TcpAgent) GetState() int32 {
	return a.agentState
}

func (a *TcpAgent) SetBindData(date interface{}) {
	a.userData = date
}

func (a *TcpAgent) Start(wg *sync.WaitGroup) {
	a.startTime = time.Now().Unix()
	a.tick = time.Now().Unix()
	go a.Run(wg)
	go a.Update(wg)
	//发送登陆消息
	a.SendLogInCheck()
}

func (a *TcpAgent) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Error("agent:UID:%v, RemoteAddr:%s, ,read message err: %v", a.UId, a.conn.RemoteAddr().String(), err)
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
	wg.Done()
	a.DoClose()
}

func (a *TcpAgent) Update(wg *sync.WaitGroup) {
	a.agentState = Agent_UnChecked
	wg.Add(1)
	t1 := time.NewTimer(time.Second * 1)
	t2 := time.NewTimer(time.Second * time.Duration(a.TickLifetime))
	for {
		select {
		case msg, isClosed := <-a.WriteChan:
			if !isClosed {
				goto Loop
			}
			if nil != msg {
				err := a.conn.DoWrite(msg)
				if err != nil {
					log.Error("conn.DoWrite err:%v", err)
					goto Loop
				}
			}
		case <-a.closeSign:
			if a.agentState != Agent_Closing {
				a.agentState = Agent_Closing
				close(a.WriteChan)
			}
		case <-t1.C:
			if a.agentState == Agent_UnChecked {
				nowTime := time.Now().Unix()
				if (a.startTime + a.UnCheckLifetime) < nowTime {
					log.Debug("outTime to not Checked: %v", a.conn.RemoteAddr())
					a.DoClose()
				}
			} else if a.agentState == Agent_Running {
				nowTime := time.Now().Unix()
				if (a.tick + a.TickLifetime*3) < nowTime {
					log.Debug("outTime to no tick: %v", a.conn.RemoteAddr())
					a.DoClose()
				}
			}
			t1.Reset(time.Second * 1)
		case <-t2.C:
			a.SendTick()
			t2.Reset(time.Second * 10)
		}
	}
Loop:
	wg.Done()
	a.Close()
}

func (a *TcpAgent) Close() {
	if a.agentState != Agent_Closed {
		a.agentState = Agent_Closed
		a.conn.Close()
	}
}

func (a *TcpAgent) DoClose() {
	if a.agentState != Agent_Closing {
		a.closeSign <- true
	}
}

func (a *TcpAgent) WriteMsg(u *CallData, msg interface{}) {
	if a.Processor != nil {
		data, err := a.Processor.Marshal(u, msg)
		if err != nil {
			log.Error("marshal message ud:%v error: %v", u, err)
			return
		}
		writeData, err := a.conn.WriteMsg(data)
		if err != nil {
			log.Error("write message ud:%v error: %v", u, err)
		} else if nil != writeData {
			if len(a.WriteChan) == cap(a.WriteChan) {
				log.Debug("close conn: channel full")
				a.DoClose()
				return
			}
			a.WriteChan <- writeData
		}
	}
}

func (a *TcpAgent) UnMarshalMsg(msg proto.Message, data []byte) error {
	return a.Processor.UnMarshalMsg(msg, data)
}

func (a *TcpAgent) SendTick() {
	t := new(base_proto.ServerTick)
	t.Time = uint32(time.Now().Unix())
	t.State = a.agentState
	u := new(CallData)
	u.Mod = common.Mod_Base
	u.Function = common.BaseMod_AgentTick
	a.SendMsg(u, t)
}

func (a *TcpAgent) SendLogInCheck() {
	t := new(base_proto.ServerLogInCheckReq)
	t.Name = a.Name
	t.Sid = a.LocalUId
	u := new(CallData)
	u.Mod = common.Mod_Base
	u.Function = common.BaseMod_AgentCheck
	a.SendMsg(u, t)
}

func (a *TcpAgent) SendMsg(u *CallData, msg interface{}) {
	if a.agentState != Agent_Closing && a.agentState != Agent_Closed {
		a.WriteMsg(u, msg)
	}
}
