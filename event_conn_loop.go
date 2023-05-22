package net

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"net"
	"wl_net/conn"
	"wl_net/enum"
	"wl_net/event"
	"wl_net/inter"
	"wl_net/listen/tcp"
	"wl_net/listen/udp"
	"wl_net/listen/web_socket"
	"wl_net/log"
	"wl_net/pool"
)

type EventConnLoop struct {
	addr        string
	connCount   int32
	ch          chan any
	evCh        chan struct{}
	conn        inter.Conn
	eventConnFn inter.EventConnHandler
	codec       inter.Codec
	*ants.Pool
	netType enum.NetType
}

func NewConnEventLoop(addr string, port int32, ev inter.EventConnHandler, netWork enum.NetType) *EventConnLoop {
	gPool, err := ants.NewPool(10)
	if err != nil {
		panic(err)
	}
	return &EventConnLoop{
		addr:        fmt.Sprintf("%s:%d", addr, port),
		eventConnFn: ev,
		codec:       pool.NewCodec(),
		ch:          make(chan any, 1),
		evCh:        make(chan struct{}, 1),
		netType:     netWork,
		Pool:        gPool,
	}
}

func (e *EventConnLoop) WithDecodeLength(headSize int) {
	enum.HeadSize = headSize
}

func (e *EventConnLoop) WithCodec(codec inter.Codec) {
	e.codec = codec
}

func (e *EventConnLoop) Run() {
	if e.netType == enum.Tcp {
		tcp.NewTcpListen(e.addr, e.accept).Start()
	} else if e.netType == enum.WebSocket {
		web_socket.NewWebSocketListen(e.addr, e.accept).Start()
	} else if e.netType == enum.Udp {
		udp.NewUcpListen(e.addr, e.ch, e.evCh, e.codec).Start()
	}
	go e.SubMsgHandel()
	e.evLoop()
}

// 主从 reactor
// main：用于处理accept连接建立及销毁
// sub通过select监听到多个conn消息，使用不同的handler进行处理，handler处理是在线程池中通过work执行的，sub可以有多个
// msg 处理和 accept拆分
// msg 处理放置线程池中执行
func (e *EventConnLoop) evLoop() {
	for {
		select {
		case <-e.evCh:
			close(e.ch)
			log.Info("evLoop close")
			return
		case message := <-e.ch:
			switch msg := message.(type) {
			case inter.Conn:
				e.eventConnFn.OnOpened(msg)
			}
		}
	}
}

func (e *EventConnLoop) SubMsgHandel() {
	for {
		select {
		case <-e.evCh:
			log.Info("SubMsgHandel close")
			return
		case message := <-e.ch:
			switch msg := message.(type) {
			case *event.ConnMsgEvent:
				e.Submit(func() {
					e.eventConnFn.React(msg.Conn, msg.Data)
					e.PutConnMsgEvent(msg)
				})
			case *event.ErrConnEvent:
				e.Submit(func() {
					e.deleteConn(msg.Conn, true)
					e.PutErrorConnEvent(msg)
				})
			}
		}
	}
}

func (e *EventConnLoop) accept(stdConn net.Conn) {
	connect := conn.NewConn(e.ch, e.evCh, stdConn, e.codec, e.netType)
	e.ch <- connect
	connect.Listen()
}

func (e *EventConnLoop) deleteConn(conn inter.Conn, isClose bool) {
	if isClose {
		e.eventConnFn.OnClose(conn)
	}
	conn.Close()
}

func (e *EventConnLoop) Close() {
	close(e.evCh)
}

func (e *EventConnLoop) PutConnMsgEvent(ev *event.ConnMsgEvent) {
	event.ConnMsgEventPool.Put(ev)
}

func (e *EventConnLoop) PutErrorConnEvent(ev *event.ErrConnEvent) {
	event.ErrConnEventPool.Put(ev)
}
