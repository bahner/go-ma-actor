package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
)

const (
	resolveUsage = "/resolve <DID|NICK>"
	resolveHelp  = "Tries to resolve the most recent version of the DID Document for the given DID or NICK."
)

func (ui *ChatUI) handleResolveCommand(args []string) {

	if len(args) == 2 {

		id := args[1]
		e, err := entity.GetOrCreate(id)
		if err != nil {
			ui.displaySystemMessage("Error fetching entity: " + err.Error())
			return
		}

		ui.displaySystemMessage("Resolving DID Document for " + e.DID.Id + "...")
		d, c, err := doc.Fetch(id, false)
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}
		ui.displaySystemMessage("Resolved DID Document for " + e.DID.Id + " (CID: " + c.String() + ")")
		e.Doc = d

	} else {
		ui.handleHelpCommand(resolveUsage, resolveHelp)
	}

}
