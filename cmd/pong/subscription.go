package main

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

// SUbscribe a to e's topic and handle messages
func subscriptionLoop(a *entity.Entity) {

	t := a.DID.String()

	log.Info("Starting to handle incoming subscription messages on topic: ", t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// STart the 2 underlying handlers which may contain some logic.
	go handleMessageEvents(ctx, a)
	go handleEnvelopeEvents(ctx, a)

	sub, err := a.Topic.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	// Create a channel for messages.
	messages := make(chan *p2ppubsub.Message, pubsubMessagesBuffersize)

	// Start an anonymous goroutine to bridge sub.Next() to the messages channel.
	go func() {
		for {
			// Assuming sub.Next() blocks until the next message is available,
			// and returns a message or an error.
			message, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., log, break, or continue based on the error type).
				log.Errorf("Error getting next message: %v", err)
				return // or continue based on your error handling policy
			}
			log.Debugf("handleSubscriptionMessages: Received message: %s", message.Data)

			// Assuming message is of the type you need; otherwise, adapt as necessary.
			select {
			case messages <- message:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debugf("Entity %s is cancelled, exiting subscription loop...", t)
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("Message channel %s closed, exiting...", t)
				return
			}

			// Firstly check if this is a public message. Its quicker.
			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Received message: %v\n", m)
				a.Messages <- m
				continue
			}

			env, err := msg.UnmarshalAndVerifyEnvelopeFromCBOR(message.Data)
			if err == nil {
				log.Debugf("handleSubscriptionMessages: Received envelope: %v\n", env)
				a.Envelopes <- env
				continue
			}
		}
	}
}
