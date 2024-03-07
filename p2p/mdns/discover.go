package mdns

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

var ErrNoProtectedPeersFound = errors.New("protected peers not found")

// DiscoverPeers starts the discovery process and connects to discovered peers until the context is cancelled.
func (m *MDNS) DiscoverPeers(ctx context.Context) error {
	log.Debugf("Discovering MDNS peers for service name: %s", m.rendezvous)

	for {
		select {
		case pai, ok := <-m.PeerChan:
			if !ok {
				log.Debug("MDNS peer channel closed.")
				return nil // Exit if the channel is closed
			}

			if pai.ID == m.h.ID() {
				continue // Skip self connection
			}

			err := m.peerConnectAndUpdateIfSuccessful(ctx, pai)
			if err != nil {
				log.Debugf("Failed connecting to MDNS peer: %s, error: %v", pai.ID, err)
				continue // Skip if connection failed
			}

			log.Info("Successfully connected to MDNS peer: ", pai.ID)

		case <-ctx.Done():
			log.Info("Context cancelled, stopping MDNS peer discovery.")
			return nil // Stop the discovery loop if the context is done
		}
	}
}
