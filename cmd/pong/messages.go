package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

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

	err = reply.Send(ctx, to.Topic)
	if err != nil {
		return fmt.Errorf("failed sending message: %w", errors.Cause(err))
	}

	log.Debugf("Sending reply to %s over %s", reply.To, to.Topic.String())

	return nil
}
