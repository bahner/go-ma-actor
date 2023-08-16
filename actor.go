package main

import (
	"context"
	"encoding/json"

	"github.com/bahner/go-ma/key"
	"github.com/bahner/go-ma/message"
	"github.com/bahner/go-space/p2p/host"
	p2pupsub "github.com/bahner/go-space/p2p/pubsub"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	log "github.com/sirupsen/logrus"
)

type Actor struct {
	Node *host.Host
	Key  *key.Key

	// Context for any operations that might need cancellation or timeouts.
	ctx context.Context

	// Pubsub service and related fields, tied together by the identity.
	ps        *p2pupsub.Service    // The pubsub service instance.
	topic     *pubsub.Topic        // The topic derived from the identity.
	broadcast *pubsub.Topic        // The topic for broadcasting messages.
	cast      *pubsub.Subscription // The subscription to the topic for sending messages to.
	listen    *pubsub.Subscription // The subscription to the topic for receiving messages.
}

func newActor(ctx context.Context, secret string) (*Actor, error) {

	keyset, err := key.UnpackKeyset(secret)
	if err != nil {
		return nil, err
	}

	node := host.New()
	node.AddOption(libp2p.Identity(keyset.IPNSKey.PrivKey))
	node.AddOption(libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/udp/0",
		"/ip6/::/tcp/0",
		"/ip6/::/udp/0"))

	ps, err := createAndInitPubSubService(ctx, node)

	a := &Actor{
		Node: node,
		Key:  keyset.IPNSKey,
		ctx:  ctx,
		ps:   ps,
	}

	if err := a.initTopicAndSubscription(); err != nil {
		return nil, err
	}

	// Start the actor's listening loop
	go a.Listen()

	return a, nil
}

func (a *Actor) initTopicAndSubscription() error {
	topicName := a.Key.IPNSName.String() // Assuming IPNSName has a String() method to convert it to a string.
	var err error

	// Create the topic.
	a.topic, err = a.ps.Sub.Join(topicName)
	if err != nil {
		return err
	}

	// Subscribe to the topic.
	a.cast, err = a.topic.Subscribe()
	if err != nil {
		return err
	}

	return nil
}

func (a *Actor) Listen() {
	for {
		msg, err := a.sendTopic.Next(a.ctx)
		if err != nil {
			// Log the error or handle it more gracefully.
			log.Errorf("Failed to get next message: %v", err)
			return
		}

		if msg.ReceivedFrom == a.Node.Node.ID() {
			continue
		}

		am := new(message.Message)
		if err := json.Unmarshal(msg.Data, am); err != nil {
			log.Debugf("Failed to unmarshal message: %v", err)
			continue
		}

		a.ProcessMessage(am)
	}
}

func (a *Actor) ProcessMessage(m *message.Message) {
	// Handle the message according to your application's logic.
	// For instance, this could involve updating the actor's state, triggering some action, etc.
}
