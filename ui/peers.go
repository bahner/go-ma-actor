package ui

import (
	"fmt"
	"sort"

	"github.com/bahner/go-ma-actor/p2p/peer"
)

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {

	// Tweak this to change the timeout for peer discovery
	peers := ui.p.GetConnectedProtectedPeersAddrInfo()

	// clear is thread-safe
	ui.peersList.Clear()

	plist := []string{}

	for _, p := range peers {

		ap, err := peer.GetOrCreate(p)

		if err == nil {
			plist = append(plist, ap.Nick)
		}
	}

	sort.Strings(plist)

	for _, p := range plist {
		fmt.Fprintln(ui.peersList, p)
	}

	ui.app.Draw()
}

func (ui *ChatUI) handleHelpRefreshCommand(args []string) {
	ui.displaySystemMessage("Usage: /refresh")
	ui.displaySystemMessage("Refreshes the list of peers")
}
