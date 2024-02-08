package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleWhereisCommand(args []string) {

	if len(args) == 2 {

		// We have to do both here, in order not to make assumptions about the input
		id := alias.LookupEntityNick(args[1]) // This should return a DID
		nick := alias.GetOrCreateEntityAlias(id)

		if !did.IsValidDID(id) {
			ui.displaySystemMessage("Invalid DID: " + id)
			return
		}

		d, err := doc.Fetch(id, false) // Don't accept cached version
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}
		if d.LastKnownLocation != "" {
			ui.displaySystemMessage("Last known location of " + nick + " is " + d.LastKnownLocation)
		} else {
			ui.displaySystemMessage("No last known location for '" + nick + "(" + id + ")" + "' found")
		}
	} else {
		ui.displaySystemMessage("Usage: /whereis <DID>")
	}

}

func (ui *ChatUI) handleHelpWhereisCommand(args []string) {
	ui.displaySystemMessage("Usage: /whereis <DID>")
	ui.displaySystemMessage("Shows the last known location of a DID")
}
