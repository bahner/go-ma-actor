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
	ui.displaySystemMessage("Looking up nick for " + ui.e.DID.Id)
	nick, err := peer.GetNickForID(ui.e.DID.Id)
	if err != nil {
		ui.displaySystemMessage(fmt.Sprintf("Error looking up nick: %s", err))
	}
	ui.displaySystemMessage("Found nick " + nick)
	ui.e.SetNick(nick)
	ui.refreshPeers()
	ui.msgBox.Clear()
	ui.setupInputField()
	ui.displaySystemMessage("Setting title to " + ui.e.Nick)
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

		ap, err := peer.GetOrCreateFromAddrInfo(p)

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
