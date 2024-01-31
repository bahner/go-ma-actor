package entity

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma/msg"
)

func (e *Entity) openEnvelopes(a *Entity) (*msg.Message, error) {

	var err error

	ctx, cancel := context.WithCancel(e.ctx)
	defer cancel()

	envelopes := e.Topic.SubscribeEnvelopes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to inbox: %w", err)
	}

	envelope := <-envelopes
	if err != nil {
		return nil, fmt.Errorf("failed to receive message from inbox: %w", err)
	}

	message, err := envelope.Open(e.Keyset.EncryptionKey.PrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to open envelope: %w", err)
	}

	return message, nil
}

// receiveMessages is a goroutine that receives messages from the Inbox subscription
// and sends them to the message channel.
// The entity is the entity we are talking to, ie. "the Room", so we need
// to pass in the actor for signing and decryption.
func (e *Entity) receiveMessages(a *Entity) {
	e.ctx, e.cancel = context.WithCancel(context.Background())

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
			if msg, err := e.openEnvelopes(a); err == nil {
				a.Messages <- msg
			}
		}
	}

}
