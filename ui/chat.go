package ui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	from := did.GetFragment(cm.From)
	prompt := withColor("green", fmt.Sprintf("<%s>:", from))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	// This could be a timeout for topic publishing
	ctx := context.Background()

	log.Debugf("Handling chatMessage: %s", input)
	msgBytes := []byte(input)
	log.Debugf("ui.a.Entity.DID.Fragment: %s", ui.a.Entity.DID.Fragment)
	log.Debugf("ui.e.ID: %s", ui.e.DID)
	msg, err := msg.New(ui.a.Entity.DID.String(), ui.e.DID, msgBytes, "text/plain", ui.a.Entity.Keyset.SigningKey.PrivKey)
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
	t, err := topic.GetOrCreate(ui.e.DID)
	if err != nil {
		log.Debugf("topic creation error: %s", err)
		return fmt.Errorf("topic creation error: %w", err)
	}

	err = msg.Send(ctx, t.Topic)
	if err != nil {
		log.Debugf("message publishing error: %s", err)
		return fmt.Errorf("message publishing error: %w", err)
	}
	log.Debugf("Message published to topic: %s", t.Topic.String())

	return nil
}
