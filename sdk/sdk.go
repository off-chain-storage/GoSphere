package sdk

// func SetupPropagationModule(ctx context.Context) (*Subscription, error) {
// 	log.Info("Make Connection with Propagation Module...")

// 	// Load Config
// 	c, err := loadConfig()
// 	if err != nil {
// 		log.Error("failed to load config: ", err)
// 		return nil, err
// 	}

// 	dialer := websocket.Dialer{
// 		ReadBufferSize:  1024 * 1024 * 1.5,
// 		WriteBufferSize: 1024 * 1024 * 1.5,
// 	}

// 	conn, _, err := dialer.Dial(c.Path, nil)
// 	if err != nil {
// 		log.Error("failed to dial to propagation module: ", err)
// 		return nil, err
// 	}

// 	sub := &Subscription{
// 		conn: conn,
// 		ctx:  ctx,
// 	}

// 	if sub.ch == nil {
// 		sub.ch = make(chan *Message, 32)
// 	}

// 	out := make(chan *Subscription, 1)

// 	log.Info("Make Connection with Propagation Module... Done")
// 	return <-out, nil
// }

// func startPropagationModule(pmAddr string) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Dial to Propagation Module
// 	dialer = websocket.Dialer{
// 		ReadBufferSize:  1024 * 1024 * 1.5,
// 		WriteBufferSize: 1024 * 1024 * 1.5,
// 	}

// 	conn, _, err := dialer.Dial(pmAddr, nil)
// 	if err != nil {
// 		log.Error("failed to dial to propagation module: ", err)
// 		return
// 	}

// 	conn.SetReadLimit(1024 * 1024 * 2)

// 	// Start the message sending loop
// 	go func() {
// 		messageLoop(ctx, conn)
// 	}()

// 	log.Info("Connection with Propagation Module is established")
// }

// func messageLoop(ctx context.Context, conn *websocket.Conn) {
// 	for {
// 		select {
// 		// case <-ctx.Done():
// 		// 	return

// 		case data := <-msgSenderChan:
// 			log.Info("hi")
// 			if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
// 				log.Errorf("Failed to write message to propagation module: %v", err)
// 				log.Info("HI")
// 				return
// 			}
// 			log.Info("Message Sent")

// 		default:
// 			log.Info("Waiting for message to send...")
// 			msgType, msgData, err := conn.ReadMessage()
// 			if err != nil {
// 				log.Error("ReadMessage error:", err)
// 				return
// 			}

// 			log.Info("Message Received")

// 			if msgType == websocket.BinaryMessage {
// 				msgReceiverChan <- msgData
// 			} else {
// 				log.Error("Received message is not validated")
// 			}
// 		}
// 	}
// }

// // Write Message
// func WriteMessage(msg []byte) {
// 	msgSenderChan <- msg
// }

// // Read Message
// func ReadMessage(ctx context.Context) chan []byte {
// 	// go readMessage(ctx)

// 	return msgReceiverChan
// }

// func readMessage(cxt context.Context) {
// 	for {
// 		select {
// 		case <-cxt.Done():
// 			return

// 		default:
// 			msgType, msgData, err := conn.ReadMessage()
// 			if err != nil {
// 				log.Error("ReadMessage error:", err)
// 				return
// 			}

// 			log.Info("Message Received")

// 			if msgType == websocket.BinaryMessage {
// 				msgReceiverChan <- msgData
// 			} else {
// 				log.Error("Received message is not validated")
// 			}
// 		}
// 	}
// }
