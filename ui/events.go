package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bahner/go-ma"
	log "github.com/sirupsen/logrus"
)

// // displaySelfMessage writes a message from ourself to the message window,
// // with our nick highlighted in yellow.
// func (ui *ChatUI) displaySelfMessage(msg string) {
// 	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.e.Nick))
// 	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
// }

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

func (ui *ChatUI) displaySystemMessage(msg string) {
	prompt := withColor("cyan", "[SYSTEM]:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.

func (ui *ChatUI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.chInput:
			if strings.HasPrefix(input, "/") {
				log.Debug("hadleEvents got command: ", input)
				ui.handleCommands(input)
			} else {
				log.Debug("hadleEvents got message: ", input)
				ui.handleChatMessage(input)
			}

		case m := <-ui.chMessage:
			if m.MimeType == ma.BROADCAST_MIME_TYPE {
				ui.displayBroadcastMessage(m)
				continue
			}
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			ui.refreshPeers()

		case <-ui.chDone:
			return
		}
	}
}
