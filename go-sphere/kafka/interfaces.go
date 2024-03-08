package kafka

type Kafka interface {
	StreamProvider
	Topic
}

type StreamProvider interface {
	Broadcast([]byte) error
}

type Topic interface {
	JoinTopic(topic string) error
	LeaveTopic(topic string) error
	PublishToTopic(topic string, data []byte) error
	SubscribeToTopic(topic string) error
}
