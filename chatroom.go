package main

import (
	"context"
	"encoding/json"

	"github.com/bahner/go-myspace/message"
	"github.com/libp2p/go-libp2p/core/peer"

	p2pPupsub "github.com/bahner/go-myspace/p2p/pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// ChatRoomBufSize is the number of incoming messages to buffer for each topic.
const ChatRoomBufSize = 128

// ChatRoom represents a subscription to a single PubSub topic. Messages
// can be published to the topic with ChatRoom.Publish, and received
// messages are pushed to the Messages channel.
type ChatRoom struct {
	// Messages is a channel of messages received from other peers in the chat room
	Messages chan *message.Message

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	roomName string
	self     peer.ID
	nick     string
}

func newChatRoom(ctx context.Context, ps *p2pPupsub.Service, nickname string, roomName string) (*ChatRoom, error) {

	id := ps.Host.Node.ID()

	return &ChatRoom{
		ctx:      ctx,
		ps:       ps.Sub,
		roomName: roomName,
		self:     id,
		nick:     nickname,
	}, nil
}

// joinChatRoom tries to subscribe to the PubSub topic for the room name, returning
// a ChatRoom on success.
func joinChatRoom(c *ChatRoom) (*ChatRoom, error) {
	// join the pubsub topic
	topic, err := ps.Sub.Join(c.roomName)
	if err != nil {
		return nil, err
	}
	log.Info("Entered room: ", c.roomName)

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	log.Infof("Subscribed to room: %s", c.roomName)

	c.topic = topic
	c.sub = sub
	c.Messages = make(chan *message.Message, ChatRoomBufSize)

	// start reading messages from the subscription in a loop
	go c.readLoop()
	return c, nil
}

// Publish sends a message to the pubsub topic.
func (cr *ChatRoom) Publish(msg string) error {

	m := message.New(cr.self.Pretty(), cr.nick, []byte(msg))
	// m := message.New(cr.self, cr.topic, message)
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.topic.Publish(cr.ctx, msgBytes)
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.topic.ListPeers()
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}

		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}
		cm := new(message.Message)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		// send valid messages onto the Messages channel
		cr.Messages <- cm
	}
}
