package ui

import (
	"context"
	"strings"

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

	if len(args) >= 2 {

		id := strings.Join(args[1:], separator)
		id = entity.DID(id)

		e, err := entity.GetOrCreate(id, false)
		if err != nil {
			ui.displaySystemMessage("Error getting entity: " + err.Error())
			return
		}

		// This function handles the verification of the entity
		err = ui.enterEntity(e, true)
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
// Input is a did, so we know it's a valid entity.
func (ui *ChatUI) enterEntity(e *entity.Entity, reEntry bool) error {

	var err = e.Verify()

	if err != nil {
		log.Errorf("Error verifying entity for entry: %s", err.Error())
		// Without an entity, we can't do anything.
		return err
	}

	// Warn if we're already here and not re-entering.
	if ui.e != nil && e.DID.Id == ui.e.DID.Id && !reEntry {
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
	ui.msgBox.SetTitle(entity.Nick(e.DID.Id))

	// Start handling the new topic
	// This *must* be called *after* the entity is set!
	// And only unless we're entering self. Then there's no need. It's already running.
	if ui.e.DID.Id != ui.a.Entity.DID.Id {
		// Let the actor subscribe any new entity, so
		// that envelopes are passed on correctly.
		go ui.a.Subscribe(ui.currentEntityCtx, ui.e)

		// Don't listen for envelopes when entering self.
		go ui.a.HandleIncomingEnvelopes(ui.currentEntityCtx, ui.chMessages)

		// Handle incoming messages to the entity, also accept messages from self.
		go ui.e.HandleIncomingMessages(ui.currentEntityCtx, ui.chMessages)

	}

	return nil
}
