package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleMsgCommand(args []string) {

	if len(args) > 1 {

		if len(args) < 3 {
			ui.displaySystemMessage("Message can't be empty")
			return
		}

		recipient := args[1]
		if !did.IsValidDID(recipient) {
			recipient = alias.LookupEntityNick(recipient)
		}
		if recipient == "" {
			ui.displaySystemMessage(fmt.Sprintf("Invalid DID: %s", args[1]))
			return
		}

		var message string
		if len(args) > 3 {
			message = strings.Join(args[2:], " ")
		} else {
			message = args[2]
		}
		msgBytes := []byte(message)
		if log.GetLevel() == log.DebugLevel {
			ui.displaySystemMessage(fmt.Sprintf("Sending message to %s: %s", recipient, message))
		} else {
			ui.displaySystemMessage(fmt.Sprintf("Sending message to %s", recipient))
		}

		msg, err := msg.New(ui.a.DID.String(), recipient, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message creation error: %s", err))
		}

		resp, err := entity.GetOrCreate(recipient)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("entity creation error: %s", err))
		}

		err = msg.Send(context.Background(), resp.Topic)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message publishing error: %s", err))
		}
		log.Debugf("Message published to topic: %s", ui.e.Topic.String())
	} else {
		ui.handleHelpMsgCommand(args)
	}
}

func (ui *ChatUI) handleHelpBroadcastCommand(args []string) {
	ui.displaySystemMessage("Usage: /broadcast <DID|NICK> <message>")
	ui.displaySystemMessage("Sends a public announcement to the specified DID")
}
