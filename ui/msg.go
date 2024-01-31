package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleMsgCommand(args []string) {

	if len(args) > 1 {

		recipient := alias.GetEntityDID(args[1])
		if !did.IsValidDID(recipient) {
			ui.displaySystemMessage(fmt.Sprintf("Invalid DID: %s", recipient))
			return
		}

		if len(args) < 3 {
			ui.displaySystemMessage("Message can't be empty")
			return
		}

		var message string
		if len(args) > 3 {
			message = strings.Join(os.Args[2:], " ")
		} else {
			message = args[2]
		}
		msgBytes := []byte(message)

		msg, err := msg.New(ui.a.DID.String(), recipient, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message creation error: %s", err))
		}

		err = msg.Send(ui.e.Ctx, ui.e.Topic.Topic)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message publishing error: %s", err))
		}
		log.Debugf("Message published to topic: %s", ui.e.Topic.Topic.String())
	} else {
		ui.displaySystemMessage("Usage: /msg <DID> <MESSAGE>")
	}
}
