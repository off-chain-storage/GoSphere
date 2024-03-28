package iface

import (
	"context"

	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

type RouterClient interface {
	SendDataToPropagationManager(ctx context.Context, in *spherePB.BlockData) (*spherePB.Response, error)
}
