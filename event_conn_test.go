package net

import (
	"encoding/binary"
	"fmt"
	"github.com/Byfengfeng/wl_net/enum"
	"github.com/Byfengfeng/wl_net/inter"
	"github.com/Byfengfeng/wl_net/log"
	"testing"
	"time"
	"unicode"
)

type Server struct {
}

func (s *Server) OnOpened(conn inter.Conn) {
	fmt.Println("connect")
}

func (s *Server) OnClose(conn inter.Conn) {
	fmt.Println("close")
}

var (
	handlerCount int64
)

func (s *Server) React(conn inter.Conn, data []byte) {
	handlerCount++
	conn.AsyncWrite([]byte("456"))
	fmt.Println(len(data), handlerCount)
}

func NewServer() inter.EventConnHandler {
	return &Server{}
}

type Client struct {
}

func (c Client) Async(data []byte) {
	//TODO implement me
	panic("implement me")
}

//260 306 80000 1分钟
//接收问题 发送问题 3G内存，内存问题

func (c Client) OnDialOpened(conn inter.Conn) {
	log.Info("conn")
}

func (c Client) OnDialClose(conn inter.Conn) {
	fmt.Println("close")
}

func (c Client) DialReact(conn inter.Conn, data []byte) {
	fmt.Println(string(data))
}

func NewClient() inter.EventDialHandler {
	return &Client{}
}

var (
	msgCount = int64(0)
)

type Err struct {
	code int32
	str  string
}

func testpanic() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic(Err{11, "测试"})
}

func TestGoto(t *testing.T) {
	currentTime := time.Now().UnixMilli()
	i := 5000
run:
	fmt.Println(i)
	i--

	if i > 0 {
		goto run
	}

	//for i := 0; i < 5000; i++ {
	//	fmt.Println(i)
	//}
	fmt.Println(time.Now().UnixMilli() - currentTime)
	time.Sleep(10 * time.Second)
}

func TestNewConnEventLoop(t *testing.T) {
	//eventLoop := NewConnEventLoop("", 9998, NewServer(), enum.Tcp)
	//eventLoop := NewConnEventLoop("", 9998, NewServer(), enum.Udp)
	//1W connect 5000
	//testpanic()
	//go func() {
	//	timer := time.NewTicker(1 * time.Second)
	//	for {
	//		select {
	//		case <-timer.C:
	//			if handlerCount > 0 {
	//				msgCount += 1
	//				log.Info(fmt.Sprintf("已处理%d个消息,平均每秒处理%d个消息", handlerCount, handlerCount/msgCount))
	//			}
	//		}
	//	}
	//}()
	eventLoop := NewConnEventLoop("", 9998, NewServer(), enum.WebSocket, 1)
	//eventLoop.WithDecodeLength(2)
	eventLoop.Run()
}

func TestNewEventDialLoop(t *testing.T) {
	//dialLoop := NewEventDialLoop("", 9998, enum.Tcp, NewClient())
	//dialLoop := NewEventDialLoop("", 9998, enum.Udp, NewClient())
	for i := 0; i < 5; i++ {
		dialLoop := NewEventDialLoop("", 9998, enum.Tcp, NewClient(), 2)
		go dialLoop.Run()
	}

	//time.Sleep(10 * time.Second)
	//log.Info("test start", zap.Int("time", time.Now().Second()))

	//log.Info("test end", zap.Int("time", time.Now().Second()))
	time.Sleep(10 * time.Minute)
}

func ChineseCount(str1 string) (count int) {
	for _, char := range str1 {
		if unicode.Is(unicode.Han, char) {
			count += 2
		} else {
			count++
		}
	}
	return
}

func TestNewEventDialLoop1(t *testing.T) {
	fmt.Println(ChineseCount("Hello你好WorldGo语言真强"))
	for j := 0; j < 10; j++ {
		go func() {
			dialLoop := NewEventDialLoop("", 9998, enum.WebSocket, NewClient(), 3)
			//dialLoop := NewEventDialLoop("", 9998, enum.Udp, NewClient())
			//dialLoop := NewEventDialLoop("", 9998, enum.WebSocket, NewClient())
			go dialLoop.Run()

			//dial, err := net.Dial("tcp", ":9998")
			//if err != nil {
			//	panic(err)
			//}
			time.Sleep(time.Second)
			for i := 0; i < 10000; i++ {
				time.Sleep(time.Nanosecond * 2)
				dialLoop.Send([]byte("123456789"))
				//fmt.Println(i)
			}
		}()

	}

	time.Sleep(10 * time.Minute)
}

func Encode(data []byte) []byte {
	length := uint32(len(data)) + uint32(enum.HeadSize)
	buf := make([]byte, length)
	binary.BigEndian.PutUint32(buf, length)
	buf = append(buf[:enum.HeadSize], data...)
	return buf
}

var str = ".jpg .jpeg .dds .psd .psd .psb .pdt .webp .xmp .gif .bmp .png .svg .tiff .tga .vda .icb .vst .eps .pns .obj .fbx .max .mb .ma .ai .3ds .chr .stl .blend .prefab  .meta  .scene"
