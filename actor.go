package main

import (
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	"github.com/bahner/go-ma/msg"
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
	// mouth chan *msg.Message // Messages to send
	ears chan *msg.Message // Received messages
	// hands chan string           // Local command input. Simple commands like /quit, /help etc.
}

func initActor(k *set.Keyset) (*Actor, error) {

	a := &Actor{}
	var err error

	log.Debugf("Setting Actor Keyset: %v", k)
	a.Keyset = k

	// Add the DID fragment as a field to the actor.
	a.DID, err = did.NewFromIPNSKey(a.Keyset.IPNSKey)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create DID: %v", err)
	}
	log.Debugf("new_actor: Created DID: %s", a.DID.String())

	// Publish the IPNSKey to IPFS for publication.
	err = k.IPNSKey.ExportToIPFS(a.DID.Fragment, *forcePublish)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to export IPNSKey to IPFS: %v", err)
	}

	// Make sure the actor has a DOC and published DIDDocument.
	a.Doc, err = doc.NewFromKeyset(a.Keyset, a.DID.String())
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to create DOC: %v", err)
	}

	_, err = a.Doc.Publish()
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to publish DOC: %v", err)
	}

	// We can now
	log.Debugf("new_actor: Joining to topic: %s", room)
	recvTopic, err := ps.Sub.Join(room)
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to join topic: %v", err)
	}

	log.Debugf("new_actor: Subscribing to topic: %s", room)
	a.From, err = recvTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("new_actor: Failed to subscribe to topic: %v", err)
	}

	log.Debugf("new_actor: Actor initialized: %s", a.DID.Fragment)
	return a, nil
}
