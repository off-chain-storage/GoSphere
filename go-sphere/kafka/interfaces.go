package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type Kafka interface {
	StreamProvider
	TopicProvider
	Producer
	Consumer
}

type StreamProvider interface {
	Broadcast([]byte) error
}

type TopicProvider interface {
	JoinTopic(topic string) (string, error)
	LeaveTopic(topic string) error
	PublishToTopic(ctx context.Context, topic string, data []byte) error
	SubscribeToTopic(topic string) (string, error)
}

type Producer interface {
	Producer() sarama.AsyncProducer
}

type Consumer interface {
	Consumer() sarama.ConsumerGroup
}
