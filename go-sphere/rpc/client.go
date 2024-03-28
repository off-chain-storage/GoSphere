package rpc

import (
	"context"

	"github.com/off-chain-storage/GoSphere/go-sphere/db"
	grpcapi "github.com/off-chain-storage/GoSphere/go-sphere/rpc/grpc-api"
	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/helpers"
	"github.com/off-chain-storage/GoSphere/go-sphere/rpc/iface"
	"github.com/off-chain-storage/GoSphere/go-sphere/sync"
	"google.golang.org/grpc"
)

type ClientConfig struct {
	Endpoints  []string
	MaxMsgSize int
	DB         db.AccessRedisDB
}

type ClientService struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *ClientConfig
	router map[string]iface.Router
	conns  map[string]helpers.NodeConnection
}

func NewClient(ctx context.Context, cfg *ClientConfig) *ClientService {
	ctx, cancel := context.WithCancel(ctx)

	cs := &ClientService{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
		router: make(map[string]iface.Router),
		conns:  make(map[string]helpers.NodeConnection),
	}

	dialOpts := ConstructDialOptions(cs.cfg.MaxMsgSize)
	if dialOpts == nil {
		return cs
	}

	for _, endpoint := range cs.cfg.Endpoints {
		grpcConn, err := grpc.DialContext(ctx, endpoint, dialOpts...)
		if err != nil {
			log.WithError(err).Fatalln("Could not connect to gRPC server")
		}

		cs.conns[endpoint] = helpers.NewNodeConnection(grpcConn)

		cs.cfg.DB.Set(endpoint, grpcConn.GetState().String())

		log.Info("Connected to gRPC ", endpoint, " Server")
	}

	if len(cs.conns) == 0 {
		log.Error("No connections were established. Exiting...")
		return nil
	}

	return cs
}

func (cs *ClientService) Start() {
	for _, endpoint := range cs.cfg.Endpoints {
		routerClient := grpcapi.NewGrpcRouterClient(cs.conns[endpoint].GetGrpcClientConn())

		routerStruct := &router{
			routerClient: routerClient,
		}

		cs.router[endpoint] = routerStruct
		sync.SetRPCServerRouterInfo(endpoint, routerStruct)
	}
}

func (cs *ClientService) Stop() error {
	cs.cancel()
	log.Info("Stopping client service")
	for _, endpoint := range cs.cfg.Endpoints {
		if err := cs.conns[endpoint].GetGrpcClientConn().Close(); err != nil {
			return err
		}
	}
	return nil
}

func ConstructDialOptions(maxCallRecvMsgSize int, extraOpts ...grpc.DialOption) []grpc.DialOption {
	if maxCallRecvMsgSize == 0 {
		maxCallRecvMsgSize = 10 * 5 << 20 // 50MB
	}

	// grpc.WithTransportCredentials(insecure.NewCredentials())

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
		),
	}

	dialOpts = append(dialOpts, extraOpts...)
	return dialOpts
}
