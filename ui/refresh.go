package ui

import (
	"fmt"
	"sort"

	"github.com/bahner/go-ma-actor/p2p/peer"
)

const (
	refreshUsage = "/refresh"
	refreshHelp  = "Refreshes the chat windows"
)

func (ui *ChatUI) handleRefreshCommand(args []string) {
	if len(args) == 1 {
		ui.handleRefresh()
	} else {
		ui.handleHelpCommand(refreshUsage, refreshHelp)
	}
}

func (ui *ChatUI) handleRefresh() {
	ui.refreshPeers()
	ui.msgBox.Clear()
	ui.setupInputField()
	ui.msgBox.SetTitle(ui.e.Nick)
	ui.app.Draw()
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {

	// Tweak this to change the timeout for peer discovery
	peers := ui.p.ConnectedProtectedPeersAddrInfo()

	// clear is thread-safe
	ui.peersList.Clear()

	plist := []string{}

	for _, p := range peers {
		n, err := peer.LookupNick(p.ID.String())
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Error looking up nick for %s: %s", p.ID, err))
			n = p.ID.ShortString()
		}
		plist = append(plist, n)
	}

	sort.Strings(plist)

	for _, p := range plist {
		fmt.Fprintln(ui.peersList, p)
	}

	ui.app.Draw()
}
