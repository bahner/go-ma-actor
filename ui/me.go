package ui

// handleMeCommands handles the "/me" commands in the chat UI.
// It takes a slice of strings as the arguments for the command.
// If the second argument is "who", it calls the handleWhoAmICommand function.
// If the second argument is "where", it calls the handleWhereAmICommand function.
// Otherwise, it displays a system message indicating an unknown alias node command.
// If the number of arguments is less than 2, it displays a system message indicating the correct usage.
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
// handleWhoAmICommand displays the ID of the current user.
// If the number of arguments is 2, it displays the user's ID.
// Otherwise, it delegates the handling to the handleHelpMeCommands function.
func (ui *ChatUI) handleWhoAmICommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.a.Entity.DID.Id)
	} else {
		ui.handleHelpMeCommands()
	}

}
func (ui *ChatUI) handleWhereAmICommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.e.DID.Id)
	} else {
		ui.handleHelpMeCommands()
	}

}

func (ui *ChatUI) handleHelpMeCommands() {
	ui.displaySystemMessage("Usage: /me who|where")
	ui.displaySystemMessage("Shows your own DID or the last known location of your DID")
}
