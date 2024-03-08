package p2p

import (
	"context"
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/peer"
	log "github.com/sirupsen/logrus"
)

const protectInterval = time.Second * 5

// This function looks through the list of peers connetced peers and add them to the list of protected peers,
// if it's known to be a protected peer.
func (p *P2P) protectLoop(ctx context.Context) {
	// Create a ticker that ticks every specified interval.
	ticker := time.NewTicker(protectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit the loop
			return
		case <-ticker.C:
			// This block is executed every time the ticker ticks.
			peers := p.DHT.Host().Network().Peers()

			for _, pid := range peers {
				idStr := pid.String()
				if peer.IsKnown(idStr) && peer.IsAllowed(idStr) {
					if !p.DHT.Host().ConnManager().IsProtected(pid, ma.RENDEZVOUS) {
						log.Printf("Protecting previously unprotected peer %s", idStr)
						p.DHT.Host().ConnManager().TagPeer(pid, ma.RENDEZVOUS, peer.DEFAULT_TAG_VALUE)
						p.DHT.Host().ConnManager().Protect(pid, ma.RENDEZVOUS)
					}
				}
			}
		}
	}
}
