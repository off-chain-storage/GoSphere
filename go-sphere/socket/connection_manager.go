package socket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	isClosing bool
	mu        sync.Mutex
}

func (s *Service) RegisterClient(connection *websocket.Conn) {
	s.clients[connection] = &Client{}
}

func (s *Service) UnRegisterClient(connection *websocket.Conn) {
	delete(s.clients, connection)
}

func (s *Service) GetClientsInfo() {

}
