package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

const (
	enterUsage = "/enter <DID>"
	enterHelp  = `Enters an entity with the specified DID
What this means is that messages will be sent to this entity.
Everybody 'in' the entity will be able to read the messages.
NB! use /msg to send encrypted messages to any recipient.`
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
		ui.displayHelpUsage(enterUsage)
	}
}

// This is *the* function that changes the entity. Do Everything™ here.
// Do *not* use this to change the actor.
// INput is the nick or DID of the entity.
func (ui *ChatUI) enterEntity(d string, reEntry bool) error {

	var (
		e   *entity.Entity
		err error
	)

	err = ui.a.Verify()
	if err != nil {
		return err
	}

	// If we have a cached entity for this nick, use it.
	e, err = entity.Lookup(d)
	if err != nil {
		// If we don't have it stored, then create it.
		e, err = entity.GetOrCreate(d)
	}

	// Without a valid entity, we can't do anything.
	if err != nil || e.Verify() != nil {
		return err
	}

	// FIXEME: hm. Why not?
	// If this is not the same as the last known location, then
	// update the last known location
	if ui.e != nil && d == ui.e.DID.Id && !reEntry {
		return errAlreadyHereWarning
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

	ui.msgBox.Clear()
	titleNick, err := entity.LookupNick(e.DID.Id)
	if err != nil {
		titleNick = e.DID.Id
		ui.displaySystemMessage("Error looking up nick: " + err.Error())
	}
	ui.msgBox.SetTitle(titleNick)

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	// And only unless we're entering self. Then there's no need. It's already running.
	if ui.e.DID.Id != ui.a.Entity.DID.Id {
		// Let the actor subscribe any new entity, so
		// that envelopes are passed on correctly.
		go ui.a.Subscribe(ui.currentEntityCtx, ui.e)

		// Handle incoming envelopes to the entity as the actor.
		// Only an actor can decrypt and handle envelopes.

		// Don't listen for envelopes when entering self.
		go ui.handleIncomingEnvelopes(ui.currentEntityCtx, ui.e, ui.a)

		// Handle incoming messages to the entity, also accept messages from self.
		go ui.handleIncomingMessages(ui.currentEntityCtx, ui.e)

	}

	// Update the location
	// If this fails - 🤷🏽
	go ui.a.UpdateLastKnowLocation(e.DID.Id)

	return nil
}
