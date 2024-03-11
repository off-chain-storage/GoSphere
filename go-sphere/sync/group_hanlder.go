package sync

import (
	"github.com/IBM/sarama"
)

var consumer = &Consumer{
	ready: make(chan bool),
}

type Consumer struct {
	ready chan bool
}

func GetConsumerHandler() *Consumer {
	return consumer
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Error("Message channel was closed")
				return nil
			}
			log.Info("Message claimed: ", string(message.Value))
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
