package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

func (s *Service) buildProducerOptions() {
	options := sarama.NewConfig()
	options.Producer.RequiredAcks = sarama.WaitForAll
	options.Producer.Retry.Max = 5
	options.Producer.Retry.Backoff = 100 * time.Millisecond
	options.Producer.Return.Successes = true
	options.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(s.cfg.BrokerList, options)
	if err != nil {
		log.WithError(err).Error("Failed to create Kafka producer")
		return
	}

	s.producer = producer
}

func (s *Service) buildConsumerGroupOptions() {
	options := sarama.NewConfig()
	options.Consumer.Offsets.Initial = sarama.OffsetOldest
	options.Consumer.Return.Errors = true
	options.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	options.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(s.cfg.BrokerList, "go-sphere", options)
	if err != nil {
		log.WithError(err).Error("Failed to create Kafka consumer group")
		return
	}

	s.consumer = consumer
}
