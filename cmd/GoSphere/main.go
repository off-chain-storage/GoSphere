package main

import (
	"context"
	"fmt"
	"os"

	"github.com/off-chain-storage/GoSphere/cmd"
	"github.com/off-chain-storage/GoSphere/go-sphere/node"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var appFlags = []cli.Flag{
	/* Common Flag */
	cmd.VerbosityFlag,

	/* RedisDB Flag */
	cmd.RedisDBAddrFlag,

	/* Conn-Router Flag */
	cmd.EnableConnRouterFlag,

	/* Websocket Flag */
	cmd.WebsocketAddrFlag,

	/* Kafka Flag */
	cmd.KafkaBrokersFlag,
	cmd.KafkaGroupIDFlag,

	/* gRPC Flag */
	cmd.RPCAddrFlag,
	cmd.EndPoint,
	cmd.GrpcMaxCallRecvMsgSizeFlag,
}

func main() {
	rCtx, cancel := context.WithCancel(context.Background())

	app := cli.App{}
	app.Name = "GoSphere"
	app.Usage = "this is a GoSphere implementation for Propagation Hub SDK"
	app.Action = func(ctx *cli.Context) error {
		if err := startSDK(ctx, cancel); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}

	app.Flags = appFlags

	if err := app.RunContext(rCtx, os.Args); err != nil {
		log.Error(rCtx, err.Error())
	}
}

func startSDK(ctx *cli.Context, cancel context.CancelFunc) error {
	verbosity := ctx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}

	// Set log level
	logrus.SetLevel(level)

	// Start SDK
	goSphere, err := node.New(ctx, cancel)
	if err != nil {
		return fmt.Errorf("unable to start GoSphere: %w", err)
	}

	goSphere.Start()
	return nil
}
