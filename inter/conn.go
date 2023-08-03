package inter

import (
	"net"
)

type Conn interface {
	Close()
	Listen()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetId(id uint64)
	ConnId() uint64
	AsyncWrite(bytes []byte)
}
