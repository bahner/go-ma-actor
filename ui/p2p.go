package ui

const (
	p2pUsage         = "/p2p discover"
	p2pHelp          = "P2P commands only feature discovery at the moment"
	p2pDiscoverUsage = "/p2p discover"
	p2pDiscoverHelp  = "Triggers a discovery of peers"
)

func (ui *ChatUI) triggerDiscovery() {

	ui.displaySystemMessage("Discovery process started...")
	ui.p.DiscoverPeers()
	ui.displaySystemMessage("Discovery process complete.")

}
func (ui *ChatUI) handleP2PDiscoverCommand(args []string) {

	if len(args) == 2 {
		ui.triggerDiscovery()
	} else {
		ui.handleHelpCommand(p2pDiscoverUsage, p2pDiscoverHelp)
	}
}

func (ui *ChatUI) handleP2PCommand(args []string) {

	if len(args) == 2 {
		command := args[1]
		switch command {
		case "discover":
			ui.handleP2PDiscoverCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown p2p command: " + command)
		}
	}

	ui.handleHelpCommand(peerUsage, peerHelp)
}
