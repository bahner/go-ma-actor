package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

const (
	msgUsage = "/@<DID|NICK> <message>"
	msgHelp  = "Sends a private message directly the specified DID"
)

func (ui *ChatUI) handleMsgCommand(input string) {

	parts := strings.SplitN(input, " ", 2)

	if len(parts) == 2 {

		recipient := parts[0][1:] // The recipient is the first argument, without the leading @
		if !did.IsValid(recipient) {
			recipient = entity.GetDID(recipient)
		}

		if recipient == "" {
			ui.displaySystemMessage(fmt.Sprintf("Invalid DID: %s", recipient))
			return
		}

		message := parts[1]
		msgBytes := []byte(message)
		if log.GetLevel() == log.DebugLevel {
			ui.displaySystemMessage(fmt.Sprintf("Sending message to %s: %s", recipient, message))
		}

		msg, err := msg.New(ui.a.Entity.DID.Id, recipient, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message creation error: %s", err))
		}

		ui.displaySentPrivateMessage(msg)

		envelope, err := msg.Enclose()
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("envelope creation error: %s", err))
		}

		recp, err := entity.GetOrCreate(recipient)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("entity creation error: %s", err))
		}

		if recp == nil {
			ui.displaySystemMessage(fmt.Sprintf("entity not found: %s", recipient))
			return
		}

		// Connect to the entity's node, so we establish contact for the future.
		// A web of nodes, like. A web of trust innit.
		_, err = recp.ConnectPeer()
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("peer connection error: %s", err))
			ui.displaySystemMessage(fmt.Sprintf("sending message through the clouds %s", recipient))
		}

		// FIXME: get direct messaging to work.
		// err = msg.Send(context.Background(), recp.Topic)
		// Send private message in the entity's context. It's a whisper.
		// But should've been sent to the actor, not the entity. A loveletter, like.
		medium := ui.e.Topic
		err = envelope.Send(context.Background(), medium)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("message publishing error: %s", err))
		}
		log.Debugf("Message published to topic: %s", medium.String())
	} else {
		ui.handleHelpCommand(msgUsage, msgHelp)
	}
}
