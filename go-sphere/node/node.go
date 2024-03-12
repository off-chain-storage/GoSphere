package node

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/off-chain-storage/GoSphere/cmd"
	"github.com/off-chain-storage/GoSphere/go-sphere/db"
	"github.com/off-chain-storage/GoSphere/go-sphere/kafka"
	"github.com/off-chain-storage/GoSphere/go-sphere/socket"
	web "github.com/off-chain-storage/GoSphere/go-sphere/socket"
	regularSync "github.com/off-chain-storage/GoSphere/go-sphere/sync"
	"github.com/off-chain-storage/GoSphere/runtime"
	"github.com/urfave/cli/v2"
)

type GoSphereNode struct {
	cliCtx              *cli.Context
	ctx                 context.Context
	cancel              context.CancelFunc
	services            *runtime.ServiceRegistry
	lock                sync.RWMutex
	stop                chan struct{}
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
		stop:     make(chan struct{}),
		services: registry,
	}

	goSphere.initialSyncComplete = make(chan struct{})

	/* Register Services */
	// Register Redis DB for propagation manager routing
	log.Debugln("Starting Redis DB")
	if err := goSphere.registerRedisDB(cliCtx); err != nil {
		return nil, err
	}

	// Register Kafka for message broker
	log.Debugln("Starting Kafka")
	if err := goSphere.registerKafka(cliCtx); err != nil {
		return nil, err
	}

	// Register Web Service for Websocket
	log.Debugln("Starting WebSocket Service")
	if err := goSphere.registerWebSocketService(cliCtx); err != nil {
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
	g.lock.Lock()

	log.Info("Start GoSphere")
	g.services.StartAll()
	stop := g.stop

	g.lock.Unlock()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc

		log.Info("Got Interrupt, shutting down...")

		go g.Close()

		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.WithField("times", i-1).Info("Already shutting down, interrupt more to panic")
			}
		}
		panic("Panic closing the GoSphere")
	}()

	<-stop
}

func (g *GoSphereNode) Close() {
	g.lock.Lock()
	defer g.lock.Unlock()

	log.Info("Closing GoSphere")
	g.services.StopAll()
	g.cancel()

	close(g.stop)
}

func (g *GoSphereNode) registerRedisDB(cliCtx *cli.Context) error {
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

func (g *GoSphereNode) registerKafka(cliCtx *cli.Context) error {
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

func (g *GoSphereNode) registerWebSocketService(cliCtx *cli.Context) error {
	wsAddr := cliCtx.String(cmd.WebsocketAddrFlag.Name)

	svc, err := socket.NewService(g.ctx, &web.Config{
		WsAddr: wsAddr,
		Router: newRouter(),
		Kafka:  g.fetchKafka(),
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect Web Service")
		return err
	}

	log.Info("Registering Web Service")

	return g.services.RegisterService(svc)
}

func (g *GoSphereNode) registerSyncService(initialSyncComplete chan struct{}) error {
	svc := regularSync.NewService(
		g.ctx,
		regularSync.WithKafka(g.fetchKafka()),
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

func newRouter() *fiber.App {
	return fiber.New(
		fiber.Config{
			WriteBufferSize: int(1.5 * 1024 * 1024), // 1.5MB
			ReadBufferSize:  int(1.5 * 1024 * 1024), // 1.5MB
		},
	)
}
