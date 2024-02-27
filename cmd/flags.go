package cmd

import "github.com/urfave/cli/v2"

var (
	/* CommonFlag */
	// VerbosityFlag defines the logrus configuration.
	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "Logging Level (trace, debug, info=default, warn, error, fatal, panic)",
		Value: "info",
	}

	/* RedisDB Flag */
	// RedisDBAddrFlag defines the address of the RedisDB.
	RedisDBAddrFlag = &cli.StringFlag{
		Name:  "redis-addr",
		Usage: "Connection 라우팅 데이터베이스 (Address of the RedisDB)",
	}

	/* Conn-Router Flag */
	// EnableConnRouterFlag defines the flag to enable the node as a connection router.
	EnableConnRouterFlag = &cli.BoolFlag{
		Name:  "enable-conn-router",
		Usage: "Connection-Router Node로 빌드 (Enable the node as a connection router)",
		Value: false,
	}
)
