package ui

import (
	"github.com/bahner/go-ma-actor/config"
)

const (
	saveUsage = "/save"
	saveHelp  = "Saves the current configuration to file"
)

func (ui *ChatUI) handleSaveCommand(args []string) {

	if len(args) == 1 {
		err := config.Save()
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		ui.displaySystemMessage("Configuration saved")
	}

	ui.handleHelpCommand(saveUsage, saveHelp)

}
