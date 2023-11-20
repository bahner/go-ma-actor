package ui

func (ui *ChatUI) handleEnterCommand(args []string) {
	if len(args) > 1 {
		ui.displaySystemMessage("TODO: enter new topic")
	} else {
		ui.displaySystemMessage("Usage: /enter [new_topic_name]")
	}
}
