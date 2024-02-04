package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleMessageEvents(a *entity.Entity) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		log.Info("Waiting for messages...")
		m, ok := <-a.Messages
		if !ok {
			log.Debugf("Message channel closed, exiting...")
			return
		}
		log.Debugf("Handling message: %v from %s to %s", string(m.Content), m.From, m.To)

		msgJSON, err := json.Marshal(m)
		if err != nil {
			log.Errorf("Error marshalling message: %v", err)
			continue
		}
		log.Debugf("Handling message: %v", string(msgJSON))
		// Check if the message is from self to prevent pong loop
		// NB! This is only for *incoming* messages, not broadcasts.
		if m.From == a.DID.String() {
			log.Debugf("Received message from self, ignoring...")
			continue
		}

		log.Debugf("Sending public announcement to %s over %s", m.From, a.DID.String())
		err = broadcast(ctx, a)
		if err != nil {
			log.Errorf("Error sending public announcement: %v", err)
		}

		log.Debugf("Sending private reply to %s over %s", m.From, a.DID.String())
		err = replyPrivately(ctx, a, m)
		if err != nil {
			log.Errorf("Error sending public announcement: %v", err)
		}
	}
}

func broadcast(ctx context.Context, a *entity.Entity) error {

	// Public announcements all go to the same topic, which is the DID of the actor.
	topic := a.DID.String()

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	r, err := msg.NewBroadcast(topic, topic, []byte("PA:"+viper.GetString("pong.msg")), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = r.Sign(a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed signing message: %w", errors.Cause(err))
	}

	err = r.Broadcast(ctx, a.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending signed broadcast over %s", topic)

	return nil
}

func replyPrivately(ctx context.Context, a *entity.Entity, m *msg.Message) error {

	// We need to reverse the to and from here. The message is from the other actor, and we are sending to them.
	to := m.From
	from := m.To

	// Broadcast are sent to the topic, and the topic is the DID of the recipient
	r, err := msg.New(from, to, []byte("private:"+viper.GetString("pong.msg")), "text/plain", a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = r.Sign(a.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed signing message: %w", errors.Cause(err))
	}

	err = r.Send(ctx, a.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending private message to %s over %s", to, a.Topic.String())

	return nil
}
