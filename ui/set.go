package ui

import (
	"context"
	"strings"

	"github.com/bahner/go-ma-actor/p2p/peer"
	"github.com/spf13/viper"
)

const (
	setUsage          = "/set broadcast|discovery"
	setHelp           = "Toggles broadcast and peer discovery on and off"
	setBroadcastUsage = "/set broadcast on|off"
	setBroadcastHelp  = "Toggles broadcast messages on and off"
	setDiscoveryUsage = "/set discovery on|off"
	setDiscoveryHelp  = "Toggles the discovery loop on and off"
	setLocationUsage  = "/set location <DID|NICK>|here"
	setLocationHelp   = "Set your actor's location\nUse 'here' to set your current location"
	setNickUsage      = "/set nick <NICK>"
	setNickHelp       = "Set your actor's nickname"
)

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) >= 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		case "discovery":
			ui.handleSetDiscoveryCommand(args)
		case "nick":
			ui.handleSetNickCommand(args)
		case "location":
			ui.handleSetLocationCommand(args)
		}
	} else {
		ui.handleHelpCommand(setUsage, setHelp)
	}

}

func (ui *ChatUI) handleSetNickCommand(args []string) {

	if len(args) >= 3 {

		nick := strings.Join(args[2:], " ")

		viper.Set("actor.nick", nick)
		ui.inputField.SetLabel(nick + ": ")
		return
	}
	ui.handleHelpCommand(setNickUsage, setNickHelp)
}

func (ui *ChatUI) handleSetLocationCommand(args []string) {

	var err error

	if len(args) == 3 {
		location := args[2]
		if location == "here" {
			location = ui.e.DID.Id
		} else {
			location, err = peer.LookupID(location)
			if err != nil {
				ui.displaySystemMessage("Error: " + err.Error())
				return
			}
		}

		viper.Set("actor.location", location)
		return
	}
	ui.handleHelpCommand(setNickUsage, setNickHelp)
}

func (ui *ChatUI) handleSetDiscoveryCommand(args []string) {

	if len(args) == 3 {

		toggle := args[2]

		switch toggle {
		case "on":

			// Now we can start continuous discovery in the background.
			ui.discoveryLoopCtx, ui.discoveryLoopCancel = context.WithCancel(context.Background())
			go ui.p.DiscoveryLoop(context.Background())
			ui.displaySystemMessage("Discovery is on")
		case "off":
			if ui.discoveryLoopCancel != nil {
				ui.discoveryLoopCancel()
				ui.displaySystemMessage("Discovery is off")
			}
		default:
			ui.handleHelpCommand(setDiscoveryUsage, setDiscoveryHelp)
		}
	} else {
		ui.handleHelpCommand(setDiscoveryUsage, setDiscoveryHelp)
	}
}

func (ui *ChatUI) handleSetBroadcastCommand(args []string) {

	if len(args) == 3 {

		toggle := args[2]

		switch toggle {
		case "on":
			go ui.subscribeBroadcasts()
			ui.displaySystemMessage(broadcastOnText)
			return
		case "off":
			if ui.broadcastCancel != nil {
				ui.broadcastCancel()
				ui.displaySystemMessage(broadcastOffText)
				return
			}
		default:
			ui.displayHelpUsage(setBroadcastUsage)
			return
		}
	}

	ui.handleHelpCommand(setBroadcastUsage, setBroadcastHelp)
}
