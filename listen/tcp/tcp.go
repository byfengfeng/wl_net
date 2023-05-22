package tcp

import (
	"go.uber.org/zap"
	"net"
	"wl_net/log"
)

type ListenTcp struct {
	*net.TCPListener
	address    string
	connHandel func(conn net.Conn)
}

func NewTcpListen(addr string, connHandel func(conn net.Conn)) *ListenTcp {
	return &ListenTcp{address: addr, connHandel: connHandel}
}

func (n *ListenTcp) Start() {
	addr, err := net.ResolveTCPAddr("tcp", n.address)
	if err != nil {
		log.Error("tcp listener addr", zap.Any("err", err.Error()))
		return
	}
	n.TCPListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("tcp listener err", zap.Any("err", err.Error()))
		return
	}
	go func() {
		for {
			tcpConn, err := n.TCPListener.Accept()
			if err != nil {
				log.Error("tcp accept channel exit", zap.Any("err", err.Error()))
				return
			}
			go n.connHandel(tcpConn)
		}
	}()
	log.Info("start tcp listen ")
	return
}
