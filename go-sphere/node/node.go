package node

import (
	"context"
	"strings"

	"github.com/off-chain-storage/GoSphere/cmd"
	"github.com/off-chain-storage/GoSphere/go-sphere/db"
	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
	regularSync "github.com/off-chain-storage/GoSphere/go-sphere/sync"
	"github.com/off-chain-storage/GoSphere/runtime"
	"github.com/urfave/cli/v2"
)

type GoSphereNode struct {
	cliCtx              *cli.Context
	ctx                 context.Context
	cancel              context.CancelFunc
	services            *runtime.ServiceRegistry
	initialSyncComplete chan struct{}
}

func New(cliCtx *cli.Context, cancel context.CancelFunc) (*GoSphereNode, error) {
	// Reflection (Service Runtime)
	registry := runtime.NewServiceRegistry()

	ctx := cliCtx.Context
	goSphere := &GoSphereNode{
		cliCtx:   cliCtx,
		ctx:      ctx,
		cancel:   cancel,
		services: registry,
	}

	goSphere.initialSyncComplete = make(chan struct{})

	/* Register Services */
	// Register Redis DB for propagation manager routing
	log.Debugln("Starting Redis DB")
	if err := goSphere.startRedisDB(cliCtx); err != nil {
		return nil, err
	}

	// Register Kafka for message broker
	log.Debugln("Starting Kafka")
	if err := goSphere.startKafka(cliCtx); err != nil {
		return nil, err
	}

	// Register Sync Service for Syncing
	log.Debugln("Starting Sync Service")
	if err := goSphere.registerSyncService(goSphere.initialSyncComplete); err != nil {
		return nil, err
	}

	return goSphere, nil
}

func (g *GoSphereNode) Start() {

}

func (g *GoSphereNode) Close() {

}

func (g *GoSphereNode) startRedisDB(cliCtx *cli.Context) error {
	dbAddr := cliCtx.String(cmd.RedisDBAddrFlag.Name)

	svc, err := db.NewRedisClient(g.ctx, &db.Config{
		DbAddr: dbAddr,
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect Redis DB")
		return err
	}

	log.Info("Connecting to Redis DB")
	return g.services.RegisterService(svc)
}

func (g *GoSphereNode) startKafka(cliCtx *cli.Context) error {
	kafkaBrokers := cliCtx.String(cmd.KafkaBrokersFlag.Name)
	brokerList := strings.Split(kafkaBrokers, ",")

	svc, err := kafka.NewKafkaService(g.ctx, &kafka.Config{
		BrokerList: brokerList,
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect Kafka")
		return err
	}

	log.Info("Connecting to Kafka")
	return g.services.RegisterService(svc)
}

func (g *GoSphereNode) registerSyncService(initialSyncComplete chan struct{}) error {
	svc := regularSync.NewService(
		g.ctx,
		regularSync.WithInitialSyncComplete(initialSyncComplete),
	)

	log.Info("Registering Sync Service")

	return g.services.RegisterService(svc)
}

func (g *GoSphereNode) fetchKafka() kafka.Kafka {
	var k *kafka.Service
	if err := g.services.FetchService(&k); err != nil {
		panic(err)
	}
	return k
}
