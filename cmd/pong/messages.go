package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	actormsg "github.com/bahner/go-ma-actor/msg"
	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func handleMessageEvents(ctx context.Context, a *actor.Actor) {
	me := a.Entity.DID.Id

	log.Debugf("Starting handleMessageEvents for %s", me)

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info("handleMessageEvents: context cancelled, exiting...")
			return

		case m, ok := <-a.Entity.Messages: // Attempt to receive a message
			if !ok {
				log.Debugf("messageEvents: channel closed, exiting...")
				return
			}

			if m == nil {
				log.Debugf("messageEvents: received nil message, ignoring...")
				continue
			}

			if m.Message.Verify() != nil {
				log.Debugf("messageEvents: failed to verify message: %v", m)
				continue
			}

			if m.Message.From == me {
				log.Debugf("Received message from self, ignoring...")
				continue
			}

			messageType, err := m.Message.MessageType()
			if err != nil {
				log.Debugf("Failed to get message type: %v", err)
				continue
			}

			if messageType != actormsg.CHAT_MESSAGE_TYPE {
				log.Debugf("Received message of unknown type: %s, ignoring...", messageType)
				continue
			}

			if m.Message.To == me {
				messageReply(ctx, a, m.Message)
			}
		}
	}
}

func messageReply(ctx context.Context, a *actor.Actor, m *msg.Message) error {

	// Switch sender and receiver. Reply back to from :-)
	replyTo := m.From
	replyToEntity, err := entity.GetOrCreate(replyTo)
	if err != nil {
		log.Errorf("messageReply: %v", err)
		return err
	}

	reply := []byte(viper.GetString("pong.reply"))
	fmt.Printf("Sending private message to %s over %s\n", replyTo, a.Entity.Topic.String())

	err = actormsg.Reply(ctx, *m, reply, a.Keyset.SigningKey.PrivKey, replyToEntity.Topic)
	if err != nil {
		log.Errorf("messageReply: %v", err)
	}

	return err

}
