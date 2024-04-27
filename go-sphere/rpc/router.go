package rpc

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/iface"
	"github.com/off-chain-storage/GoSphere/message"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

type router struct {
	client        iface.RouterClient
	clientService *ClientService
}

// C-R이 P-M으로 데이터를 전송하는 함수
func (r *router) SendDataToPropagationManager(ctx context.Context, blockData []byte) error {
	blk, err := message.BuildData(blockData)
	if err != nil {
		log.WithError(err).Error("Failed to build block data")
		return err
	}

	// Send Acknowledgement for time measurement - temp service
	go r.clientService.SendUDPMessage("Data sent to P-M from C-R")

	stream, err := r.client.SendDataToPropagationManager(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to send data to propagation manager")
		return err
	}

	chunkSize := 1024 * 512
	chunks := chunkByteSlice(blk.Data, chunkSize)

	tempBlock := spherePB.BlockData{}
	for i, chunk := range chunks {
		tempBlock.Data = chunk
		if i == len(chunks)-1 {
			tempBlock.IsLast = true
		} else {
			tempBlock.IsLast = false
		}

		if err := stream.Send(&tempBlock); err != nil {
			log.WithError(err).Error("Failed to send block data to propagation manager")
			return err
		}
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		log.WithError(err).Error("Failed to close and receive stream")
		return err
	}

	log.Info("Data sent to propagation manager successfully")
	return nil
}

func chunkByteSlice(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for {
		if len(data) == 0 {
			break
		}

		if len(data) < chunkSize {
			chunks = append(chunks, data)
			break
		}

		chunks = append(chunks, data[:chunkSize])
		data = data[chunkSize:]
	}

	return chunks
}
