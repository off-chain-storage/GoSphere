package socket

import (
	"context"
	"net"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/off-chain-storage/GoSphere/go-sphere/db"
	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
)

type Config struct {
	WsAddr string
	Router *fiber.App
	Kafka  kafka.StreamProvider
	DB     db.ReadOnlyRedisDB
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

	// UDP Server - temp service
	udpServer_1 *net.UDPAddr
	udpServer_2 *net.UDPAddr
	conn_1      *net.UDPConn
	conn_2      *net.UDPConn
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

	// UDP Server - temp service
	if err := s.buildUDPAddr(); err != nil {
		log.WithError(err).Fatal("Could not build UDP address")
	}

	if err := s.Conn(); err != nil {
		log.WithError(err).Fatal("Could not start UDP listener")
	}

	return s, nil
}

func (s *Service) buildUDPAddr() error {
	udpServer_1, err := net.ResolveUDPAddr("udp4", "43.200.145.206:30004")
	if err != nil {
		return err
	}
	s.udpServer_1 = udpServer_1

	udpServer_2, err := net.ResolveUDPAddr("udp4", "43.200.145.206:30006")
	if err != nil {
		return err
	}
	s.udpServer_2 = udpServer_2

	return nil
}

func (s *Service) Conn() error {
	conn_1, err := net.DialUDP("udp4", nil, s.udpServer_1)
	if err != nil {
		return err
	}
	s.conn_1 = conn_1

	conn_2, err := net.DialUDP("udp4", nil, s.udpServer_2)
	if err != nil {
		return err
	}
	s.conn_2 = conn_2

	return nil
}

func (s *Service) SendUDPMessage(version int, msg string) error {
	if version == 1 {
		_, err := s.conn_1.Write([]byte(msg + "\n"))
		if err != nil {
			return err
		}
	} else {
		_, err := s.conn_2.Write([]byte(msg + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
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
