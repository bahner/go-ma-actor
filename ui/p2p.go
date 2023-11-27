package ui

func (ui *ChatUI) triggerDiscovery() {

	// go ui.n.StartPeerDiscovery(ui.ctx, config.GetRendezvous())
	ui.p.DiscoverPeers()
	ui.displaySystemMessage("Discovery process started...")
}
