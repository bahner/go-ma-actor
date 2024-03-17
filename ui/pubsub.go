package ui

import (
	"context"
	"time"

	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

const DISCOVERY_INTERVAL = time.Minute

func (ui *ChatUI) pubsubPeersLoop(ctx context.Context) {

	ticker := time.NewTicker(DISCOVERY_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debugf("Checking for new peers to add to pubsub")
			for _, p := range ui.p.ConnectedProtectedPeersAddrInfo() {
				err := ui.addPeerToPubsub(p)
				if err == peer.ErrAddrInfoAddrsEmpty || err == peer.ErrAlreadyConnected {
					continue
				}
				log.Debugf("Added peer %s to pubsub", p.ID)
			}
		case <-ctx.Done():
			log.Warn("pubsubPeersLoop: context done")
			return
		}
	}
}

func (ui *ChatUI) addPeerToPubsub(pai p2peer.AddrInfo) error {

	if len(pai.Addrs) == 0 {
		return peer.ErrAddrInfoAddrsEmpty
	}

	// Don't add already connected peers
	for _, p := range ui.e.Topic.ListPeers() {
		if p == pai.ID {
			return peer.ErrAlreadyConnected
		}
	}

	pubsub.AddPeer(pai.ID)

	return nil

}
