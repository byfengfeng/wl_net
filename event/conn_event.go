package event

import (
	"wl_net/inter"
	"wl_net/pool"
)

var (
	ConnMsgEventPool = pool.NewPool(func() any {
		return &ConnMsgEvent{}
	})
	ErrConnEventPool = pool.NewPool(func() any {
		return &ErrConnEvent{}
	})
)

// ErrConnEvent conn err event
type ErrConnEvent struct {
	Conn    inter.Conn
	IsClose bool
}

func NewErrConnEvent(conn inter.Conn, IsClose bool) *ErrConnEvent {
	connMsgEvent := ErrConnEventPool.Get().(*ErrConnEvent)
	connMsgEvent.Conn = conn
	connMsgEvent.IsClose = IsClose
	return connMsgEvent
}

// ConnMsgEvent conn read byte
type ConnMsgEvent struct {
	Data []byte
	Conn inter.Conn
}

func NewConnMsgEvent(conn inter.Conn, data []byte) *ConnMsgEvent {
	connMsgEvent := ConnMsgEventPool.Get().(*ConnMsgEvent)
	connMsgEvent.Conn = conn
	connMsgEvent.Data = data
	return connMsgEvent
}
