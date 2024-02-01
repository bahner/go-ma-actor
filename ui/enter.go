package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/did"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) == 2 {

		_did := args[1]
		// If id is not a valid did, then try to find it in the aliases
		if !did.IsValidDID(_did) {
			_did = alias.GetEntityDID(_did)
		}

		// If it is still not a valid did, then return
		if _did == "" {
			ui.displaySystemMessage("Invalid DID")
			return
		}

		log.Debugf("Trying to find: %s", _did)

		// If the DID is our own identity that is already handled.
		if _did == ui.a.DID.String() {
			ui.displaySystemMessage("You can't enter yourself.")
			return
		}

		// If this is not the same as the last known location, then
		// update the last known location
		if ui.e.DID.String() == _did {
			ui.displaySystemMessage("You are already here.")
			return
		}

		// Update the UI
		err := ui.changeEntity(_did)
		if err != nil {
			ui.displaySystemMessage("Error changing entity: " + err.Error())
			return
		}
		ui.msgBox.SetTitle(ui.e.Nick)
		ui.displaySystemMessage("Entered: " + _did)

		// Update the location
		err = ui.a.UpdateLastKnowLocation(_did)
		if err != nil {
			ui.displaySystemMessage("Error updating last known location: " + err.Error())
			return
		}

	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}
