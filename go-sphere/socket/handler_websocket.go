package socket

import "github.com/gofiber/contrib/websocket"

func (s *Service) wsHandler(c *websocket.Conn) {
	defer func() {
		s.unregister <- c
		c.Close()
	}()

	// Register the client
	s.register <- c

	for {
		messageType, message, err := c.ReadMessage()
		go s.SendUDPMessage(1, "Received message from blockchain node")

		if err != nil {
			switch {
			case websocket.IsCloseError(err, websocket.CloseNormalClosure):
				return

			case websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure):
				log.Errorf("Unexpected WebSocket close error: %v", err)
				return

			default:
				log.Errorf("WebSocket read error: %v", err)
				return
			}
		}

		// Send block data to Kafka
		if messageType == websocket.BinaryMessage {
			s.propose <- message
		} else {
			log.Error("websocket message received of type", messageType)
		}
	}
}
