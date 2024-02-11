package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/curve25519"
)

// handleIncomingEnvelopes handles incoming envelopes to an entity. The actor
// is responsible for decrypting the envelope. The entity
// is only provided in order to decide whether to accept the message or not.
func (ui *ChatUI) handleIncomingEnvelopes(ctx context.Context, a *entity.Entity) {
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

			if a.Keyset == nil {
				log.Errorf("Actor %s has no keyset, cannot open envelope", a.DID)
				continue
			}

			// Check if privkey is a non-zero byte array.
			if a.Keyset.EncryptionKey.PrivKey == [curve25519.ScalarSize]byte{} {
				log.Errorf("Actor %s has zero-byte privkey. Unable to decrypt envelope.", a.DID)
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
			ui.displayPrivateMessage(m)
		}
	}
}
