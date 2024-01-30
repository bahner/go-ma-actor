package ui

func (ui *ChatUI) triggerDiscovery() {

	ui.p.DiscoverPeers()
	ui.displaySystemMessage("Discovery process started...")
}
