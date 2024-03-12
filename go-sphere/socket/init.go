package socket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/pkg/errors"
)

type Client struct {
	isClosing bool
	mu        sync.Mutex
}

func (s *Service) InitMessageChannel() {
	s.clients = make(map[*websocket.Conn]*Client)
	s.register = make(chan *websocket.Conn)
	s.unregister = make(chan *websocket.Conn)
	s.broadcast = make(chan []byte)
	s.propose = make(chan []byte)
}

func (s *Service) InitRouter() error {
	// Register Proposer Web Server's Router
	if s.router == nil {
		return errors.New("no fiber router on server")
	}

	// Start the websocket msg handler
	go s.run()

	s.router.Get("/ws", websocket.New(s.wsHandler))

	log.Info("Initialize Proposer REST API Routes")

	return nil
}
