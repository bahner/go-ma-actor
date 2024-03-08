package mdns

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// DiscoveryLoop starts the discovery process and connects to discovered peers until the context is cancelled.
func (m *MDNS) DiscoveryLoop(ctx context.Context) error {
	log.Debugf("Discovering MDNS peers for service name: %s", m.rendezvous)

	for {
		select {
		case pai, ok := <-m.PeerChan:
			if !ok {
				if !(ctx == nil) { // conext.Bacground() is nil
					log.Fatalf("MDNS peer channel closed, ut was supposed to be running in the background.")
				}
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
