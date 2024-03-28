package sdk

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// env Config Structure
type Config struct {
	Path string `mapstructure:"PM_ADDR"`
}

/* Variables about */
// Propagation Manager
var err error
var conn *websocket.Conn

// Message Write & Read Channel
var msgSenderChan chan []byte
var msgReceiverChan chan []byte

var wg sync.WaitGroup

func SetupPropagationModule() {
	log.Info("Setting up Propagation Module...")

	// Load Config
	c, err := loadConfig()
	if err != nil {
		log.Error("failed to load config: ", err)
		return
	}

	// Set up Propagation Module
	pmAddr := c.Path
	msgSenderChan = make(chan []byte)
	msgReceiverChan = make(chan []byte)

	startPropagationModule(pmAddr)
}

func startPropagationModule(pmAddr string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dial to Propagation Module
	dialer := websocket.Dialer{
		ReadBufferSize:  1024 * 1024 * 1.5,
		WriteBufferSize: 1024 * 1024 * 1.5,
	}

	conn, _, err = dialer.Dial(pmAddr, nil)
	if err != nil {
		log.Error("failed to dial to propagation module: ", err)
		return
	}

	conn.SetReadLimit(1024 * 1024 * 2)

	// Start the message sending loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		messageLoop(ctx)
	}()

	wg.Wait()
	conn.Close()
}

func messageLoop(ctx context.Context) {
	for {
		select {
		case data := <-msgSenderChan:
			if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Errorf("Failed to write message to propagation module: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// Write Message
func WriteMessage(msg []byte) {
	msgSenderChan <- msg
}

// Read Message
func ReadMessage(ctx context.Context) chan []byte {
	go readMessage(ctx)

	return msgReceiverChan
}

func readMessage(cxt context.Context) {
	for {
		select {
		case <-cxt.Done():
			return

		default:
			msgType, msgData, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage error:", err)
				return
			}

			log.Info("Message Received")

			if msgType == websocket.BinaryMessage {
				msgReceiverChan <- msgData
			} else {
				log.Error("Received message is not validated")
			}
		}
	}
}

func loadConfig() (c *Config, err error) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		log.Error("failed to read config file: ", err)
		return nil, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Error("failed to unmarshal config: ", err)
		return nil, err
	}

	return c, nil
}
