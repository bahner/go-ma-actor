package ui

import (
	"github.com/bahner/go-ma/did/doc"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleWhereisCommand(args []string) {

	if len(args) == 2 {

		d, err := doc.Fetch(args[1], false) // Don't accept cached version
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}
		if d.LastKnownLocation != "" {
			ui.displaySystemMessage("Last known location of " + args[1] + " is " + d.LastKnownLocation)
		} else {
			ui.displaySystemMessage("No last known location for " + args[1])
		}
	} else {
		ui.displaySystemMessage("Usage: /whereis <DID>")
	}

}
