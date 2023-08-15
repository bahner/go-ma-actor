package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bahner/go-myspace/message"
	"github.com/libp2p/go-libp2p/core/peer"

	p2pPupsub "github.com/bahner/go-myspace/p2p/pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const ChatRoomBufSize = 128

type ChatRoom struct {
	Messages chan *message.Message
	ctx      context.Context
	ps       *p2pPupsub.Service
	topic    *pubsub.Topic
	sub      *pubsub.Subscription
	roomName string
	self     peer.ID
	nick     string
}

func newChatRoom(ctx context.Context, ps *p2pPupsub.Service, nickname, roomName string) (*ChatRoom, error) {
	cr := &ChatRoom{
		ctx:      ctx,
		ps:       ps,
		roomName: roomName,
		self:     ps.Host.Node.ID(),
		nick:     nickname,
	}

	// Try to join the topic immediately upon creation
	if err := cr.join(); err != nil {
		return nil, err
	}

	return cr, nil
}

func (cr *ChatRoom) join() error {
	var err error

	cr.topic, err = cr.ps.Sub.Join(cr.roomName)
	if err != nil {
		return fmt.Errorf("failed to join room: %v", err)
	}

	cr.sub, err = cr.topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to room: %v", err)
	}

	cr.Messages = make(chan *message.Message, ChatRoomBufSize)

	// Start reading messages asynchronously
	go cr.readLoop()
	return nil
}

func (cr *ChatRoom) Publish(content string) error {
	m := message.New(cr.self.Pretty(), cr.nick, []byte(content))
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	if err = cr.topic.Publish(cr.ctx, msgBytes); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.topic.ListPeers()
}

func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			// Optionally log the error or handle it more gracefully
			close(cr.Messages)
			return
		}

		if msg.ReceivedFrom == cr.self {
			continue
		}

		cm := new(message.Message)
		if err := json.Unmarshal(msg.Data, cm); err != nil {
			// Optionally log the error
			continue
		}

		cr.Messages <- cm
	}
}
