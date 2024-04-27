package socket

import "github.com/gofiber/contrib/websocket"

func (s *Service) Broadcast(msgData []byte) {
	for connection, c := range s.clients {
		go func(connection *websocket.Conn, c *Client) {
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.isClosing {
				return
			}
			if err := connection.WriteMessage(websocket.BinaryMessage, msgData); err != nil {
				c.isClosing = true
				log.Error("write error:", err)

				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
				s.unregister <- connection
			}
		}(connection, c)
	}
}
