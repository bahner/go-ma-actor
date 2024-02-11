package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
)

func (ui *ChatUI) handleResolveCommand(args []string) {

	if len(args) == 2 {

		id := alias.LookupEntityNick(args[1])

		e, err := entity.GetOrCreate(id)
		if err != nil {
			ui.displaySystemMessage("Error fetching entity: " + err.Error())
			return
		}

		ui.displaySystemMessage("Resolving DID Document for " + e.DID.String() + "...")
		d, err := doc.Fetch(id, false)
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}

		ui.displaySystemMessage("Resolved DID Document for " + e.DID.String())
		e.Doc = d

	} else {
		ui.handleHelpResolveCommand(args)
	}

}

func (ui *ChatUI) handleHelpResolveCommand(args []string) {
	ui.displaySystemMessage("Usage: /resolve <DID|NICK>")
	ui.displaySystemMessage("Tries to resolve the most recent version of the DID Document for the given DID or NICK.")
}
