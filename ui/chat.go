package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	from := alias.GetOrCreateEntityAlias(cm.From)
	prompt := withColor("black", fmt.Sprintf("<%s>:", from))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displayBroadcastMessage(cm *msg.Message) {
	from := alias.GetOrCreateEntityAlias(cm.From)
	prompt := withColor("blue", fmt.Sprintf("<%s>:", from))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displayPrivateMessage(cm *msg.Message) {
	from := alias.GetOrCreateEntityAlias(cm.From)
	prompt := withColor("green", fmt.Sprintf("<%s>:", from))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displaySentPrivateMessage(cm *msg.Message) {
	to := alias.GetOrCreateEntityAlias(cm.To)
	prompt := withColor("green", fmt.Sprintf("@%s:", to))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	// This could be a timeout for topic publishing
	ctx := context.Background()

	if ui.a == nil {
		ui.displaySystemMessage(ErrYouDontExist.Error())
		return ErrYouDontExist
	}

	if ui.e == nil {
		ui.displaySystemMessage(ErrNoEntitySelected.Error())
		return ErrNoEntitySelected
	}

	from := ui.a.Entity.DID.Id
	to := ui.e.DID.Id

	log.Debugf("Handling chatMessage: %s, from %s to %s", input, from, to)
	msgBytes := []byte(input)

	msg, err := msg.New(from, to, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
	if err != nil {
		log.Debugf("message creation error: %s", err)
		return fmt.Errorf("message creation error: %w", err)
	}

	err = msg.Send(ctx, ui.e.Topic)
	if err != nil {
		log.Debugf("message publishing error: %s", err)
		return fmt.Errorf("message publishing error: %w", err)
	}
	log.Debugf("Message published to topic: %s", ui.e.Topic.String())

	return nil
}
