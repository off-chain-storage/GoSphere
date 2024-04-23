package grpcapi

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/iface"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
	"google.golang.org/grpc"
)

type grpcRouterClient struct {
	routerClient spherePB.PropagationManagerClient
}

func (c *grpcRouterClient) SendDataToPropagationManager(ctx context.Context) (spherePB.PropagationManager_SendDataToPropagationManagerClient, error) {
	return c.routerClient.SendDataToPropagationManager(ctx)
}

func NewGrpcRouterClient(cc grpc.ClientConnInterface) iface.RouterClient {
	return &grpcRouterClient{spherePB.NewPropagationManagerClient(cc)}
}
