package main

import (
	"crypto/ed25519"
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/message"
	p2ppupsub "github.com/libp2p/go-libp2p-pubsub"
)

type Room struct {
	Subscription *p2ppupsub.Subscription
	Topic        *p2ppupsub.Topic
	DID          *did.DID // NB! "Room must have published it's DIDDocument."
	Doc          *doc.Document
	// The signing and encryption keys are used to verify and encrypt messages.
	SigningKey    *ed25519.PublicKey
	EncryptionKey *ed25519.PublicKey
	Messages      chan *message.Message
	// We can add objects etc here.
	Actor *Actor
	nick  string
}

// Override ProcessMessage to handle room-specific actions
func (r *Room) ProcessMessage(m *message.Message) {
	// Add the message to the Room's Messages channel
	r.Messages <- m
}

func NewRoom(d string) (*Room, error) {

	var r *Room
	var err error

	// Set the DID
	r.DID, err = did.NewFromDID(d)
	if err != nil {
		return nil, fmt.Errorf("room: failed to create DID: %v", err)
	}

	// Set nick to DID fragment
	r.nick = r.DID.Fragment

	// Fetch the public keys need to send and receive messages to the room
	r.Doc, err = doc.New(r.DID.String(), r.DID.String())
	if err != nil {
		return nil, fmt.Errorf("room: failed to create DOC: %v", err)
	}

	r.Doc, err = doc.Fetch(r.DID.String())
	if err != nil {
		return nil, fmt.Errorf("room: failed to fetch DOC: %v", err)
	}

	// Subscribe to the recipients topic
	r.Topic, err = ps.Sub.Join(r.DID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %v", err)
	}
	r.Subscription, err = r.Topic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %v", err)
	}

	// Create a new channel for messages
	r.Messages = make(chan *message.Message)

	return r, nil
}

// func (r *Room) Publish(content string) error {
// 	m, err := message.New()
// 	if err != nil {
// 		k, err := key.NewFromEncodedPrivKey(identity)
// 		if err != nil {
// 			return fmt.Errorf("failed to create key: %v", err)
// 		}
// 		m.Sign(k.PrivKey)
// 		msgBytes, err := json.Marshal(m)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal message: %v", err)
// 		}
// 		log.Debugf("Publishing message: %s", string(msgBytes))

// 		if err = r.Actor.Topic.Publish(r.Actor.ctx, msgBytes); err != nil {
// 			return fmt.Errorf("failed to publish message: %v", err)
// 		}

// 		return nil
// 	}
// }

// // Example to show how you can use Actor's function in Room.
// func (r *Room) ListPeers() []peer.ID {
// 	return r.Actor.PubSubService.Sub.ListPeers(r.DID.String())
// }

// func (r *Room) JoinRoom() error {
// 	return r.Actor.InitTopicAndSubscription()
// }

// func (r *Room) ListenToMessages() {
// 	for {
// 		msg, err := r.Actor.sendTopic.Next(r.Actor.ctx)
// 		if err != nil {
// 			// Handle error
// 			log.Errorf("Failed to get next message: %v", err)
// 			return
// 		}
// 		if msg.ReceivedFrom == r.Actor.node.Node.ID() {
// 			continue
// 		}

// 		roomMessage := new(message.Message)
// 		if err := json.Unmarshal(msg.Data, roomMessage); err != nil {
// 			log.Debugf("Failed to unmarshal message: %v", err)
// 			continue
// 		}
// 		r.Messages <- roomMessage
// 	}
// }

// // You can further extend the Room's functionalities as needed.

func (r *Room) Enter(a *Actor) {

	a.ears = r.Messages

}
