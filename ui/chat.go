package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	actormsg "github.com/bahner/go-ma-actor/msg"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
// func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	e, err := entity.GetOrCreate(cm.From)
	if err != nil {
		log.Debugf("entity lookup error: %s", err)
		return
	}
	prompt := withColor("black", fmt.Sprintf("<%s>:", e.Nick()))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displayBroadcastMessage(cm *msg.Message) {
	e, err := entity.GetOrCreate(cm.From)
	if err != nil {
		log.Debugf("entity lookup error: %s", err)
		return
	}
	prompt := withColor("blue", fmt.Sprintf("<%s>:", e.Nick()))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displayPrivateMessage(cm *msg.Message) {
	e, err := entity.GetOrCreate(cm.From)
	if err != nil {
		log.Debugf("entity lookup error: %s", err)
		return
	}
	prompt := withColor("green", fmt.Sprintf("<%s>:", e.Nick()))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) displaySentPrivateMessage(cm *msg.Message) {
	e, err := entity.GetOrCreate(cm.To)
	if err != nil {
		log.Debugf("entity lookup error: %s", err)
		return
	}
	prompt := withColor("green", fmt.Sprintf(".%s:", e.Nick()))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, string(cm.Content))
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	// This could be a timeout for topic publishing
	ctx := context.Background()

	if ui.a == nil {
		ui.displaySystemMessage(errYouDontExist.Error())
		return errYouDontExist
	}

	if ui.e == nil {
		ui.displaySystemMessage(errNoEntitySelected.Error())
		return errNoEntitySelected
	}

	from := ui.a.Entity.DID.Id
	to := ui.e.DID.Id

	log.Debugf("Handling chatMessage: %s, from %s to %s", input, from, to)
	msgBytes := []byte(input)

	msg, err := actormsg.Chat(from, to, msgBytes, ui.a.Keyset.SigningKey.PrivKey)
	if err != nil {
		log.Debugf("failed to create chat message: %s", err)
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	err = msg.Send(ctx, ui.e.Topic)
	if err != nil {
		log.Debugf("message publishing error: %s", err)
		return fmt.Errorf("message publishing error: %w", err)
	}
	log.Debugf("Message published to topic: %s", ui.e.Topic.String())

	return nil
}
