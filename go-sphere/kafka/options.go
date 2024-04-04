package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

func (s *Service) buildProducerOptions() {
	options := sarama.NewConfig()
	// options.Producer.RequiredAcks = sarama.WaitForAll
	options.Producer.Retry.Max = 5
	options.Producer.Retry.Backoff = 100 * time.Millisecond
	options.Producer.Return.Successes = true
	options.Producer.Return.Errors = true
	options.Producer.MaxMessageBytes = 1024 * 1024 * 10
	// options.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewSyncProducer(s.cfg.BrokerList, options)
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

	consumer, err := sarama.NewConsumerGroup(s.cfg.BrokerList, s.cfg.GroupID, options)
	if err != nil {
		log.WithError(err).Error("Failed to create Kafka consumer group")
		return
	}

	s.consumer = consumer
}
