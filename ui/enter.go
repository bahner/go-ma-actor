package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) == 2 {

		_did := args[1]

		// This function handles the verification of the entity
		err := ui.enterEntity(_did, false) // force = false
		if err != nil {
			ui.displaySystemMessage("Error entering entity: " + err.Error())
			return
		}

	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}

// This is *the* function that changes the entity. Do Everything™ here.
// Do *not* use this to change the actor.
// INput is the nick or DID of the entity.
func (ui *ChatUI) enterEntity(d string, force bool) error {

	err := ui.a.Verify()
	if err != nil {
		return err
	}

	// First lookup any possible alias for the entity
	d = alias.LookupEntityNick(d)
	me := ui.a.Entity.DID.Id

	if d == me {
		return ErrSelfEntryWarning
	}

	e, err := entity.GetOrCreate(d)
	// Without a valid entity, we can't do anything.
	if err != nil || e == nil || e.Verify() != nil {
		return err
	}

	// FIXEME: hm. Why not?
	// If this is not the same as the last known location, then
	// update the last known location
	if d == e.DID.Id && !force {
		return ErrAlreadyHereWarning
	}

	// Here we go. This is the real deal.
	// Cancel the old entity.
	if ui.currentEntityCancel != nil {
		log.Debugf("Cancelling old entity: %s", ui.e.DID.Id)
		// Cancel the old entity
		ui.currentEntityCancel()
	}

	// set the new entity
	ui.e = e

	// Set the new entity context.
	log.Debugf("Setting new entity context for %s", e.DID.Id)
	ui.currentEntityCtx, ui.currentEntityCancel = context.WithCancel(context.Background())

	entityNick := alias.LookupEntityNick(e.DID.Id)
	if entityNick != e.DID.Id {
		log.Debugf("Changing entity nick from %s to %s", e.DID.Id, entityNick)
	}
	ui.msgBox.Clear()
	ui.msgBox.SetTitle(entityNick)

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	// Let the actor subscribe to the new entity, so
	// that envelopes are passed on correctly.
	go ui.a.Subscribe(ui.currentEntityCtx, ui.e)
	// Handle incoming envelopes to the entity as the actor.
	// Only an actor can decrypt and handle envelopes.
	go ui.handleIncomingEnvelopes(ui.currentEntityCtx, ui.e, ui.a)
	// Handle incoming messages to the entity
	go ui.handleIncomingMessages(ui.currentEntityCtx, ui.e)

	// Update the location
	// If this fails - 🤷🏽
	go ui.a.UpdateLastKnowLocation(e.DID.Id)

	return nil
}

func (ui *ChatUI) handleHelpEnterCommand(args []string) {
	ui.displaySystemMessage("Usage: /enter <DID>")
	ui.displaySystemMessage("Enters a chat room with the specified DID")
}
