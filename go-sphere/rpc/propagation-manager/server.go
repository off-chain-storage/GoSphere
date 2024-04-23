package propamanager

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/socket"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

type Manager struct {
	Ctx    context.Context
	Socket socket.Socket
}

// 여기는 P-M이 C-R로부터 데이터를 수신하는 부분 - 여기서 청크로 받은 파일들 이어 붙여서 각 소켓으로 전파해야 함
func (pm *Manager) SendDataToPropagationManager(stream spherePB.PropagationManager_SendDataToPropagationManagerServer) error {
	log.Info("Received Block Data Request from Connection Router")

	var tempBlockData []byte
	for {
		req, err := stream.Recv()
		if err != nil {
			log.WithError(err).Error("Failed to receive chunk from C-R")
			return err
		}

		tempBlockData = append(tempBlockData, req.Data...)
		if req.IsLast {
			// Send Acknowledgement for time measurement - temp service
			go pm.Socket.SendUDPMessage(2, "Receive data from C-R")

			pm.Socket.Broadcast(tempBlockData)
			return stream.SendAndClose(&spherePB.Response{
				Message: "Successfully Send Block Data",
			})
		}
	}
}
