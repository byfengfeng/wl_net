package net

import (
	"fmt"
	"github.com/Byfengfeng/wl_net/conn"
	"github.com/Byfengfeng/wl_net/enum"
	"github.com/Byfengfeng/wl_net/event"
	"github.com/Byfengfeng/wl_net/inter"
	"github.com/Byfengfeng/wl_net/listen"
	"github.com/Byfengfeng/wl_net/pool"
	"net"
)

type EventDialLoop struct {
	addr        string
	connCount   int32
	ch          chan any
	evCh        chan struct{}
	netType     enum.NetType
	conn        inter.Conn
	eventConnFn inter.EventDialHandler
	codec       inter.Codec
}

func NewEventDialLoop(addr string, port int32, netType enum.NetType, ev inter.EventDialHandler) *EventDialLoop {
	return &EventDialLoop{
		addr:        fmt.Sprintf("%s:%d", addr, port),
		eventConnFn: ev,
		codec:       pool.NewCodec(),
		netType:     netType,
		ch:          make(chan any, 20),
		evCh:        make(chan struct{}),
	}
}

func (e *EventDialLoop) Run() {
	go e.evLoop()
	e.Accept(listen.Dial(e.netType, e.addr))
}

func (e *EventDialLoop) evLoop() {
	for {
		select {
		case message := <-e.ch:
			switch event := message.(type) {
			case *event.ErrConnEvent:
				e.deleteConn(event.Conn)
				return
			case *conn.Conn:
				e.addConn(event)
			case *event.ConnMsgEvent:
				e.eventConnFn.DialReact(event.Conn, event.Data)
				e.PutConnMsgEvent(event)
			}
		}
	}
}

func (e *EventDialLoop) Send(byte []byte) {
	if e.conn != nil {
		e.conn.AsyncWrite(byte)
	}
}

func (e *EventDialLoop) Accept(stdConn net.Conn) {
	connect := conn.NewConn(e.ch, e.evCh, stdConn, e.codec, e.netType)
	e.ch <- connect
	connect.Listen()
}

func (e *EventDialLoop) addConn(conn inter.Conn) {
	e.conn = conn
	e.eventConnFn.OnDialOpened(conn)
}

func (e *EventDialLoop) deleteConn(conn inter.Conn) {
	e.eventConnFn.OnDialClose(conn)
	conn.Close()
	e.Run()
}

func (e *EventDialLoop) Close() {
	e.ch <- event.NewErrConnEvent(e.conn, true)
}

func (e *EventDialLoop) PutConnMsgEvent(ev *event.ConnMsgEvent) {
	event.ConnMsgEventPool.Put(ev)
}
