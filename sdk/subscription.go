package sdk

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type Subscription struct {
	ctx      context.Context
	err      error
	ch       chan *Message
	cancelCh chan<- *Subscription
	conn     *websocket.Conn
	once     sync.Once
}

func (s *Subscription) ReadMessage(ctx context.Context) (*Message, error) {
	select {
	case msg, ok := <-s.ch:
		if !ok {
			return nil, s.err
		}

		return msg, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *Subscription) Cancel() {
	s.conn.Close()

	select {
	case s.cancelCh <- s:
	case <-s.ctx.Done():
	}
}

func (s *Subscription) close() {
	s.once.Do(func() {
		close(s.ch)
	})
}
