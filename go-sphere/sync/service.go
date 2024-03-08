package sync

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
)

type config struct {
	kafka kafka.Kafka
}

type Service struct {
	cfg                 *config
	ctx                 context.Context
	cancel              context.CancelFunc
	initialSyncComplete chan struct{}
}

func NewService(ctx context.Context, opts ...Option) *Service {
	ctx, cancel := context.WithCancel(ctx)
	r := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    &config{},
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil
		}
	}

	return r
}

func (s *Service) Start() {
	log.Info("Start Sync Service")

}

func (s *Service) Stop() error {
	s.cancel()

	return nil
}
