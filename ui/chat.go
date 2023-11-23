package ui

import (
	"fmt"

	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.From))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Body))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	log.Debugf("Handling chatMessage: %s", input)
	msgBytes := []byte(input)
	log.Debugf("ui.a.Entity.DID.Fragment: %s", ui.a.Entity.DID.Fragment)
	log.Debugf("ui.e.ID: %s", ui.e.DID)
	msg, err := msg.New(ui.a.Entity.DID.Fragment, ui.e.DID, msgBytes, "text/plain")
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
	ui.displayChatMessage(msg)

	// // FIXME. This should be done in the message.New function
	m, err := msg.MarshalToJSON()
	if err != nil {
		log.Debugf("message serialization error: %s", err)
		return fmt.Errorf("message serialization error: %s", err)
	}
	ui.displaySelfMessage(string(m))

	return nil
}
