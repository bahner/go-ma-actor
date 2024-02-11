package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleBroadcastCommand(args []string) {

	if len(args) > 1 {

		recipient := ui.e.DID.String()

		var message string
		if len(args) > 2 {
			message = strings.Join(args[1:], " ")
		} else {
			message = args[1]
		}
		msgBytes := []byte(message)
		if log.GetLevel() == log.DebugLevel {
			ui.displaySystemMessage(fmt.Sprintf("Broadcasting %s to %s", message, recipient))
		} else {
			ui.displaySystemMessage(fmt.Sprintf("Broadcasting to %s", recipient))
		}

		msg, err := msg.NewBroadcast(ui.a.DID.String(), recipient, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Broadcast creation error: %s", err))
		}

		resp, err := entity.GetOrCreate(recipient, false) // Get the latest version of the entity.
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Entity creation error: %s", err))
		}

		err = msg.Broadcast(context.Background(), resp.Topic)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Broadcast error: %s", err))
		}

		log.Debugf("Message broadcasted to topic: %s", ui.e.Topic.String())
	} else {
		ui.handleHelpBroadcastCommand(args)
	}

}

func (ui *ChatUI) handleHelpBroadcastCommand(args []string) {
	ui.displaySystemMessage("Usage: /broadcast <message>")
	ui.displaySystemMessage("Sends a public announcement to the current entity")
}
