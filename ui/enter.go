package ui

import (
	"context"

	"errors"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

var (
	errSelfEntryError   = errors.New("you can't enter yourself")
	errAlreadyHereError = errors.New("you are already here")
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) == 2 {

		_did := args[1]

		// Update the UI
		e, err := entity.GetOrCreate(_did)
		if err != nil {
			ui.displaySystemMessage("Error getting entity: " + err.Error())
			return
		}

		// This function handles the verification of the entity
		err = ui.enterEntity(e)
		if err != nil {
			ui.displaySystemMessage("Error entering entity: " + err.Error())
			return
		}

	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}

// This is *the* function that changes the entity. Do Everything™ here.
func (ui *ChatUI) enterEntity(e *entity.Entity) error {

	// First lookup any possible alias for the entity
	d := alias.LookupEntityNick(e.DID.String())
	log.Debugf("Alias for %s: %s", e.DID.String(), d)

	if d == ui.a.DID.String() {
		ui.displaySystemMessage(errSelfEntryError.Error())
		return errSelfEntryError
	}

	// FIXEME: hm. Why not?
	// If this is not the same as the last known location, then
	// update the last known location
	if d == e.DID.String() {
		ui.displaySystemMessage(errAlreadyHereError.Error())
		return errAlreadyHereError
	}

	// Here we go. This is the real deal.
	// Cancel the old entity.
	if ui.currentEntityCancel != nil {
		// Cancel the old entity
		ui.currentEntityCancel()
	}

	// set the new entity
	ui.e = e

	// Set the new entity context.
	ui.currentEntityCtx, ui.currentEntityCancel = context.WithCancel(context.Background())

	// Look up the nick for the entity. Just a nicety, really.
	ui.e.Nick = alias.LookupEntityNick(ui.e.DID.String())

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	// NB! This is just for the entity. No envelopes are being handled here.
	go ui.subscribeToEntityPubsubMessages(e)
	go ui.handleIncomingMessages(e)

	// Update the UI
	err := ui.enterEntity(e)
	if err != nil {
		ui.displaySystemMessage("Error changing entity: " + err.Error())
		return err
	}

	// Update the location
	// If this fails - 🤷🏽
	go ui.a.UpdateLastKnowLocation(e.DID.String())

	return nil
}
