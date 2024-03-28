package connrouter

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/socket"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

type Router struct {
	Ctx    context.Context
	Socket socket.Socket
}

func (cr *Router) SendDataToPropagationManager(ctx context.Context, req *spherePB.BlockData) (*spherePB.Response, error) {
	log.Info("Received Block Data Request from Connection Router")

	// 여기에 데이터 자기가 물고 있는 Propagation Manager로 전파하는 코드 추가
	cr.Socket.Broadcast(req.Data)

	return &spherePB.Response{
		Message: "Successfully Send Block Data",
	}, nil
}
