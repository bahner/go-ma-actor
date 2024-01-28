package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
	log "github.com/sirupsen/logrus"
)

func handleEvents(ctx context.Context, a *actor.Actor) {
	envelopes := a.Topic.SubscribeEnvelopes(ctx)
	for {
		fmt.Println("Waiting for messages...")
		select {
		case e, ok := <-envelopes:
			if !ok {
				fmt.Printf("Envelope channel closed, exiting...")
				return
			}
			fmt.Printf("Received envelope: %v", e)

			// Process the envelope and send a pong response
			m, err := e.Open(a.Entity.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				fmt.Printf("Error opening envelope: %v\n", err)
				continue
			}

			fmt.Printf("Received message: %v\n", string(m.Content))

			// Check if the message is from self to prevent pong loop
			if m.From != a.Entity.DID.String() {
				log.Debugf("Sending pong to %s over %s", m.From, a.Entity.DID.String())
				err := pong(ctx, a, m)
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
