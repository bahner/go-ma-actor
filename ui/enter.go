package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/did"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) > 1 {

		// Look up the DID, if possible and convert it to a proper did.DID
		d, err := did.New(
			alias.GetEntityDID(args[1]),
		)
		if err != nil {
			ui.displaySystemMessage("Invalid DID: " + err.Error())
			return
		}
		log.Debugf("Trying to find: %s", d.String())

		// If the DID is our own identity that is already handled.
		if d.String() == ui.a.DID.String() {
			ui.displaySystemMessage("You can't enter yourself. You are already here.")
			return
		}

		// If this is not the same as the last known location, then
		// update the last known location
		if ui.a.Doc.LastKnownLocation == d.String() {
			ui.displaySystemMessage("You are already there.")
			return
		}

		// Update the UI
		ui.changeEntity(ui.e.DID.String())
		ui.msgBox.SetTitle(ui.e.Nick)
		ui.displaySystemMessage("Entered: " + d.String())

		// Update the location
		err = ui.a.UpdateLastKnowLocation(d.String())
		if err != nil {
			ui.displaySystemMessage("Error updating last known location: " + err.Error())
			return
		}

	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}
