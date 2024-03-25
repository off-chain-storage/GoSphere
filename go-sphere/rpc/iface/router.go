package iface

import "context"

type Router interface {
	SendDataToPropagationManager(ctx context.Context, blockData []byte) error
}
