package sdk

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Message struct {
	// ReceivedFrom string
	Data []byte
}

type PManager struct {
	ctx      context.Context
	pmAddr   string
	conn     *websocket.Conn
	mySubs   map[*Subscription]struct{}
	sendMsg  chan *Message
	cancelCh chan *Subscription
	addSub   chan *addSubReq
}

func NewPManager(ctx context.Context) (*PManager, error) {
	pm := &PManager{
		ctx:      ctx,
		mySubs:   make(map[*Subscription]struct{}),
		sendMsg:  make(chan *Message, 32),
		cancelCh: make(chan *Subscription),
		addSub:   make(chan *addSubReq),
	}

	pm.dialToPropagationManager()

	go pm.processPManager(ctx)

	go pm.readMessage()

	return pm, nil
}

func (pm *PManager) WriteMessage(msgData []byte) {
	msg := &Message{
		Data: msgData,
	}

	pm.sendMsg <- msg
}

func (pm *PManager) Subscribe() (*Subscription, error) {
	sub := &Subscription{
		ctx:  pm.ctx,
		conn: pm.conn,
	}

	if sub.ch == nil {
		sub.ch = make(chan *Message, 32)
	}

	out := make(chan *Subscription, 1)

	select {
	case pm.addSub <- &addSubReq{
		sub:  sub,
		resp: out,
	}:

	case <-pm.ctx.Done():
		return nil, pm.ctx.Err()
	}

	return <-out, nil
}

func (pm *PManager) processPManager(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-pm.sendMsg:
			pm.publishMessageToGoSphere(msg)

		case sub := <-pm.addSub:
			pm.handleAddSubscription(sub)

		case sub := <-pm.cancelCh:
			pm.handleRemoveSubscription(sub)

		case <-ctx.Done():
			log.Info("processing PManager loop shut down")
			return
		}
	}
}

func (pm *PManager) dialToPropagationManager() error {
	return pm.tryDial(0, 3)
}

func (pm *PManager) tryDial(attempt int, maxAttempts int) error {
	c, err := loadConfig()
	if err != nil {
		log.Println("failed to load config: ", err)
		return err
	}
	pm.pmAddr = c.Path

	dialer := websocket.Dialer{
		ReadBufferSize:  1024 * 1024 * 20, // 2MB
		WriteBufferSize: 1024 * 1024 * 20, // 2MB
	}

	conn, _, err := dialer.Dial(pm.pmAddr, nil)
	if err != nil {
		log.Errorf("failed to dial to propagation module: %v", err)
		if attempt < maxAttempts {
			time.Sleep(1 * time.Second)
			return pm.tryDial(attempt+1, maxAttempts)
		}
		return err
	}

	pm.conn = conn
	return nil
}

func (pm *PManager) readMessage() {
	for {
		select {
		case <-pm.ctx.Done():
			return

		default:
			msgType, msgData, err := pm.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, 1006) {
					for {
						if err := pm.dialToPropagationManager(); err == nil {
							log.Info("Reconnection successful.")
							break
						}

						time.Sleep(1 * time.Second)
					}
				} else {
					log.Printf("ReadMessage error: %v", err)
					return
				}
			}

			if msgType == websocket.BinaryMessage {
				subNode := pm.mySubs

				for sub := range subNode {
					select {
					case sub.ch <- &Message{Data: msgData}:
					default:
						log.Error("Failed to send message to subscriber")
					}
				}
			}
		}
	}
}

func (pm *PManager) publishMessageToGoSphere(msg *Message) error {
	if err := pm.conn.WriteMessage(websocket.BinaryMessage, msg.Data); err != nil {
		log.Error("Failed to write message to propagation module: ", err)
		return err
	}
	return nil
}

type addSubReq struct {
	sub  *Subscription
	resp chan *Subscription
}

func (pm *PManager) handleAddSubscription(req *addSubReq) {
	sub := req.sub

	sub.cancelCh = pm.cancelCh

	pm.mySubs[sub] = struct{}{}

	req.resp <- sub
}

func (pm *PManager) handleRemoveSubscription(sub *Subscription) {
	sub.err = errors.New("Cancelled the subscription")
	sub.close()

	delete(pm.mySubs, sub)
}

func (pm *PManager) reConnection() error {
	if err := pm.dialToPropagationManager(); err == nil {
		log.Info("Success to reconnected w/propagation module")
		return nil
	} else {
		return err
	}
}
