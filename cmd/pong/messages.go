package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleMessageEvents(ctx context.Context, a *actor.Actor) {
	log.Debugf("Starting handleMessageEvents for %s", a.Entity.DID.Id)

	me := a.Entity.DID.Id

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info("handleMessageEvents: context cancelled, exiting...")
			return

		// Lemme think about this.
		// case <-a.Ctx.Done(): // Check for cancellation signal from the actor as well
		// 	log.Info("handleMessageEvents: actor context cancelled, exiting...")
		// 	return

		case m, ok := <-a.Entity.Messages: // Attempt to receive a message
			if !ok {
				log.Debugf("messageEvents: channel closed, exiting...")
				return
			}

			if m == nil {
				log.Debugf("messageEvents: received nil message, ignoring...")
				continue
			}

			if m.Verify() != nil {
				log.Debugf("messageEvents: failed to verify message: %v", m)
				continue
			}

			log.Debugf("Handling message: %v from %s to %s", string(m.Content), m.From, m.To)

			if m.From == me {
				log.Debugf("Received message from self, ignoring...")
				continue
			}

			// Only broadcast to broadcasts. Reply to messages.
			if m.To == me && m.MimeType == ma.BROADCAST_MIME_TYPE {
				log.Debugf("Received broadcast from %s to %s", m.From, m.To)
				log.Debugf("Sending broadcast announcement to %s over %s", m.From, me)
				err := broadcast(ctx, a)
				if err != nil {
					log.Errorf("Error sending public announcement: %v", err)
				}
				continue
			}
			if m.To == me && m.MimeType == ma.MESSAGE_MIME_TYPE {
				log.Debugf("Received message from %s to %s", m.From, m.To)
				log.Debugf("Sending reply to %s over %s", m.From, me)
				err := reply(ctx, a, m)
				if err != nil {
					log.Errorf("Error sending message: %v", err)
				}
				continue
			}
		}
	}
}

func broadcast(ctx context.Context, a *actor.Actor) error {

	me := a.Entity.DID.Id

	// Public announcements all go to the same topic, which is the DID of the actor's entity.
	topic := me

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	r, err := msg.NewBroadcast(topic, []byte("Public Announcment: "+viper.GetString("pong.msg")), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = r.Broadcast(ctx, a.Entity.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending signed broadcast over %s", topic)

	return nil
}

func reply(ctx context.Context, a *actor.Actor, m *msg.Message) error {

	// We need to reverse the to and from here. The message is from the other actor, and we are sending to them.
	to := m.From
	from := m.To

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	r, err := msg.New(from, to, []byte("Private reply: "+viper.GetString("pong.msg")), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = r.Send(ctx, a.Entity.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending private message to %s over %s", to, a.Entity.Topic.String())

	return nil
}
