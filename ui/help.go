package ui

func (ui *ChatUI) handleHelpCommands(args []string) {

	if len(args) == 1 {
		ui.displaySystemMessage("Usage: /help [command]")
		ui.displaySystemMessage("Available commands:")
		ui.displaySystemMessage("/help status")
		ui.displaySystemMessage("/help msg")
		ui.displaySystemMessage("/help discover")
		ui.displaySystemMessage("/help enter")
		ui.displaySystemMessage("/help alias")
		ui.displaySystemMessage("/help aliases")
		ui.displaySystemMessage("/help whereis")
		ui.displaySystemMessage("/help me # Pun intended")
		ui.displaySystemMessage("/help refresh")
		ui.displaySystemMessage("/help")
		ui.displaySystemMessage("/help quit")
		ui.displaySystemMessage("Type /help [command] for more information")
	} else {
		switch args[1] {
		case "status":
			ui.handleHelpStatusCommands(args)
		case "msg":
			ui.handleHelpMsgCommand(args)
		case "discover":
			ui.handleHelpDiscoverCommand(args)
		case "enter":
			ui.handleHelpEnterCommand(args)
		case "alias":
			ui.handleHelpAliasCommands(args)
		case "aliases":
			ui.handleHelpAliasesCommand(args)
		case "whereis":
			ui.handleHelpWhereisCommand(args)
		case "me":
			ui.handleHelpMeCommands(args)
		case "refresh":
			ui.handleHelpRefreshCommand(args)
		case "quit":
			ui.handleHelpQuitCommand(args)
		default:
			ui.handleHelpUnknownCommand(args)
		}
	}
}

func (ui *ChatUI) handleHelpQuitCommand(args []string) {
	ui.displaySystemMessage("Usage: /quit")
	ui.displaySystemMessage("Quits the chat client")
}

func (ui *ChatUI) handleHelpUnknownCommand(args []string) {
	ui.displaySystemMessage("Unknown command: " + args[0])
}
