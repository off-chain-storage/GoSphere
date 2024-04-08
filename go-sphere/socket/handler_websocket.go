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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("read error:", err)
			}
			return
		}

		if messageType == websocket.BinaryMessage {
			// Send block data to Kafka
			s.propose <- message

		} else {
			log.Error("websocket message received of type", messageType)
		}
	}
}
