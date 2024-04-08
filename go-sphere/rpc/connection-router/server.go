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

	// Send Acknowledgement for time measurement - temp service
	cr.Socket.SendUDPMessage(2, "Receive data from C-R")

	cr.Socket.Broadcast(req.Data)

	return &spherePB.Response{
		Message: "Successfully Send Block Data",
	}, nil
}
