package rpc

import (
	"context"
	"net"

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
	ctx     context.Context
	cancel  context.CancelFunc
	cfg     *ClientConfig
	manager map[string]iface.ManagerServer
	conns   map[string]helpers.NodeConnection

	// UDP Server - temp service
	udpServer *net.UDPAddr
	conn      *net.UDPConn
}

// NewClient creates a new gRPC client - Connection Router gRPC Client
func NewClient(ctx context.Context, cfg *ClientConfig) *ClientService {
	ctx, cancel := context.WithCancel(ctx)

	cs := &ClientService{
		ctx:     ctx,
		cancel:  cancel,
		cfg:     cfg,
		manager: make(map[string]iface.ManagerServer),
		conns:   make(map[string]helpers.NodeConnection),
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

	// UDP Server - temp service
	if err := cs.buildUDPAddr(); err != nil {
		log.WithError(err).Fatal("Could not build UDP address")
	}

	if err := cs.Conn(); err != nil {
		log.WithError(err).Fatal("Could not start UDP listener")
	}

	return cs
}

func (cs *ClientService) buildUDPAddr() error {
	udpServer, err := net.ResolveUDPAddr("udp4", "3.35.85.78:30006")
	if err != nil {
		return err
	}

	cs.udpServer = udpServer
	return nil
}

func (cs *ClientService) Conn() error {
	conn, err := net.DialUDP("udp4", nil, cs.udpServer)
	if err != nil {
		return err
	}
	cs.conn = conn
	return nil
}

func (cs *ClientService) SendUDPMessage(msg string) error {
	_, err := cs.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (cs *ClientService) Start() {
	for _, endpoint := range cs.cfg.Endpoints {
		crClient := grpcapi.NewGrpcRouterClient(cs.conns[endpoint].GetGrpcClientConn())

		managerStruct := &router{ // router means P-M gRPC client
			client:        crClient,
			clientService: cs,
		}

		cs.manager[endpoint] = managerStruct
		sync.SetRPCServerRouterInfo(endpoint, managerStruct)
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
