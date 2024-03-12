package socket

import "github.com/gofiber/contrib/websocket"

/* Block Data Flow */
// 1. Send to Kafka
// 2. Connection-Router Consume Block Data
// 3. Connection-Router Send to each Propagation Manager by tcp(socket) or gRPC
// 4. Propagation Manager Broadcast to each Clients(WebSocket)

func (s *Service) run() {
	for {
		select {
		// Register a new client in Connection-Router
		case connection := <-s.register:
			s.clients[connection] = &Client{}
			log.Info("Connection Registered")

		// Unregister a client in Connection-Router
		case connection := <-s.unregister:
			delete(s.clients, connection)
			log.Info("Connection Unregistered")

		// Send to Kafka
		case message := <-s.broadcast:
			if err := s.cfg.Kafka.Produce(message); err != nil {
				log.Error("Kafka Produce Error:", err)
			}

		// Broadcast a message to each Clients(WebSocket)
		case broadcast := <-s.broadcast:
			for connection, c := range s.clients {
				go func(connection *websocket.Conn, c *Client) {
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}
					if err := connection.WriteMessage(websocket.BinaryMessage, broadcast); err != nil {
						c.isClosing = true
						log.Println("write error:", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						s.unregister <- connection
					}
				}(connection, c)
			}
			log.Debugln("Message Broadcast:", broadcast)
		}
	}
}
