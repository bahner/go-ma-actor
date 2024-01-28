package ui

func (ui *ChatUI) handleEnterCommand(args []string) {
	if len(args) > 1 {
		ui.changeTopic(args[1])
		ui.msgBox.SetTitle(ui.e.DID)
		ui.displaySystemMessage("Entered: " + args[1])
	} else {
		ui.displaySystemMessage("Usage: /enter [new_topic_name]")
	}
}
