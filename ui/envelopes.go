package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	log "github.com/sirupsen/logrus"
)

// handleIncomingEnvelopes handles incoming envelopes to an entity. The actor
// is responsible for decrypting the envelope. The entity
// is only provided in order to decide whether to accept the message or not.
// The original Subscribe features the actor, So envelopes are sent here.
func handleIncomingEnvelopes(ctx context.Context, e *entity.Entity, a *actor.Actor) {
	log.Debugf("Waiting for actor envelopes")
	for {
		select {
		case <-ctx.Done():
			log.Debug("Context cancelled, exiting handleIncomingEnvelopes...")
			return
		case envelope, ok := <-a.Envelopes: // Envelopes should always have been sent to the actor.
			if !ok {
				log.Debug("Actor envelope channel closed, exiting...")
				return
			}
			log.Debugf("Received actor envelope: %v", envelope)

			err := a.Keyset.Verify()
			if err != nil {
				log.Errorf("handleIncomingEnvelope: %s: %v", a.Entity.DID.String(), err)
				continue
			}

			// Process the envelope and send a pong response
			m, err := envelope.Open(a.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("Error opening actor envelope: %v\n", err)
				continue
			}

			log.Debugf("Opened envelope and found message: %v\n", string(m.Content))

			// Send the message to the actor for processing. It can decide to ignore it.
			e.Messages <- m
		}
	}
}
