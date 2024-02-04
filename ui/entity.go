package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) setEntity(did string) error {

	log.Debugf("Setting entity to %s", did)

	var err error

	log.Debugf("Creating entity for topic %s", did)
	// e, err = getOrCreateEntity(did)
	e, err := entity.GetOrCreate(did)
	if err != nil {
		return fmt.Errorf("error getting or creating entity: %w", err)
	}

	// Loog up the nick for the entity
	e.Nick = alias.LookupEntityDID(did)

	// Now pivot to the new entity
	// and cancel the old.
	// old_nick := ui.e
	ui.e = e
	// Leave the subscription to the old topic for now.
	// old_nick.Subscription.Cancel()

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	go ui.subscribeToEntityPubsubMessages(e)
	go ui.handleIncomingMessages(e)

	log.Infof("Entity set to %s", ui.e.Topic.String())

	return nil

}
