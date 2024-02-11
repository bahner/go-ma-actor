package ui

import (
	"context"
	"fmt"

	"errors"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
)

var (
	errSelfEntryError = errors.New("you can't enter yourself")
	// errAlreadyHereError = errors.New("you are already here")
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) == 2 {

		_did := args[1]

		// This function handles the verification of the entity
		err := ui.enterEntity(_did)
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
func (ui *ChatUI) enterEntity(d string) error {

	// First lookup any possible alias for the entity
	d = alias.LookupEntityNick(d)

	e, err := entity.GetOrCreate(d, false)
	// Without a valid entity, we can't do anything.
	if err != nil || e == nil || e.Verify() != nil {
		return fmt.Errorf("enterEntity: failed to get or create entity: %v", err)
	}

	if d == ui.a.DID.String() {
		ui.displaySystemMessage(errSelfEntryError.Error())
		return errSelfEntryError
	}

	// // FIXEME: hm. Why not?
	// // If this is not the same as the last known location, then
	// // update the last known location
	// if d == e.DID.String() {
	// 	ui.displaySystemMessage(errAlreadyHereError.Error())
	// 	return errAlreadyHereError
	// }

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
	ui.e.Nick = alias.GetOrCreateEntityAlias(ui.e.DID.String())

	ui.msgBox.Clear()
	ui.msgBox.SetTitle(ui.e.Nick)

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	go ui.e.Subscribe(ui.currentEntityCtx, ui.e)
	go ui.handleIncomingMessages(ui.currentEntityCtx, ui.e)
	go ui.handleIncomingEnvelopes(ui.currentEntityCtx, ui.e)

	go ui.e.Subscribe(ui.currentEntityCtx, ui.a)
	go ui.handleIncomingEnvelopes(ui.currentEntityCtx, ui.a)

	// Update the location
	// If this fails - 🤷🏽
	go ui.a.UpdateLastKnowLocation(e.DID.String())

	return nil
}

func (ui *ChatUI) handleHelpEnterCommand(args []string) {
	ui.displaySystemMessage("Usage: /enter <DID>")
	ui.displaySystemMessage("Enters a chat room with the specified DID")
}
