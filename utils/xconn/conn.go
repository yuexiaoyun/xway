package xconn

import (
	"net"
)

type Tcp struct {
	Addr string
	net.Conn
}

func NewTcp(addr string) (*Tcp, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Tcp{Addr: addr, Conn: conn}, nil
}

func (xc *Tcp) Write(p []byte) (int, error) {
	var n int
	var err error
	n, err = xc.Conn.Write(p)
	if err != nil {
		// 尝试重连
		var conn net.Conn
		conn, err = net.Dial("tcp", xc.Addr)
		if err != nil {
			return 0, err
		}
		xc.Conn = conn
		n, err = xc.Conn.Write(p)
	}
	return n, err
}
