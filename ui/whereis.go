package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleWhereisCommand(args []string) {

	if len(args) == 2 {

		// We have to do both here, in order not to make assumptions about the input
		e, err := entity.GetOrCreate(args[1])
		if err != nil {
			ui.displaySystemMessage("Failed to generate entity for request: " + err.Error())
			return
		}

		d, err := doc.Fetch(e.DID.Id, false) // Don't accept cached version
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}
		if d.LastKnownLocation != "" {
			ui.displaySystemMessage("Last known location of " + e.Nick + " is " + d.LastKnownLocation)
		} else {
			ui.displaySystemMessage("No last known location for '" + e.Nick + "(" + e.DID.Id + ")" + "' found")
		}
	} else {
		ui.displaySystemMessage("Usage: /whereis <DID>")
	}

}

func (ui *ChatUI) handleHelpWhereisCommand() {
	ui.displaySystemMessage("Usage: /whereis <DID>")
	ui.displaySystemMessage("Shows the last known location of a DID")
}
