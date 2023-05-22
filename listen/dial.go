package listen

import (
	"fmt"
	"github.com/Byfengfeng/wl_net/enum"
	"github.com/Byfengfeng/wl_net/log"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
	"net"
	"time"
)

func Dial(netType enum.NetType, addr string) (conn net.Conn) {
	switch netType {
	case enum.Tcp:
		conn = tcpOrUdp("tcp", addr, reConn)
	case enum.Udp:
		conn = tcpOrUdp("udp", addr, reConn)
	case enum.WebSocket:
		conn = tcpOrUdp("websocket", addr, reConnWs)
	}

	return conn
}

const maxRetryCount = 10

func tcpOrUdp(network, addr string, connectFun func(network, addr string) net.Conn) (conn net.Conn) {
	retryCount := 1
	for {
		conn = connectFun(network, addr)
		if conn == nil {
			time.Sleep(time.Duration(retryCount) * time.Second)
			if retryCount < maxRetryCount {
				retryCount <<= 1
			} else if retryCount > maxRetryCount {
				retryCount = maxRetryCount
			}
			continue
		}
		return conn
	}
}

func reConn(network, addr string) net.Conn {
	conn, err := net.Dial(network, addr)
	if err != nil {
		log.Warn(fmt.Sprintf("%s channel dial err, In progress reConn %s", network, addr), zap.Error(err))
		return nil
	}
	return conn
}

func reConnWs(network, addr string) net.Conn {
	var url = fmt.Sprintf("ws://%s/ws", addr)
	var origin = fmt.Sprintf("http://%s", addr)
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Warn(fmt.Sprintf("websocket channel dial err, In progress reConn %s", addr), zap.Error(err))
		return nil
	}
	return conn
}
