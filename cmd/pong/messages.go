package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func handleMessageEvents(ctx context.Context, a *actor.Actor) {
	me := a.Entity.DID.Id

	log.Debugf("Starting handleMessageEvents for %s", me)

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info("handleMessageEvents: context cancelled, exiting...")
			return

		case m, ok := <-a.Entity.Messages: // Attempt to receive a message
			if !ok {
				log.Debugf("messageEvents: channel closed, exiting...")
				return
			}

			if m == nil {
				log.Debugf("messageEvents: received nil message, ignoring...")
				continue
			}

			if m.Message.Verify() != nil {
				log.Debugf("messageEvents: failed to verify message: %v", m)
				continue
			}

			content := string(m.Message.Content)
			from := m.Message.From
			to := m.Message.To
			_type := m.Message.Type

			log.Debugf("Handling message: %v from %s to %s", content, from, to)

			if from == me {
				log.Debugf("Received message from self, ignoring...")
				continue
			}

			if to == me && _type == msg.PLAIN {
				messageReply(ctx, a, m)
			}
		}
	}
}

func messageReply(ctx context.Context, a *actor.Actor, m *entity.Message) error {

	// Switch sender and receiver. Reply back to from :-)
	replyFrom := m.Message.To
	replyTo := m.Message.From
	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	r, err := msg.New(replyFrom, replyTo, reply(m), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = r.Send(ctx, a.Entity.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	fmt.Printf("Sending private message to %s over %s\n", replyTo, a.Entity.Topic.String())

	return nil
}
