package net

import (
	"errors"
	"net"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg(args []byte) ([]byte, error)
	DoWrite(data []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
}

type ConnList map[net.Conn]struct{}

type TCPConn struct {
	conn      net.Conn
	msgParser MessageParser
}

func newTCPConn(conn net.Conn) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.msgParser = DefaultMsgParser
	return tcpConn
}

func (c *TCPConn) DoWrite(data []byte) error {
	if data != nil {
		_, err := c.conn.Write(data)
		return err
	}
	return errors.New("nil data")
}

func (c *TCPConn) Close() {
	_ = c.conn.Close()
}

//为了保证 符合 Reader  接口类型
func (c *TCPConn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *TCPConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TCPConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TCPConn) ReadMsg() ([]byte, error) {
	return c.msgParser.Read(c)
}

func (c *TCPConn) WriteMsg(data []byte) ([]byte, error) {
	return c.msgParser.Write(data)
}
