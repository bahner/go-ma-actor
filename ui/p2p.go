package ui

import "github.com/bahner/go-ma/p2p"

func (ui *ChatUI) triggerDiscovery() {

	// go ui.n.StartPeerDiscovery(ui.ctx, config.GetRendezvous())
	p2p.StartPeerDiscovery(ui.ctx, ui.n)
	ui.displaySystemMessage("Discovery process started...")
}
