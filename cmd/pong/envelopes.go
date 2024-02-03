package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func handleEnvelopeEvents(ctx context.Context, e *entity.Entity) {
	err := e.Topic.Subscribe(ctx, e.Messages, e.Envelopes)
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
	}
	for {
		fmt.Println("Waiting for envelopes...")
		select {
		case envelope, ok := <-e.Envelopes:
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

			log.Debugf("Received envelope from %s:", string(m.From))
			e.Messages <- m

		case <-ctx.Done():
			fmt.Println("Context done, exiting...")
			return
		}
	}
}
