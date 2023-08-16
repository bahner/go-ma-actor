package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bahner/go-ma/key"
	"github.com/bahner/go-ma/message"
	"github.com/bahner/go-space/p2p/host"
	p2pPupsub "github.com/bahner/go-space/p2p/pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

type Room struct {
	actor    *Actor
	Messages chan *message.Message
	roomName string
	nick     string
}

func newRoom(ctx context.Context, ps *p2pPupsub.Service, node host.Host, k key.Key, nickname, roomName string) (*Room, error) {
	r := &Room{
		roomName: roomName,
		nick:     nickname,
		actor: &Actor{
			Node: node,
			Key:  k,
			ctx:  ctx,
			ps:   ps,
		},
	}

	if err := r.actor.InitTopicAndSubscription(); err != nil {
		return nil, err
	}

	// Start the actor's listening loop
	go r.actor.Listen()

	return r, nil
}

// Override ProcessMessage to handle room-specific actions
func (r *Room) ProcessMessage(m *message.Message) {
	// Add the message to the Room's Messages channel
	r.Messages <- m
}

func (r *Room) Publish(content string) error {
	m, err := message.New()
	if err != nil {
	k, err := key.NewFromEncodedPrivKey(secret)
	if err != nil {
		return fmt.Errorf("failed to create key: %v", err)
	}
	m.Sign(k.PrivKey)
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	log.Debugf("Publishing message: %s", string(msgBytes))

	if err = r.actor.topic.Publish(r.actor.ctx, msgBytes); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// Example to show how you can use Actor's function in Room.
func (r *Room) ListPeers() []peer.ID {
	return r.actor.ps.ListPeers()
}

func (r *Room) JoinRoom() error {
	return r.actor.InitTopicAndSubscription()
}

func (r *Room) ListenToMessages() {
	for {
		msg, err := r.actor.sendTopic.Next(r.actor.ctx)
		if err != nil {
			// Handle error
			log.Errorf("Failed to get next message: %v", err)
			return
		}
		if msg.ReceivedFrom == r.actor.Node.Node.ID() {
			continue
		}

		roomMessage := new(message.Message)
		if err := json.Unmarshal(msg.Data, roomMessage); err != nil {
			log.Debugf("Failed to unmarshal message: %v", err)
			continue
		}
		r.Messages <- roomMessage
	}
}

// You can further extend the Room's functionalities as needed.
