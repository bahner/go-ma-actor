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

	if len(args) >= 3 {

		recipient := args[1]
		if !did.IsValid(recipient) {
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

		msg, err := msg.New(ui.a.Entity.DID.Id, recipient, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message creation error: %s", err))
		}

		recp, err := entity.GetOrCreate(recipient)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("entity creation error: %s", err))
		}

		if recp == nil {
			ui.displaySystemMessage(fmt.Sprintf("entity not found: %s", recipient))
			return
		}

		// FIXME: get direct messaging to work.
		// err = msg.Send(context.Background(), recp.Topic)
		// Send private message in the entity's context. It's a whisper.
		// But should've been sent to the actor, not the entity. A loveletter, like.
		medium := ui.e.Topic
		err = msg.Send(context.Background(), medium)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message publishing error: %s", err))
		}
		log.Debugf("Message published to topic: %s", medium.String())
	} else {
		ui.handleHelpMsgCommand(args)
	}
}

func (ui *ChatUI) handleHelpMsgCommand(args []string) {
	ui.displaySystemMessage("Usage: /msg <DID|NICK> <message>")
	ui.displaySystemMessage("Sends a private message to the specified DID")
}
