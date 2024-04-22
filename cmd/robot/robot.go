package main

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// A simple robot that can reply to messages.
type RobotStruct struct {

	// The location of the robot, where messages are sent to.
	Location *entity.Entity
	// The actor we send messages from.
	Robot *actor.Actor
	// Messages chan
	Messages chan *entity.Message
}

func NewRobot() (i *RobotStruct, err error) {

	i = &RobotStruct{}

	// Init of actor requires P2P to be initialized
	i.Robot = actor.Init()
	if i.Robot == nil {
		return nil, fmt.Errorf("failed to initialize actor")
	}

	messages := i.Robot.Entity.Messages

	// Subscribe to messages to self
	go i.Robot.Subscribe(context.Background(), i.Robot.Entity)
	go i.Robot.HandleIncomingEnvelopes(context.Background(), messages)
	go i.Robot.Entity.HandleIncomingMessages(context.Background(), messages)

	// Subscribe to message at location
	i.Location, err = entity.GetOrCreate(config.ActorLocation())
	if err != nil {
		return nil, fmt.Errorf("failed to get or create actor location: %w", errors.Cause(err))
	}

	go i.Robot.Subscribe(context.Background(), i.Location)
	go i.Robot.HandleIncomingEnvelopes(context.Background(), messages)

	go i.Location.HandleIncomingMessages(context.Background(), messages)

	go i.handleEntityMessageEvents()

	return i, err
}

func (i *RobotStruct) handleEntityMessageEvents() {
	ctx := context.Background()
	me := i.Robot.Entity.DID.Id
	myMessages := i.Robot.Entity.Messages
	errPrefix := fmt.Sprintf("handleEntityMessageEvents (%s): ", me)

	log.Debugf("Starting handleMessageEvents for %s", me)

	for {
		select {
		case <-ctx.Done(): // Check for cancellation signal
			log.Info(errPrefix + "context cancelled, exiting...")
			return

		case m, ok := <-myMessages: // Attempt to receive a message
			if !ok {
				log.Debugf(errPrefix + "channel closed, exiting...")
				return
			}

			if m == nil {
				log.Debugf(errPrefix + "received nil message, ignoring...")
				continue
			}

			if m.Message.Verify() != nil {
				log.Debugf(errPrefix+"failed to verify message: %v", m)
				continue
			}

			content := string(m.Message.Content)
			from := m.Message.From
			to := m.Message.To

			log.Debugf(errPrefix+"Handling message: %v from %s to %s", content, from, to)

			if from == me {
				log.Debugf(errPrefix + "Received message from self, ignoring...")
				continue
			}

			if m.Message.Type == msg.CHAT {
				i.handleMessage(ctx, m)
			}
		}
	}
}

func (i *RobotStruct) handleMessage(ctx context.Context, m *entity.Message) error {

	// Switch sender and receiver. Reply back to from :-)
	replyTo, err := entity.GetOrCreate(m.Message.From)
	if err != nil {
		return fmt.Errorf("failed to create new entity: %w", errors.Cause(err))
	}
	replyToEntity, err := entity.GetOrCreate(replyTo.DID.Id)
	if err != nil {
		log.Errorf("failed to create new entity: %v", err)
	}

	replyBytes := reply(m)

	err = m.Message.Reply(ctx, replyBytes, i.Robot.Keyset.SigningKey.PrivKey, replyToEntity.Topic)
	if err != nil {
		log.Errorf("failed to reply to message: %v", err)
	}

	return err
}

func reply(m *entity.Message) []byte {

	c := client()
	content := string(m.Message.Content)

	ctx := context.Background()
	res, err := c.SimpleSend(ctx, content)
	if err != nil {
		log.Fatal(err)
	}

	return []byte(res.Choices[0].Message.Content)
}
