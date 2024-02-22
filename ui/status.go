package ui

import (
	"fmt"
)

func (ui *ChatUI) handleStatusCommand(args []string) {

	if len(args) == 1 {
		ui.displaySystemMessage(ui.getStatusHost())
		ui.displaySystemMessage(ui.getStatusTopic())
		return
	}

	if len(args) == 2 {
		switch args[1] {
		case "sub":
			ui.displayStatusMessage(ui.getStatusSub())
		case "topics":
			ui.displayStatusMessage(ui.getStatusTopic())
		case "host":
			ui.displayStatusMessage(ui.getStatusHost())
		default:
			ui.displaySystemMessage("Unknown status type: " + args[1])
		}
		return
	}

	ui.handleHelpStatusCommands(args)
}

// displayStatusMessage writes a status message to the message window.
func (ui *ChatUI) displayStatusMessage(msg string) {
	prompt := withColor("cyan", "<STATUS>:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) getStatusSub() string {
	return "not implemented yet"
}

func (ui *ChatUI) getStatusTopic() string {
	// Return whatever status you'd like about the topic.
	// Fetching peers as an example below:
	// peers := ui.keyAgreement.ListPeers()
	aConnected := ui.a.Entity.Topic.ListPeers()
	eConnected := ui.e.Topic.ListPeers()
	bConnected := ui.b.ListPeers()
	return fmt.Sprintf("\nEntity: %s\n%s\nActor: %s\n%s\nBroadcast: %s\n%s",
		ui.e.Topic.String(), eConnected[:],
		ui.a.Entity.Topic.String(), aConnected[:],
		ui.b.String(), bConnected[:],
	)
}

func (ui *ChatUI) getStatusHost() string {
	// Return whatever status you'd like about the host.
	// Just an example below:
	var result string
	result += "Peer ID: " + ui.p.Node.ID().String() + "\n"
	result += fmt.Sprintf("Peers no# %d\n", len(ui.p.Node.Network().Peers()))
	return result
}

func (ui *ChatUI) handleHelpStatusCommands(args []string) {
	ui.displaySystemMessage("Usage: /status topics|host")
	ui.displaySystemMessage("Displays the current status of the chat client")
}
