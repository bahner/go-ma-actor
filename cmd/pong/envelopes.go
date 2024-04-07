package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleEnvelopeEvents(ctx context.Context, a *actor.Actor) {
	me := a.Entity.DID.Id

	log.Debugf("Starting handleEnvelopeEvents for %s", me)

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info("handleEnvelopeEvents: context cancelled, exiting...")
			return
		case env, ok := <-a.Envelopes: // Attempt to receive an envelope
			if !ok {
				log.Debugf("Envelope channel closed, exiting...")
				return
			}

			m, err := env.Open(a.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("Error opening envelope: %v", err)
				if m != nil && m.Verify() != nil {
					log.Debugf("Failed to open envelope and verify message: %v", m)
				}
				continue
			}

			log.Debugf("Replying privately to message %v from %s", string(m.Content), m.From)
			err = envelopeReply(ctx, a, m)
			if err != nil {
				log.Errorf("Error replying to message: %v", err)
			}
		}
	}
}

func envelopeReply(ctx context.Context, a *actor.Actor, m *msg.Message) error {

	// We need to reverse the to and from here. The message is from the other actor, and we are sending to them.
	replyTo := m.From
	replyFrom := m.To
	replyMsg := []byte(viper.GetString("pong.reply"))

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	reply, err := msg.New(replyFrom, replyTo, replyMsg, "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	envelope, err := reply.Enclose()
	if err != nil {
		return fmt.Errorf("failed creating envelope: %w", errors.Cause(err))
	}

	err = envelope.Send(ctx, a.Entity.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	fmt.Printf("Sending private envelope to %s over %s\n", replyTo, a.Entity.Topic.String())

	return nil
}
