package ui

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) == 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		}
	} else {
		ui.handleHelpSetCommands(args)
	}

}

func (ui *ChatUI) handleHelpSetCommands(args []string) {
	ui.displaySystemMessage("Usage: /set broadcast on|off")
	ui.displaySystemMessage("For now toggles broadcast messages on and off")
}
