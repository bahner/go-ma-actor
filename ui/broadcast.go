package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleBroadcastCommand(args []string) {

	if len(args) > 1 {

		me := ui.a.Entity.DID.Id

		var message string
		if len(args) > 2 {
			message = strings.Join(args[1:], " ")
		} else {
			message = args[1]
		}
		msgBytes := []byte(message)
		if log.GetLevel() == log.DebugLevel {
			ui.displaySystemMessage(fmt.Sprintf("Broadcasting %s", message))
		}

		msg, err := msg.NewBroadcast(me, msgBytes, "text/plain", ui.a.Keyset.SigningKey.PrivKey)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Broadcast creation error: %s", err))
		}

		err = msg.Broadcast(context.Background(), ui.b)
		if err != nil {
			ui.displaySystemMessage(fmt.Sprintf("Broadcast error: %s", err))
		}

		log.Debugf("Message broadcasted to topic: %s", ui.b)
	} else {
		ui.handleHelpBroadcastCommand(args)
	}

}

func (ui *ChatUI) handleHelpBroadcastCommand(args []string) {
	ui.displaySystemMessage("Usage: /broadcast <message>")
	ui.displaySystemMessage("Sends a public announcement to the current entity")
}

// This is *the* function that changes the entity. Do Everything™ here.
// Do *not* use this to change the actor.
// INput is the nick or DID of the entity.
func (ui *ChatUI) initBroadcast() error {

	var err error

	if ui.p == nil {
		return fmt.Errorf("initBroadcast: pubsub is nil")
	}

	ui.b, err = ui.p.PubSub.Join(ma.BROADCAST_TOPIC)
	if err != nil {
		return fmt.Errorf("initBroadcast: failed to join broadcast topic: %v", err)
	}

	// We don't want to respond to broadcasts in pong mode.
	// That would create a broadcast loop.
	if !config.PongMode() {
		go ui.subscribeBroadcasts()
	}

	return nil

}

// Subscribe to the entity's topic and handle incoming messages.
// the actor is the entity that will receive the messages.
// It may well be the entity itself.
// The context sho
func (ui *ChatUI) subscribeBroadcasts() {

	// We only ever want one subscription to the broadcast topic at a time
	if ui.broadcastCancel != nil {
		ui.broadcastCancel()
	}

	// This should be cancelled until, well - it's cancelled.
	ctx := context.Background()
	ui.broadcastCtx, ui.broadcastCancel = context.WithCancel(ctx)
	defer ui.broadcastCancel()

	sub, err := ui.b.Subscribe()
	if err != nil {
		log.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	messages := make(chan *p2ppubsub.Message, PUBSUB_MESSAGES_BUFFERSIZE)

	// Start an anonymous goroutine to bridge sub.Next() to the messages channel.
	go func() {
		for {
			// Assuming sub.Next() blocks until the next message is available,
			// and returns a message or an error.
			message, err := sub.Next(ctx)
			if err != nil {
				// Handle error (e.g., log, break, or continue based on the error type).
				log.Errorf("Error getting next message: %v", err)
				return // or continue based on your error handling policy
			}
			log.Debugf("handleSubscriptionMessages: Received message: %s", message.Data)

			// Assuming message is of the type you need; otherwise, adapt as necessary.
			select {
			case messages <- message:
			case <-ui.broadcastCtx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debugf("Broadcast reception cancelled")
			return
		case message, ok := <-messages:
			if !ok {
				log.Debugf("Broadcast channel closed, exiting...")
				return
			}

			m, err := msg.UnmarshalAndVerifyMessageFromCBOR(message.Data)
			if err == nil {
				if m.MimeType == ma.BROADCAST_MIME_TYPE {
					log.Debugf("handleBroadcastMessages: Received broadcast message: %v\n", m)
					ui.displayBroadcastMessage(m)
					continue
				}

				log.Error("handleBroadcastMessages: Received message is not a valid broadcast message")
			}
		}
	}
}

func (ui *ChatUI) handleSetBroadcastCommand(args []string) {

	if len(args) == 3 {

		toggle := args[2]

		switch toggle {
		case "on":
			ui.subscribeBroadcasts()
			ui.displaySystemMessage("Broadcasts are now on")
		case "off":
			if ui.broadcastCancel != nil {
				ui.broadcastCancel()
				ui.displaySystemMessage("Broadcasts are now off")
			}
		default:
			ui.handleHelpSetBroadcastCommand(args)
		}
	} else {
		ui.handleHelpSetBroadcastCommand(args)
	}
}

func (ui *ChatUI) handleHelpSetBroadcastCommand(args []string) {
	ui.displaySystemMessage("Usage: /set broadcast on|off")
	ui.displaySystemMessage("For now toggles broadcast messages on and off")
}
