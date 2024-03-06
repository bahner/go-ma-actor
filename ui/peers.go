package ui

import (
	"fmt"
	"sort"
	"time"

	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/spf13/viper"
)

const (
	refreshUsage = "/refresh"
	refreshHelp  = "Refreshes the list of peers"
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

func getUIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
func getUIPeersRefreshInterval() time.Duration {
	return viper.GetDuration("ui.refresh")
}
