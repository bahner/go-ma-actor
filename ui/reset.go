package ui

const (
	resetUsage = "/reset"
	resetHelp  = "Resets the chat application and P2P network"
)

func (ui *ChatUI) handleResetCommand(args []string) {
	if len(args) == 1 {
		ui.handleReset()
	} else {
		ui.handleHelpCommand(resetUsage, resetHelp)
	}
}

func (ui *ChatUI) handleReset() {

	ui.msgBox.Clear()

	// Cancel the broadcast loop and start it again
	ui.displaySystemMessage("Resetting broadcast channel...")
	ui.broadcastCancel()
	ui.initBroadcast()

	// Reset the actor
	ui.displaySystemMessage("Resetting actor...")
	ui.currentActorCancel()
	ui.startActor()

	// Refresh the peers
	ui.displaySystemMessage("Refreshing peers...")
	ui.refreshPeers()

	// Reenter the current entity
	ui.displaySystemMessage("Reentering entity...")
	ui.currentEntityCancel()
	ui.enterEntity(ui.e.DID.Id, true)

	ui.app.Draw()

}
