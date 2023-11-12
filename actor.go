package main

import (
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/message"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	log "github.com/sirupsen/logrus"
)

type Actor struct {
	Keyset *set.Keyset
	DID    *did.DID
	Doc    *doc.Document

	// PubSub Attributes
	From *pubsub.Subscription // The subscription to the topic for receiving messages.

	// By using antropomorphic terms, we underline that this is a special use case.
	// mouth chan *message.Message // Messages to send
	ears chan *message.Message // Received messages
	// hands chan string           // Local command input. Simple commands like /quit, /help etc.
}

func initActor(k *set.Keyset) (*Actor, error) {

	var a *Actor
	var err error

	// Add the DID fragment as a field to the actor.
	a.DID, err = did.NewFromIPNSKey(k.IPNSKey)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create DID: %v", err)
	}
	log.Debugf("new_actor: Created DID: %s", a.DID.String())

	// Make sure the actor has a DOC and published DIDDocument.
	a.Doc, err = doc.New(a.DID.String(), a.DID.String())
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create DOC: %v", err)
	}

	_, err = a.Doc.Publish()
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to publish DOC: %v", err)
	}

	// We can now
	recvTopic, err := ps.Sub.Join(node.Node.ID().String()) // The ipnskey is the id of the actor.
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to join topic: %v", err)
	}

	a.From, err = recvTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to subscribe to topic: %v", err)
	}

	return a, nil
}
