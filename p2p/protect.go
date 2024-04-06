package p2p

import (
	"context"
	"time"

	"github.com/bahner/go-ma-actor/p2p/peer"
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
			peers := p.Host.Network().Peers()

			for _, pid := range peers {
				if peer.IsKnown(pid) {
					peer.Protect(p.Host, pid)
				}
			}
		}
	}
}
