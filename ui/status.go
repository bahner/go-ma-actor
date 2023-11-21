package ui

import (
	"fmt"
)

func (ui *ChatUI) handleStatusCommand(args []string) {
	if len(args) > 1 {
		switch args[1] {
		case "sub":
			ui.displayStatusMessage(ui.getStatusSub())
		case "topic":
			ui.displayStatusMessage(ui.getStatusTopic())
		case "host":
			ui.displayStatusMessage(ui.getStatusHost())
		default:
			ui.displaySystemMessage("Unknown status type: " + args[1])
		}
	} else {
		ui.displaySystemMessage("Usage: /status [sub|topic|host]")
	}
}

// displayStatusMessage writes a status message to the message window.
func (ui *ChatUI) displayStatusMessage(msg string) {
	prompt := withColor("cyan", "<STATUS>:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) getStatusSub() string {
	// status := fmt.Sprintf("keyAgreement: %s", ui.a.keyAgreement.String())
	// status += fmt.Sprintf("assertionMethod: %s", ui.assertionMethod.String())
	return "not implemented yet"
}

func (ui *ChatUI) getStatusTopic() string {
	// Return whatever status you'd like about the topic.
	// Fetching peers as an example below:
	// peers := ui.keyAgreement.ListPeers()
	return "not implemented yet"
}

func (ui *ChatUI) getStatusHost() string {
	// Return whatever status you'd like about the host.
	// Just an example below:
	var result string
	result += "Peer ID: " + ui.n.ID().String() + "\n"
	result += "Peers:\n"
	for _, p := range ui.n.Network().Peers() {
		result += p.String() + "\n"
	}
	return result
}
