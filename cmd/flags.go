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

	/* gRPC Flag */
	// RPCAddrFlag defines the address of the gRPC server.
	RPCAddrFlag = &cli.StringFlag{
		Name:  "grpc-server-addr",
		Usage: "gRPC 서버 주소 (\"localhost:8080\")",
	}

	// EndPointFlag defines the address of the gRPC server for Client.
	EndPoint = &cli.StringFlag{
		Name:  "endpoint",
		Usage: "gRPC 서버 주소 (쉼표로 구분, 예: \"localhost:9001,localhost:9002,localhost:9003\")",
	}

	GrpcMaxCallRecvMsgSizeFlag = &cli.IntFlag{
		Name:  "grpc-max-msg-size",
		Usage: "gRPC 서버의 최대 수신 메시지 크기 (bytes)",
		Value: 1024 * 1024 * 1024, // 1GB
	}

	/* Websocket Flag */
	// WebsocketPortFlag defines the address of the websocket server.
	WebsocketAddrFlag = &cli.StringFlag{
		Name:  "websocket-addr",
		Usage: "Websocket 서버 주소 (\"localhost:8080\")",
	}

	/* Kafka Flag */
	// KafkaBrokersFlag defines the address of the Kafka brokers.
	KafkaBrokersFlag = &cli.StringFlag{
		Name:  "kafka-brokers",
		Usage: "Kafka 브로커 주소 (쉼표로 구분, 예: \"localhost:9001,localhost:9002,localhost:9003\")",
	}
)
