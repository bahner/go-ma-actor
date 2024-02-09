package main

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func handleEnvelopeEvents(ctx context.Context, a *entity.Entity) {

	log.Debugf("Starting handleEnvelopeEvents for %s", a.DID.String())

	for {
		log.Info("Waiting for messages...")
		env, ok := <-a.Envelopes
		if !ok {
			log.Debugf("Message channel closed, exiting...")
			return
		}

		m, err := env.Open(a.Keyset.EncryptionKey.PrivKey[:])
		if err != nil {
			log.Errorf("Error opening envelope: %v", err)
			if m.Verify() != nil {
				log.Debugf("Failed to open envelope and verify message: %v", m)
				continue
			}
		}

		log.Debugf("Replying privately to message %v from %s", m.Content, m.From)
		reply(ctx, a, m)

	}
}
