package main

import (
	"context"

	"github.com/bahner/go-ma-actor/entity/actor"
	log "github.com/sirupsen/logrus"
)

func handleEnvelopeEvents(ctx context.Context, a *actor.Actor) {
	me := a.Entity.DID.Id

	log.Debugf("Starting handleEnvelopeEvents for %s", me)

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info("handleEnvelopeEvents: context cancelled, exiting...")
			return
		case env, ok := <-a.Envelopes: // Attempt to receive an envelope
			if !ok {
				log.Debugf("Envelope channel closed, exiting...")
				return
			}

			m, err := env.Open(a.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("Error opening envelope: %v", err)
				// Ensure m is not nil before calling Verify to avoid a panic
				if m != nil && m.Verify() != nil {
					log.Debugf("Failed to open envelope and verify message: %v", m)
				}
				continue
			}

			log.Debugf("Replying privately to message %v from %s", string(m.Content), m.From)
			err = reply(ctx, a, m)
			if err != nil {
				log.Errorf("Error replying to message: %v", err)
			}
		}
	}
}
