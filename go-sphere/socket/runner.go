package socket

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
			s.RegisterClient(connection)
			log.Info("Connection Registered")

		// Unregister a client in Connection-Router
		case connection := <-s.unregister:
			s.UnRegisterClient(connection)
			log.Info("Connection Unregistered")

		// Send to Kafka
		case message := <-s.propose:
			if err := s.ProcessMsg(message); err != nil {
				log.Error("Kafka Produce Error:", err)
			}

			log.Info("Message Produce to Kafka")

		// Broadcast a message to each Clients(WebSocket)
		case broadcast := <-s.broadcast:
			s.Broadcast(broadcast)
			log.Info("Message Broadcast")
		}
	}
}
