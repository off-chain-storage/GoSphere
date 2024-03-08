package kafka

import "context"

type Kafka interface {
	StreamProvider
	Topic
}

type StreamProvider interface {
	Broadcast([]byte) error
}

type Topic interface {
	JoinTopic(topic string) (string, error)
	LeaveTopic(topic string) error
	PublishToTopic(ctx context.Context, topic string, data []byte) error
	SubscribeToTopic(topic string) (string, error)
}
