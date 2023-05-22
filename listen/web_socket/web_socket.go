package web_socket

import (
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"wl_net/log"
)

type WebSocketListen struct {
	address       string
	channelHandel func(conn net.Conn)
}

func NewWebSocketListen(addr string, channelHandel func(conn net.Conn)) *WebSocketListen {
	return &WebSocketListen{address: addr, channelHandel: channelHandel}
}

func (w *WebSocketListen) Start() {
	handler := websocket.Handler(func(conn *websocket.Conn) {
		w.channelHandel(conn)
	})
	http.Handle("/", handler)
	go func() {
		log.Info("websocket listen start success")
		err := http.ListenAndServe(w.address, nil)
		if err != nil {
			log.Error("http serve listen exit ", zap.Any("err", err.Error()))
			return
		}
	}()
}
