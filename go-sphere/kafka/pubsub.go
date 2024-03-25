package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

func (s *Service) JoinTopic(topic string) (string, error) {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	// Check if already joined the topic
	if _, ok := s.joinedTopics[topic]; !ok {
		s.joinedTopics[topic] = topic
	}

	return s.joinedTopics[topic], nil
}

func (s *Service) LeaveTopic(topic string) error {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	delete(s.joinedTopics, topic)

	return nil
}

func (s *Service) ProduceToTopic(ctx context.Context, topic string, data []byte) error {
	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return err
	}

	log.WithField("topic", topic).Debug("publishing message to topic")
	// s.producer.Input() <- &sarama.ProducerMessage{
	// 	Topic: topicHandle,
	// 	Key:   nil,
	// 	// Value: sarama.StringEncoder("HI"),
	// 	// 뭔가 보내는 메시지 용량과 관련있는 듯함...
	// 	Value: sarama.ByteEncoder(data),
	// }

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicHandle,
		Key:   nil,
		Value: sarama.ByteEncoder(data),
	})

	if err != nil {
		log.WithError(err).Error("failed to send message")
		return err
	}

	return nil
}

// 현재 미완
func (s *Service) SubscribeToTopic(topic string) (string, error) {
	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return "", err
	}

	log.WithField("topic", topic).Debug("subscribing to topic")

	return topicHandle, nil
}
