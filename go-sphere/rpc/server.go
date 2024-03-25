package rpc

import (
	"context"
	"net"

	connrouter "github.com/off-chain-storage/GoSphere/go-sphere/rpc/connection-router"
	"github.com/off-chain-storage/GoSphere/go-sphere/socket"
	spherePB "github.com/off-chain-storage/GoSphere/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Addr       string
	MaxMsgSize int
	Socket     socket.Socket
}

type ServerService struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *ServerConfig
	listener   net.Listener
	grpcServer *grpc.Server
	grpcClient map[net.Addr]bool
}

func NewServer(ctx context.Context, cfg *ServerConfig) *ServerService {
	ctx, cancel := context.WithCancel(ctx)

	rs := &ServerService{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		grpcClient: make(map[net.Addr]bool),
	}

	// gRPC Server Address
	address := cfg.Addr
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.WithError(err).Fatalln("Could not listen to port in Start()", address)
	}

	rs.listener = lis
	log.WithField("Address", address).Info("gRPC server listening on port")

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.MaxMsgSize),
	}

	rs.grpcServer = grpc.NewServer(opts...)

	return rs
}

func (rs *ServerService) Start() {
	connectionRouter := &connrouter.Router{
		Ctx:    rs.ctx,
		Socket: rs.cfg.Socket,
	}

	spherePB.RegisterConnectionRouterServer(rs.grpcServer, connectionRouter)

	reflection.Register(rs.grpcServer)

	go func() {
		if rs.listener != nil {
			if err := rs.grpcServer.Serve(rs.listener); err != nil {
				log.WithError(err).Error("gRPC server failed to serve")
			}
		}
	}()
}

func (rs *ServerService) Stop() error {
	rs.cancel()
	if rs.listener != nil {
		rs.grpcServer.GracefulStop()
		log.Debugln("gRPC server stopped")
	}
	return nil
}
