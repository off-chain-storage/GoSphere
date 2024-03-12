package socket

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
)

type Config struct {
	WsAddr string
	Router *fiber.App
	Kafka  kafka.StreamProvider
}

type Service struct {
	started    bool
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *Config
	router     *fiber.App
	clients    map[*websocket.Conn]*Client
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	propose    chan []byte
	broadcast  chan []byte
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel

	s := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
		router: cfg.Router,
	}

	// Initialize the message channel
	s.InitMessageChannel()

	// Register the websocket handler
	s.InitRouter()

	return s, nil
}

func (s *Service) Start() {
	if s.started {
		return
	}

	addr := s.cfg.WsAddr

	// Start the Web Server
	go s.router.Listen(addr)

	log.WithField("Address", addr).Info("http listening on address")

	s.started = true
}

func (s *Service) Stop() error {
	defer s.cancel()
	s.started = false

	return nil
}
