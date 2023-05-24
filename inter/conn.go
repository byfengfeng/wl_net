package inter

import (
	"net"
)

type Conn interface {
	Close()
	Listen()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	ConnId() int64
	AsyncWrite(bytes []byte)
}
