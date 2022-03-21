package net

import (
	"encoding/binary"
	"errors"
	"github.com/mikeqiao/newworld/log"
	"io"
	"math"
)

type MessageParser interface {
	//设置 消息生成配置
	SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32)
	//设置大小端
	SetByteOrder(littleEndian bool)
	//读取规则 读取数据
	Read(conn *TCPConn) ([]byte, error)
	//写入规则，生成数据
	Write(data []byte) ([]byte, error)
}

type MsgParser struct {
	lenMsgLen    int
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

var DefaultMsgParser *MsgParser

func init() {
	DefaultMsgParser = NewMsgParser()
}

func NewMsgParser() *MsgParser {
	p := &MsgParser{
		lenMsgLen:    4,
		minMsgLen:    1,
		maxMsgLen:    4096,
		littleEndian: true,
	}
	return p
}

func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	} else {
		log.Error("Invalid msgLen value :%v", lenMsgLen)
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	} else {
		log.Warning("Not set minMsgLen value, and default value :%v", lenMsgLen)
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	} else {
		log.Warning("Not set maxMsgLen value, and default value :%v", maxMsgLen)
	}
	var max uint32

	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) ([]byte, error) {
	var b [4]byte
	bufMsgLen := b[:p.lenMsgLen]
	//read len (io.ReadFull // 读取指定长度的字节)
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}
	//parse len
	var msgLen uint32
	switch p.lenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if p.littleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if p.littleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}
	//check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	//data
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}
	return msgData, nil
}

// goroutine safe
func (p *MsgParser) Write(data []byte) ([]byte, error) {
	//get len
	var msgLen = uint32(len(data))
	//check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	msg := make([]byte, uint32(p.lenMsgLen)+msgLen)
	// write len
	switch p.lenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if p.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if p.littleEndian {
			binary.LittleEndian.PutUint32(msg, uint32(msgLen))
		} else {
			binary.BigEndian.PutUint32(msg, uint32(msgLen))
		}
	}
	// write data
	l := p.lenMsgLen
	copy(msg[l:], data)
	return msg, nil
}
