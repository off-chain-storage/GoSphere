package iface

import "context"

type ManagerServer interface {
	SendDataToPropagationManager(ctx context.Context, blockData []byte) error
}
