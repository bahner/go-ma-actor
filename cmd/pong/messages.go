package main

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleMessageEvents(a *entity.Entity) {
	for {
		log.Info("Waiting for messages...")
		select {

		case <-a.Ctx.Done():
			log.Errorf("pong/handleMessageEvents: Actor context done, exiting...")
			return
		case m, ok := <-a.Messages:
			if !ok {
				log.Debugf("Message channel closed, exiting...")
				return
			}
			log.Debugf("Handling message: %v from %s to %s", string(m.Content), m.From, m.To)

			// Check if the message is from self to prevent pong loop
			if m.From == a.DID.String() {
				log.Debugf("Received message from self, ignoring...")
				continue
			}

			log.Debugf("Sending pong to %s over %s", m.From, a.DID.String())
			err := reply(a, m)
			if err != nil {
				log.Errorf("Error sending pong: %v", err)
			}
		}
	}
}

func reply(a *entity.Entity, m *msg.Message) error {

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	reply, err := msg.New(m.To, m.To, []byte(viper.GetString("pong.msg")), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = reply.Sign(a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed signing message: %w", errors.Cause(err))
	}

	err = reply.SendPublic(a.Ctx, a.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending broadcast over %s", a.Topic.String())

	return nil
}
