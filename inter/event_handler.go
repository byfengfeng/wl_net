package inter

type EventConnHandler interface {
	OnOpened(conn Conn)
	OnClose(conn Conn)
	React(conn Conn, data []byte)
}

type EventDialHandler interface {
	OnDialOpened(conn Conn)
	OnDialClose(conn Conn)
	DialReact(conn Conn, data []byte)
}

