package kafka

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Config struct {
	BrokerList []string
}

type Service struct {
	started          bool
	ctx              context.Context
	cancel           context.CancelFunc
	cfg              *Config
	producer         sarama.AsyncProducer
	consumer         sarama.ConsumerGroup
	joinedTopics     map[string]string
	joinedTopicsLock sync.RWMutex
}

func NewKafkaService(ctx context.Context, cfg *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel

	s := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
	}

	// Build Sarama Async Producer
	s.buildProducerOptions()
	s.buildConsumerGroupOptions()

	return s, nil
}

func (s *Service) Start() {
	if s.started {
		return
	}

	s.started = true
}

func (s *Service) Stop() error {
	defer s.cancel()
	s.started = false

	return nil
}
