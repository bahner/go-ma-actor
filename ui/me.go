package ui

const (
	meUsage = "/me nick|who|where"
	meHelp  = "Shows your own DID or the last known location of your DID"
)

// handleMeCommands handles the "/me" commands in the chat UI.
// It takes a slice of strings as the arguments for the command.
// If the second argument is "who", it calls the handleWhoAmICommand function.
// If the second argument is "where", it calls the handleWhereAmICommand function.
// Otherwise, it displays a system message indicating an unknown alias node command.
// If the number of arguments is less than 2, it displays a system message indicating the correct usage.
func (ui *ChatUI) handleMeCommands(args []string) {

	if len(args) == 2 {

		switch args[1] {
		case "who":
			ui.handleMeWhoCommand(args)
			return
		case "where":
			ui.handleMeWhereCommand(args)
			return
		case "nick":
			ui.handleMeNickCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown /me command: " + args[1])
		}
	}
	ui.displayHelpUsage(meUsage)
}

// handleAliasCommand handles the /alias command
// handleMeWhoCommand displays the ID of the current user.
// If the number of arguments is 2, it displays the user's ID.
// Otherwise, it delegates the handling to the handleHelpMeCommands function.
func (ui *ChatUI) handleMeWhoCommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.a.Entity.DID.Id)
	} else {
		ui.handleHelpCommand(meUsage, meHelp)
	}

}
func (ui *ChatUI) handleMeWhereCommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.e.DID.Id)
	} else {
		ui.handleHelpCommand(meUsage, meHelp)
	}

}
func (ui *ChatUI) handleMeNickCommand(args []string) {

	if len(args) == 2 {
		ui.displaySystemMessage(ui.a.Entity.Nick)
	} else {
		ui.handleHelpCommand(meUsage, meHelp)
	}

}
