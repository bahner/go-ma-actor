package ui

import (
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) > 1 {

		// If the DID is our own identity that is already handled.
		if args[1] == ui.a.DID.String() {
			ui.displaySystemMessage("You can't enter yourself. You are already here.")
			return
		}

		// If this is not the same as the last known location, then
		// update the last known location
		if ui.a.Doc.LastKnownLocation != args[1] {
			am, err := ui.a.Doc.GetAssertionMethod()
			if err != nil {
				ui.displaySystemMessage("Error getting assertion method: " + err.Error())
				return
			}

			e, err := entity.GetOrCreate(args[1])
			if err != nil {
				ui.displaySystemMessage("Error getting entity: " + err.Error())
				return
			}

			// Now that we have the new entity we can cancel the old one
			ui.e.Topic.Subscription.Cancel()
			// And set the new one
			ui.e = e

			ui.a.Doc.SetLastKnowLocation(args[1])
			ui.a.Doc.UpdateVersion()
			ui.a.Doc.Sign(ui.a.Keyset.SigningKey, am)
			opts := doc.DefaultPublishOptions()
			opts.Force = true
			go ui.a.Doc.Publish(opts)
		}

		ui.changeEntity(args[1])
		ui.msgBox.SetTitle(ui.e.Nick)
		ui.displaySystemMessage("Entered: " + args[1])
	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}
