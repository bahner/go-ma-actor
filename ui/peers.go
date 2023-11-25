package ui

import (
	"fmt"
	"sort"

	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/peer"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleAliasCommand(args []string) {

	if len(args) == 3 {
		p := peer.GetByAlias(args[1])
		p.Alias = args[2]
		ui.displaySystemMessage("Peer " + p.ID + " is now known as " + p.Alias)
	} else {
		ui.displaySystemMessage("Usage: /alias <current alias> <alias>")
	}

}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {
	peers := p2p.GetConnectedPeers()

	// clear is thread-safe
	ui.peersList.Clear()

	// Create a slice for aliases
	var aliases []string
	for _, p := range peers {
		ap := peer.GetOrCreate(p)
		aliases = append(aliases, ap.Alias)
	}

	// Sort the aliases
	sort.Strings(aliases)

	// Display sorted aliases
	for _, alias := range aliases {
		fmt.Fprintln(ui.peersList, alias)
	}

	ui.app.Draw()
}

// func (ui *ChatUI) handleStatusCommand(args []string) {
// 	if len(args) > 1 {
// 		switch args[1] {
// 		case "sub":
// 			ui.displayStatusMessage(ui.getStatusSub())
// 		case "topic":
// 			ui.displayStatusMessage(ui.getStatusTopic())
// 		case "host":
// 			ui.displayStatusMessage(ui.getStatusHost())
// 		default:
// 			ui.displaySystemMessage("Unknown status type: " + args[1])
// 		}
// 	} else {
// 		ui.displaySystemMessage("Usage: /status [sub|topic|host]")
// 	}
// }
