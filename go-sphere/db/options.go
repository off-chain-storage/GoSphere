package db

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) buildOptions() (options *redis.Options) {
	cfg := s.cfg

	cfg.PoolFIFO = true

	cfg.Dialer = func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, 5*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	options = &redis.Options{
		Addr:     cfg.DbAddr,
		PoolFIFO: cfg.PoolFIFO,
		Dialer:   cfg.Dialer,
	}

	return
}
