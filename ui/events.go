package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bahner/go-ma"
	log "github.com/sirupsen/logrus"
)

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

func (ui *ChatUI) displaySystemMessage(msg string) {
	out := withColor("purple", msg)
	fmt.Fprintf(ui.msgW, "%s\n", out)
}

func (ui *ChatUI) displayDebugMessage(m string) {
	if log.GetLevel() == log.DebugLevel {
		ui.displaySystemMessage(fmt.Sprintf("DEBUG: %s", m))
	}
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

			input = strings.TrimSpace(input)

			if strings.HasPrefix(input, "/") {
				log.Debug("handleEvents got command: ", input)
				ui.handleCommands(input)
				continue
			}
			if strings.HasPrefix(input, "@") {
				log.Debug("handleEvents got command: ", input)
				ui.handleMsgCommand(input)
				continue
			}
			ui.handleChatMessage(input)

		case m := <-ui.chMessages:
			if m.Type == ma.BROADCAST_MESSAGE_TYPE {
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
