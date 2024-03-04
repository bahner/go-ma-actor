package ui

const (
	discoverUsage = "/discover"
	discoverHelp  = "Triggers a discovery of peers"
)

func (ui *ChatUI) triggerDiscovery() {

	ui.displaySystemMessage("Discovery process started...")
	ui.p.DiscoverPeers()
	ui.displaySystemMessage("Discovery process complete.")

}
