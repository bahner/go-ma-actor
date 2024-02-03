package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleMessageEvents(ctx context.Context, e *entity.Entity) {
	err := e.Topic.Subscribe(ctx, e.Messages, e.Envelopes)
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
	}
	for {
		fmt.Println("Waiting for messages...")
		select {
		case m, ok := <-e.Messages:
			if !ok {
				fmt.Printf("Envelope channel closed, exiting...")
				return
			}
			fmt.Printf("Received envelope: %v", e)
			fmt.Printf("Received message: %v\n", string(m.Content))

			// Check if the message is from self to prevent pong loop
			if m.From != e.DID.String() {
				log.Debugf("Sending pong to %s over %s", m.From, e.DID.String())
				err := reply(ctx, e, m)
				if err != nil {
					log.Errorf("Error sending pong: %v", err)
				}
			} else {
				fmt.Println("Received message from self, ignoring...")
			}

		case <-ctx.Done():
			fmt.Println("Context done, exiting...")
			return
		}
	}
}

func reply(ctx context.Context, ent *entity.Entity, m *msg.Message) error {

	// Answer in the same channel, ie. my address. It's kinda like a broadcast to a "room"
	to, err := topic.GetOrCreate(ent.DID.String())
	if err != nil {
		return fmt.Errorf("failed subscribing to recipients topic: %w", errors.Cause(err))
	}

	reply, err := msg.New(m.To, m.From, []byte(viper.GetString("pong.msg")), "text/plain", ent.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	err = reply.SendPublic(ctx, to.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending reply to %s over %s", reply.To, to.Topic.String())

	return nil
}
