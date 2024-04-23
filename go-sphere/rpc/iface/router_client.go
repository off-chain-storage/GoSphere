package iface

import (
	"context"

	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

type RouterClient interface {
	SendDataToPropagationManager(ctx context.Context) (spherePB.PropagationManager_SendDataToPropagationManagerClient, error)
}
