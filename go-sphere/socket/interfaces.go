package socket

import "github.com/gofiber/contrib/websocket"

type Socket interface {
	Broadcaster
	ClientManager
	MsgProcessor
	ConnectionManager
}

type Broadcaster interface {
	Broadcast([]byte)
}

type ClientManager interface {
	GetClientsInfo()
}

type MsgProcessor interface {
	ProcessMsg([]byte) error
}

type ConnectionManager interface {
	RegisterClient(*websocket.Conn)
	UnRegisterClient(*websocket.Conn)
}
