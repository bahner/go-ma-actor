package ui

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	setUsage = "/set broadcast|discovery"
	setHelp  = "Toggles broadcast and peer discovery on and off"
)

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) == 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		case "discovery":
			ui.handleSetDiscoveryCommand(args)
		case "nick":
			ui.handleSetNickCommand(args)
		}
	} else {
		ui.handleHelpCommand(setUsage, setHelp)
	}

}

const (
	setNickUsage = "/set nick <NICK>"
	setNickHelp  = "Set your actor's nickname"
)

func (ui *ChatUI) handleSetNickCommand(args []string) {

	if len(args) == 3 {

		nick := strings.Join(args[2:], " ")

		viper.Set("actor.nick", nick)
		ui.inputField.SetLabel(nick + ": ")
		return
	}
	ui.handleHelpCommand(setNickUsage, setNickHelp)
}
