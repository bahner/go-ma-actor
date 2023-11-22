package actor

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/topic"
	"github.com/bahner/go-ma/entity"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"

	log "github.com/sirupsen/logrus"
)

const MESSAGES_BUFFERSIZE = 100

var err error

type Actor struct {

	// This context is used to cancel the Listen() function.
	ctx context.Context

	// All actors must be entities.
	// Ideally they should be the same, but then ma becomes a bit too opinionated.
	Entity *entity.Entity

	// Inbox is the topic where we receive envelopes from other actors.
	// It's basically a private channel with the DIDDocument keyAgreement as topic.
	Inbox *topic.Topic

	// Incoming messages from the actor to AssertionMethod topic. It's bascially a broadcast channel.
	// But you could use it to send messages to a specific actor or to all actors in a group.
	// This is a public channel. There will need to be some generic To (recipients) in the mesage
	// for example "broadcast", so that one actor can send a message to everybody in the room.
	// That is a TODO.
	// We receive the message contents here after verification or decryption.
	Messages chan *msg.Message
}

// Creates a new actor from an entity.
// Takes an entity and a forcePublish flag.
// The forcePublish is to override existing keys in IPFS.
func New(e *entity.Entity, forcePublish bool) (*Actor, error) {

	log.Debugf("actor.New: Setting Actor Entity: %v", e)

	a := new(Actor)

	// Firstly create assign entity to actor
	a.Entity = e

	// Create topic for incoming envelopes
	a.Inbox, err = topic.GetOrCreate(a.Entity.DID.String())
	if err != nil {
		if err.Error() != "topic already exists" {
			return nil, fmt.Errorf("actor.New: Failed to join topic: %w", err)
		}
	}

	// Set the messages channel
	a.Messages = make(chan *msg.Message, MESSAGES_BUFFERSIZE)

	// Publish the entity
	err = a.Entity.Publish(forcePublish)
	if err != nil {
		return nil, fmt.Errorf("actor.New: Failed to publish Entity: %w", err)
	}

	log.Debugf("actor.New: Actor initialized: %s", a.Entity.DID.Fragment)
	return a, nil

}

// Creates a new actor from a keyset.
// Takes a context, a keyset and a forcePublish flag.
// If ctx is nil a background context is used.
func NewFromKeyset(k *set.Keyset, forcePublish bool) (*Actor, error) {

	log.Debugf("Setting Actor Entity: %v", k)
	e, err := entity.NewFromKeyset(k)
	if err != nil {
		return nil, fmt.Errorf("actor.NewFromKeyset: Failed to create Entity: %w", err)
	}

	return New(e, forcePublish)
}

// Listen for incoming messages.
// The ctx is used to cancel the Listen() function.
func (a *Actor) Listen(ctx context.Context) {

	// We are the Cuckooes.
	if a.ctx != nil {
		a.ctx.Done()
	}

	a.ctx = ctx

	s, err := a.Inbox.Topic.Subscribe()
	if err != nil {
		log.Errorf("actor/listen: Failed to subscribe to topic: %v", err)
		return
	}

	// Listen for incoming messages
	go a.openEnvelopes(s)

}
