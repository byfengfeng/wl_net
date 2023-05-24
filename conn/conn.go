package conn

import (
	"fmt"
	"github.com/Byfengfeng/wl_net/enum"
	"github.com/Byfengfeng/wl_net/event"
	"github.com/Byfengfeng/wl_net/inter"
	"github.com/Byfengfeng/wl_net/log"
	"github.com/Byfengfeng/wl_net/pool"
	"github.com/Byfengfeng/wl_net/snowflake"
	"go.uber.org/zap"
	"io"
	"net"
)

type Conn struct {
	ch            chan<- any
	conn          net.Conn
	writeCh       chan []byte
	oCh           <-chan struct{}
	closeReadChan chan struct{}
	connType      enum.ConnType
	network       enum.NetType
	id            int64
	inter.Codec
}

var (
	connPool = pool.NewPool(func() any {
		return &Conn{}
	})
)

func NewConn(ch chan<- any, oCh <-chan struct{}, conn net.Conn, codec inter.Codec, network enum.NetType) inter.Conn {
	connect := connPool.Get().(*Conn)
	connect.conn = conn
	connect.ch = ch
	connect.id = snowflake.GenSnowflakeRegionNodeId()
	connect.oCh = oCh
	connect.writeCh = make(chan []byte)
	connect.closeReadChan = make(chan struct{})
	connect.connType = enum.Ready
	connect.network = network
	connect.Codec = codec
	return connect
}

func (c *Conn) Listen() {
	go c.write()
	switch c.network {
	case enum.Tcp:
		go c.read()
	case enum.WebSocket:
		c.read()
	case enum.Udp:
		go c.readUdp()
	}

}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) read() {
	headLen := make([]byte, enum.HeadSize)
	for {
		select {
		case <-c.closeReadChan:
			return
		case <-c.oCh:
			return
		default:
			_, err := io.ReadFull(c.conn, headLen)
			if err != nil {
				c.Action(event.NewErrConnEvent(c, true))
				return
			}
			length, data := c.Decode(headLen)
			if length > uint32(enum.HeadSize) {
				_, err = io.ReadFull(c.conn, data)
				if err != nil {
					c.Action(event.NewErrConnEvent(c, true))
					return
				}
				c.Action(event.NewConnMsgEvent(c, data))
			}
		}

	}
}

func (c *Conn) readUdp() {
	headLen := make([]byte, 65536)
	for {
		select {
		case <-c.closeReadChan:
			return
		case <-c.oCh:
			return
		default:
			n, err := c.conn.Read(headLen)
			if err != nil {
				c.Action(event.NewErrConnEvent(c, true))
				return
			}
			length := c.DecodeLength(headLen)
			if length > uint32(enum.HeadSize) {
				c.Action(event.NewConnMsgEvent(c, headLen[enum.HeadSize:n]))
			}
		}
	}
}

func (c *Conn) write() {
	for {
		select {
		case <-c.oCh:
			c.Close()
			return
		case <-c.closeReadChan:
			return
		case data := <-c.writeCh:
			_, err := c.conn.Write(c.Encode(data))
			if err != nil {
				c.Action(event.NewErrConnEvent(c, true))
			}
		}
	}
}

func (c *Conn) Close() {
	c.connType = enum.Close
	err := c.conn.Close()
	if err != nil {
		log.Error(fmt.Sprintf("%s channel close fail", c.conn.RemoteAddr()), zap.Error(err))
	}

	close(c.closeReadChan)
	close(c.writeCh)
	c.conn.Close()
	connPool.Put(c)
}

func (c *Conn) AsyncWrite(bytes []byte) {
	if c.connType == enum.Ready {
		c.writeCh <- bytes
	}
}

func (c *Conn) ConnId() int64 {
	return c.id
}

func (c *Conn) Action(event any) {
	if c.connType != enum.Close {
		c.ch <- event
	}
}
