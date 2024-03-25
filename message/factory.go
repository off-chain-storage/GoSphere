package message

import (
	spherePB "github.com/off-chain-storage/GoSphere/proto"
)

func BuildData(msg []byte) (*spherePB.BlockData, error) {
	BlockData := &spherePB.BlockData{
		Data: msg,
	}

	return BlockData, nil
}
