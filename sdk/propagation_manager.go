package sdk

import (
	"context"

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

func (pm *PManager) dialToPropagationManager() error {
	c, err := loadConfig()
	if err != nil {
		log.Error("failed to load config: ", err)
		return err
	}
	pm.pmAddr = c.Path

	dialer := websocket.Dialer{
		ReadBufferSize:  1024 * 1024 * 2, // 2MB
		WriteBufferSize: 1024 * 1024 * 2, // 2MB
	}

	conn, _, err := dialer.Dial(pm.pmAddr, nil)
	if err != nil {
		log.Error("failed to dial to propagation module: ", err)
		return err
	}

	pm.conn = conn

	return nil
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

func (pm *PManager) readMessage() {
	for {
		select {
		case <-pm.ctx.Done():
			return

		default:
			log.Info("ERROR 01")
			msgType, msgData, err := pm.conn.ReadMessage()
			log.Info("ERROR 02")
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					log.Info("ERROR 1")
				} else if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Info("ERROR 2")
				} else {
					log.Errorf("ReadMessage error: %v", err)
					log.Info("ERROR 3")
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
			} else {
				log.Error("Received message is not proper type")
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
