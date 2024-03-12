package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
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

	for {
		log.WithField("topic", topic).Debug("publishing message to topic")
		s.producer.Input() <- &sarama.ProducerMessage{
			Topic: topicHandle,
			Value: sarama.ByteEncoder(data),
		}

		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "unable to find requisite number of peers for topic")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
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
