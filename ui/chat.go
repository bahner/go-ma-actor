package ui

import (
	"encoding/json"
	"fmt"

	"github.com/bahner/go-ma/msg"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.From))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, cm.Body)
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	msgBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("message serialization error: %s", err)
	}

	msg, err := msg.New(ui.nick, ui.nick, string(msgBytes), "application/json")
	if err != nil {
		return fmt.Errorf("message creation error: %s", err)
	}

	// FIXME. This should be done in the message.New function
	m, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("message serialization error: %s", err)
	}

	err = ui.a.Public.Publish(ui.ctx, m)
	if err != nil {
		return fmt.Errorf("publish error: %s", err)
	}
	ui.displaySelfMessage(input)

	return nil
}
