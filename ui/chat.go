package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
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
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Body))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	// This could be a timeout for topic publishing
	ctx := context.Background()

	log.Debugf("Handling chatMessage: %s", input)
	msgBytes := []byte(input)
	log.Debugf("ui.a.Entity.DID.Fragment: %s", ui.a.Entity.DID.Fragment)
	log.Debugf("ui.e.ID: %s", ui.e.DID)
	msg, err := msg.New(ui.a.Entity.DID.String(), ui.e.DID, msgBytes, "text/plain")
	if err != nil {
		log.Debugf("message creation error: %s", err)
		return fmt.Errorf("message creation error: %w", err)
	}

	log.Debugf("Signing message")
	err = msg.Sign(ui.a.Entity.Keyset.SigningKey.PrivKey)
	if err != nil {
		log.Debugf("message signing error: %s", err)
		return fmt.Errorf("message signing error: %w", err)
	}
	msgJson, _ := msg.MarshalToJSON()
	log.Debugf("Message signed: %s", msgJson)
	err = msg.VerifySignature()
	if err != nil {
		log.Debugf("failed to verify my own message: %s", err)
		return fmt.Errorf("message verification error: %w", err)
	} else {
		log.Debugf("Message signature verified")
	}
	ui.displayChatMessage(msg)
	topic, err := topic.GetOrCreate(ui.e.DID)
	if err != nil {
		log.Debugf("topic creation error: %s", err)
		return fmt.Errorf("topic creation error: %w", err)
	}

	e, err := msg.Enclose()
	if err != nil {
		return fmt.Errorf("envelope creation error: %w", err)
	}

	letter, err := e.MarshalToCBOR()
	if err != nil {
		log.Debugf("message serialization error: %s", err)
		return fmt.Errorf("message serialization error: %s", err)
	}

	err = topic.Topic.Publish(ctx, letter)
	if err != nil {
		log.Debugf("message publishing error: %s", err)
		return fmt.Errorf("message publishing error: %w", err)
	}
	log.Debugf("Message published to topic: %s", topic.Topic.String())

	// // FIXME. This should be done in the message.New function
	m, err := msg.MarshalToJSON()
	if err != nil {
		log.Debugf("message serialization error: %s", err)
		return fmt.Errorf("message serialization error: %s", err)
	}

	if config.GetLogLevel() == "debug" {
		ui.displaySelfMessage(string(m))
	}

	return nil
}
