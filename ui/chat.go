package ui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	from := alias.LookupEntityDID(cm.From)
	prompt := withColor("green", fmt.Sprintf("<%s>:", from))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	// This could be a timeout for topic publishing
	ctx := context.Background()

	log.Debugf("Handling chatMessage: %s", input)
	msgBytes := []byte(input)
	log.Debugf("ui.a.DID.Fragment: %s", ui.a.DID.Fragment)
	log.Debugf("ui.e.ID: %s", ui.e.DID)
	msg, err := msg.NewBroadcast(ui.a.DID.String(), ui.e.DID.String(), msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
	if err != nil {
		log.Debugf("message creation error: %s", err)
		return fmt.Errorf("message creation error: %w", err)
	}

	if log.GetLevel() == log.DebugLevel {
		msgJson, _ := json.Marshal(msg)
		log.Debugf("Signed message: %s", msgJson)
		err = msg.Verify()
		if err != nil {
			log.Debugf("failed to verify my own message: %s", err)
			return fmt.Errorf("message verification error: %w", err)
		} else {
			log.Debugf("Message signature verified")
		}

		ui.displaySelfMessage(string(msgJson))
	}
	ui.displayChatMessage(msg)

	err = msg.Broadcast(ctx, ui.e.Topic)
	if err != nil {
		log.Debugf("message publishing error: %s", err)
		return fmt.Errorf("message publishing error: %w", err)
	}
	log.Debugf("Message published to topic: %s", ui.e.Topic.String())

	return nil
}
