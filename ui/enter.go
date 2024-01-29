package ui

import "github.com/bahner/go-ma/did/doc"

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) > 1 {

		// If the DID is our own identity that is already handled.
		if args[1] == ui.a.Entity.DID.String() {
			ui.displaySystemMessage("You can't enter yourself. You are already here.")
			return
		}

		// If this is not the same as the last known location, then
		// update the last known location
		if ui.a.Entity.Doc.LastKnownLocation != args[1] {
			am, err := ui.a.Entity.Doc.GetAssertionMethod()
			if err != nil {
				ui.displaySystemMessage("Error getting assertion method: " + err.Error())
				return
			}

			ui.a.Entity.Doc.SetLastKnowLocation(args[1])
			ui.a.Entity.Doc.UpdateVersion()
			ui.a.Entity.Doc.Sign(ui.a.Entity.Keyset.SigningKey, am)
			opts := doc.DefaultPublishOptions()
			opts.Force = true
			go ui.a.Entity.Doc.Publish(opts)
		}

		ui.changeTopic(args[1])
		ui.msgBox.SetTitle(ui.e.DID)
		ui.displaySystemMessage("Entered: " + args[1])
	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}
