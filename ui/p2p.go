package ui

func (ui *ChatUI) triggerDiscovery() {

	ui.displaySystemMessage("Discovery process started...")
	ui.p.DiscoverPeers()
	ui.displaySystemMessage("Discovery process complete.")

}

func (ui *ChatUI) handleHelpDiscoverCommand(args []string) {
	ui.displaySystemMessage("Usage: /discover")
	ui.displaySystemMessage("Triggers a discovery of peers")
}
