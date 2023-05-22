package udp

import (
	"context"
	"github.com/Byfengfeng/wl_net/conn"
	"github.com/Byfengfeng/wl_net/enum"
	"github.com/Byfengfeng/wl_net/event"
	"github.com/Byfengfeng/wl_net/inter"
	"github.com/Byfengfeng/wl_net/log"
	"go.uber.org/zap"
	"net"
)

type ListenUdp struct {
	address string
	*AddrConn
}

type AddrConn struct {
	conn          *net.UDPConn
	ctx           context.Context
	cancelFunc    context.CancelFunc
	ch            chan<- any
	oCh           <-chan struct{}
	writeCh       chan *AddrMsg
	Id            uint64
	connType      enum.ConnType
	closeReadChan chan struct{}
	codec         inter.Codec
}

type AddrMsg struct {
	Addr  net.Addr
	Bytes []byte
}

func NewUcpListen(addr string, ch chan<- any, oCh <-chan struct{}, codec inter.Codec) *ListenUdp {
	return &ListenUdp{address: addr, AddrConn: NewAddrConn(ch, oCh, nil, codec)}
}

func NewAddrConn(ch chan<- any, oCh <-chan struct{}, conn *net.UDPConn, codec inter.Codec) *AddrConn {
	return &AddrConn{
		conn:          conn,
		ch:            ch,
		oCh:           oCh,
		writeCh:       make(chan *AddrMsg),
		closeReadChan: make(chan struct{}),
		connType:      enum.Ready,
		codec:         codec,
	}
}

func (u *ListenUdp) Start() {
	addr, err := net.ResolveUDPAddr("udp", u.address)
	if err != nil {
		log.Error("udp listener addr exit", zap.Any("err", err.Error()))
		return
	}
	u.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Error("udp listener exit", zap.Any("err", err.Error()))
		return
	}
	log.Logger.Info("udp listen start success")
	u.Listen()
	return
}

func (a *AddrConn) Listen() {
	go a.read()
	go a.write()
}

func (a *AddrConn) read() {
	headLen := make([]byte, 65535)
	for {
		select {
		case <-a.oCh:
			return
		case <-a.closeReadChan:
			return
		default:
			_, udpAddr, err := a.conn.ReadFrom(headLen)
			if err != nil {
				a.ch <- event.NewErrConnEvent(conn.DelAddr(udpAddr), true)
				continue
			}

			length, data := a.codec.Decode(headLen[:enum.HeadSize])

			if length > uint32(enum.HeadSize) {
				copy(data, headLen[enum.HeadSize:])
			}
			connect, b := conn.GetAddr(a.conn.LocalAddr(), udpAddr, a.codec, a.AsyncUdpWrite)
			if b {
				a.ch <- connect
			}
			a.ch <- event.NewConnMsgEvent(connect, data)
		}
	}
}

func (a *AddrConn) write() {
	for {
		select {
		case <-a.closeReadChan:
			return
		case msg := <-a.writeCh:
			_, err := a.conn.WriteTo(msg.Bytes, msg.Addr)
			if err != nil {
				a.ch <- event.NewErrConnEvent(conn.DelAddr(msg.Addr), true)
			}
		}
	}

}

func (a *AddrConn) AsyncUdpWrite(addr net.Addr, bytes []byte) {
	a.writeCh <- &AddrMsg{addr, bytes}
}

func (a *AddrConn) Action(event any) {
	a.ch <- event
}
