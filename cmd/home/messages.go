package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// func openEnvelope(a *actor.Actor, e *msg.Envelope) (*msg.Message, error) {
// 	m, err := e.Open(a.Entity.Keyset.EncryptionKey.PrivKey[:])
// 	if err != nil {
// 		return nil, fmt.Errorf("failed unpacking message: %w", errors.Cause(err))
// 	}

// 	if log.GetLevel() == log.DebugLevel {
// 		fmt.Printf("Received message: %v\n", m.Content)
// 		jsonData, err := json.Marshal(m)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed packing message: %w", errors.Cause(err))
// 		}
// 		fmt.Printf("Received message: %v\n", string(jsonData))
// 	}

// 	err = m.Verify()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed verifying message: %w", errors.Cause(err))
// 	}

// 	return m, nil
// }

func pong(ctx context.Context, a *actor.Actor, m *msg.Message) error {
	to, err := topic.GetOrCreate(m.From)
	if err != nil {
		return fmt.Errorf("failed subscribing to recipients topic: %w", errors.Cause(err))
	}

	p, err := msg.New(m.To, m.From, []byte("Pong!"), "text/plain", a.Entity.Keyset.SigningKey.PrivKey)
	if err != nil {
		return fmt.Errorf("failed creating new message: %w", errors.Cause(err))
	}

	p.Send(ctx, to.Topic)

	log.Debugf("Sending pong to %s", p.To)

	return nil
}
