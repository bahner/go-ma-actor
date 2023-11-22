package actor

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (a *Actor) receiveEnvelopes() (*msg.Message, error) {

	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()

	envelopes := a.Topic.SubscribeEnvelopes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to inbox: %w", err)
	}

	e := <-envelopes
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from inbox: %w", err)
	}

	message, err := e.Open(a.Entity.Keyset.EncryptionKey.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to open envelope: %w", err)
	}

	return message, nil
}

func (a *Actor) openEnvelopes(sub *pubsub.Subscription) {
	for {
		select {
		case <-a.ctx.Done():
			// Exit goroutine when context is cancelled
			return
		default:
			// Read message from Inbox subscription
			if msg, err := a.receiveEnvelopes(); err == nil {
				a.Messages <- msg
			}
		}
	}
}
