package ui

import (
	"context"
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
			ui.handleHelpSetDiscoveryCommand(args)
		}
	} else {
		ui.handleHelpSetDiscoveryCommand(args)
	}
}

func (ui *ChatUI) handleHelpSetDiscoveryCommand(args []string) {
	ui.displayHelpUsage("/set discovery on|off")
	ui.displayHelpText("Toggles the discovery loop on and off")
}
