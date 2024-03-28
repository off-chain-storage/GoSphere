package grpcapi

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/iface"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
	"google.golang.org/grpc"
)

type grpcRouterClient struct {
	routerClient spherePB.ConnectionRouterClient
}

func (c *grpcRouterClient) SendDataToPropagationManager(ctx context.Context, in *spherePB.BlockData) (*spherePB.Response, error) {
	return c.routerClient.SendDataToPropagationManager(ctx, in)
}

func NewGrpcRouterClient(cc grpc.ClientConnInterface) iface.RouterClient {
	return &grpcRouterClient{spherePB.NewConnectionRouterClient(cc)}
}
