package sync

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
	"github.com/pkg/errors"
)

func (s *Service) registerConsumerGroup() {
	s.consume(
		s.ctx,
		kafka.OriginalBlockTopicFormat,
		GetConsumerHandler(),
	)
}

func (s *Service) consume(ctx context.Context, topic string, handler *Consumer) sarama.ConsumerGroup {
	topics := []string{topic}

	return s.consumeWithBase(ctx, topics, handler)
}

func (s *Service) consumeWithBase(ctx context.Context, topics []string, handler *Consumer) sarama.ConsumerGroup {
	client := s.cfg.kafka.Consumer()

	messageLoop := func() {
		for {
			if err := client.Consume(ctx, topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			GetConsumerHandler().ready = make(chan bool)
		}
	}

	// within go routine 문제가 발생할 수 있음 - 일단 문제 X
	go messageLoop()

	log.WithField("topics", topics).Info("Consumer started")
	return client
}
