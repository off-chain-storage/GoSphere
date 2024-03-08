package sync

import "github.com/off-chain-storage/GoSphere/go-sphere/kafka"

type Option func(s *Service) error

func WithKafka(kafka kafka.Kafka) Option {
	return func(s *Service) error {
		s.cfg.kafka = kafka
		return nil
	}
}

func WithInitialSyncComplete(c chan struct{}) Option {
	return func(s *Service) error {
		s.initialSyncComplete = c
		return nil
	}
}
