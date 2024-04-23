package sync

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/iface"
)

var consumer = &Consumer{
	ready: make(chan bool),
}

type Consumer struct {
	ready   chan bool
	manager map[string]iface.ManagerServer
}

func init() {
	consumer.manager = make(map[string]iface.ManagerServer)
}

func SetRPCServerRouterInfo(endpoint string, newManager iface.ManagerServer) {
	consumer.manager[endpoint] = newManager
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
			log.Info("Received message")
			session.MarkMessage(message, "")

			// rpc/router.go 안에서
			// SendDataToPropagationManager() 함수를
			// 현재 gRPC Conn 유지되어있는 Router로 브로드캐스팅하는 코드 추가
			for _, value := range consumer.manager {
				go func(v iface.ManagerServer) {
					v.SendDataToPropagationManager(context.Background(), message.Value)
				}(value)
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
