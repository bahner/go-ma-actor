package ui

import (
	"fmt"
	"sort"
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
	ui.msgBox.SetTitle(ui.e.Nick())
	ui.app.Draw()
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {

	plist := ui.p.ConnectedProctectedPeersNickList()
	sort.Strings(plist)

	// clear is thread-safe
	ui.peersList.Clear()

	for _, p := range plist {
		fmt.Fprintln(ui.peersList, p)
	}

	ui.app.Draw()
}

func (ui *ChatUI) refreshTitle() {
	ui.msgBox.SetTitle(ui.e.Nick())
	ui.app.Draw()
}
