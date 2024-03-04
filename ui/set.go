package ui

const (
	setUsage = "/set broadcast|discovery"
	setHelp  = "Toggles broadcast and peer discovery on and off"
)

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) == 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		case "discovery":
			ui.handleSetDiscoveryCommand(args)
		}
	} else {
		ui.handleHelpCommand(setUsage, setHelp)
	}

}
