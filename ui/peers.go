package ui

import "fmt"

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {
	peers := ui.n.Network().Peers()

	// clear is thread-safe
	ui.peersList.Clear()

	for _, p := range peers {
		fmt.Fprintln(ui.peersList, p.ShortString())
	}

	ui.app.Draw()
}
