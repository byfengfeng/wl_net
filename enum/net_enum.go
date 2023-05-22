package enum

type ConnType uint8

const (
	Ready ConnType = iota
	Close
)

type NetType uint8 //通信监听类型

const (
	Tcp NetType = iota + 1
	Udp
	WebSocket
	Http
)

type ConnHandlerType uint8 //通信链接处理类型

const (
	Dial ConnHandlerType = iota + 1
)

type ErrSourceType uint8

const (
	ErrRead ErrSourceType = iota + 1
	ErrWrite
	ErrServerOut
	ErrTickerTimeOut
)

type HttpType uint8

const (
	GET HttpType = iota + 1
	POST
	PUT
	DELETE
)

var HeadSize = 4
