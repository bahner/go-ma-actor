package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func handleEvents(ctx context.Context, e *entity.Entity) {
	envelopes := e.Topic.SubscribeEnvelopes(ctx)
	for {
		fmt.Println("Waiting for messages...")
		select {
		case envelope, ok := <-envelopes:
			if !ok {
				fmt.Printf("Envelope channel closed, exiting...")
				return
			}
			fmt.Printf("Received envelope: %v", e)

			// Process the envelope and send a pong response
			m, err := envelope.Open(e.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				fmt.Printf("Error opening envelope: %v\n", err)
				continue
			}

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
