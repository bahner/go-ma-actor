package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/p2p"
)

var aliases map[string]string

func init() {
	aliases = make(map[string]string)
}

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleAliasCommand(args []string) {

}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {
	peers := p2p.GetConnectedPeers()

	// clear is thread-safe
	ui.peersList.Clear()

	for _, p := range peers {
		fmt.Fprintln(ui.peersList, p)
	}

	ui.app.Draw()
}
