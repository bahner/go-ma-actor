package ui

func (ui *ChatUI) handleMeCommands(args []string) {

	if len(args) >= 2 {
		switch args[1] {
		case "who":
			ui.handleWhoAmICommand(args)
		case "where":
			ui.handleWhereAmICommand(args)
		default:
			ui.displaySystemMessage("Unknown alias node command: " + args[2])
		}
	} else {

		ui.displaySystemMessage("Usage: /me who|where")

	}

}

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleWhoAmICommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.a.DID.String())
	} else {
		ui.handleHelpMeCommands(args)
	}

}
func (ui *ChatUI) handleWhereAmICommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.e.DID.String())
	} else {
		ui.handleHelpMeCommands(args)
	}

}

func (ui *ChatUI) handleHelpMeCommands(args []string) {
	ui.displaySystemMessage("Usage: /me who|where")
	ui.displaySystemMessage("Shows your own DID or the last known location of your DID")
}
