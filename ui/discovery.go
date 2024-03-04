package ui

import (
	"context"
)

const (
	setDiscoveryUsage = "/set discovery on|off"
	setDiscoveryHelp  = "Toggles the discovery loop on and off"
)

func (ui *ChatUI) handleSetDiscoveryCommand(args []string) {

	if len(args) == 3 {

		toggle := args[2]

		switch toggle {
		case "on":

			// Now we can start continuous discovery in the background.
			ui.discoveryLoopCtx, ui.discoveryLoopCancel = context.WithCancel(context.Background())
			go ui.p.DiscoveryLoop(context.Background())
			ui.displaySystemMessage("Discovery is on")
		case "off":
			if ui.discoveryLoopCancel != nil {
				ui.discoveryLoopCancel()
				ui.displaySystemMessage("Discovery is off")
			}
		default:
			ui.handleHelpSetDiscoveryCommand()
		}
	} else {
		ui.handleHelpSetDiscoveryCommand()
	}
}

func (ui *ChatUI) handleHelpSetDiscoveryCommand() {
	ui.displayHelpUsage(setDiscoveryUsage)
	ui.displayHelpText(setDiscoveryHelp)
}
