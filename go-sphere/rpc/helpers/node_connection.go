package helpers

import "google.golang.org/grpc"

type NodeConnection interface {
	GetGrpcClientConn() *grpc.ClientConn
}

type nodeConnection struct {
	grpcClientConn *grpc.ClientConn
}

func (c *nodeConnection) GetGrpcClientConn() *grpc.ClientConn {
	return c.grpcClientConn
}

func NewNodeConnection(grpcConn *grpc.ClientConn) NodeConnection {
	conn := &nodeConnection{}
	conn.grpcClientConn = grpcConn
	return conn
}
