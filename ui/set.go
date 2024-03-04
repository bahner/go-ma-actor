package ui

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) == 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		case "discovery":
			ui.handleSetDiscoveryCommand(args)
		}
	} else {
		ui.handleHelpSetCommands()
	}

}

func (ui *ChatUI) handleHelpSetCommands() {
	ui.displayHelpUsage("/set broadcast|discovery on|off")
	ui.displayHelpText("Toggles broadcast and peer discovery on and off")
}
