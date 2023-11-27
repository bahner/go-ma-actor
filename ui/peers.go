package ui

import (
	"fmt"
	"sort"

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

	// Tweak this to change the timeout for peer discovery
	peers := ui.p.GetConnectedProtectedPeersAddrInfo()

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
