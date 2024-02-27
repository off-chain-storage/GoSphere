package node

import (
	"context"

	"github.com/off-chain-storage/GoSphere/cmd"
	"github.com/off-chain-storage/GoSphere/go-sphere/db"
	"github.com/off-chain-storage/GoSphere/runtime"
	"github.com/urfave/cli/v2"
)

type GoSphereNode struct {
	cliCtx   *cli.Context
	ctx      context.Context
	cancel   context.CancelFunc
	services *runtime.ServiceRegistry
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

	/* Register Services */
	// Register Redis DB for propagation manager routing
	log.Debugln("Starting Redis DB")
	if err := goSphere.startRedisDB(cliCtx); err != nil {
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
