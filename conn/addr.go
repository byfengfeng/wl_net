package conn

import (
	"github.com/Byfengfeng/wl_net/inter"
	"net"
	"sync"
)

var (
	addrMap = sync.Map{}
)

type addr struct {
	localAddr  net.Addr
	remoteAddr net.Addr
	inter.Codec
	handler func(addr net.Addr, data []byte)
}

func GetAddr(local, remote net.Addr, codec inter.Codec, handler func(addr net.Addr, data []byte)) (inter.Conn, bool) {
	conn, ok := addrMap.Load(remote.String())
	if ok {
		return conn.(inter.Conn), false
	}
	conn = newAddr(local, remote, codec, handler)
	addrMap.Store(remote.String(), conn)
	return conn.(inter.Conn), true
}

func DelAddr(remote net.Addr) inter.Conn {
	conn, ok := addrMap.LoadAndDelete(remote.String())
	if ok {
		return conn.(inter.Conn)
	}
	return &addr{remoteAddr: remote}
}

func newAddr(local, remote net.Addr, codec inter.Codec, handler func(addr net.Addr, data []byte)) inter.Conn {
	return &addr{local, remote, codec, handler}
}

func (a *addr) Close() {
	DelAddr(a.remoteAddr)
}

func (a *addr) Listen() {}

func (a *addr) LocalAddr() net.Addr {
	return a.localAddr
}

func (a *addr) RemoteAddr() net.Addr {
	return a.remoteAddr
}

func (a *addr) AsyncWrite(bytes []byte) {
	a.handler(a.remoteAddr, a.Encode(bytes))
}
