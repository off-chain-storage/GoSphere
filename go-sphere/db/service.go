package db

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	DbAddr   string
	PoolFIFO bool
	Dialer   func(ctx context.Context, network, address string) (net.Conn, error)
}

type Service struct {
	ctx         context.Context
	cancel      context.CancelFunc
	cfg         *Config
	redisClient *redis.Client
	conn        *redis.Conn
}

func NewRedisClient(ctx context.Context, cfg *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel

	s := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
	}

	opts := s.buildOptions()
	r := redis.NewClient(opts)

	s.redisClient = r

	return s, nil
}

func (s *Service) Start() {

}

func (s *Service) Stop() error {
	defer s.cancel()
	defer s.conn.Close()

	return nil
}

func (s *Service) RedisClient() redis.Client { return *s.redisClient }
