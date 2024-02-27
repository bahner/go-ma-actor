package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/peer"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleHelpPeerCommands(args []string) {
	ui.displaySystemMessage("Usage: /peer remove|show|nick")
	ui.displaySystemMessage("Manages peer info")
}

func (ui *ChatUI) handlePeerCommand(args []string) {

	if len(args) > 2 {
		command := args[1]
		switch command {
		case "nick":
			ui.handlePeerNickCommand(args)
			return
		case "remove":
			ui.handlePeerRemoveCommand(args)
			return
		case "show":
			ui.handlePeerShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpPeerCommands(args)

}

func (ui *ChatUI) handleHelpPeerRemoveCommand() {
	ui.displaySystemMessage("/peer remove <id|nick>")
}

func (ui *ChatUI) handlePeerRemoveCommand(args []string) {

	if len(args) == 3 {
		err := peer.Remove(args[3])
		if err != nil {
			ui.displaySystemMessage("Error removing peer: " + err.Error())
			return
		}
	} else {
		ui.handleHelpPeerRemoveCommand()
	}

}

func (ui *ChatUI) handleHelpPeerNickCommand(args []string) {
	ui.displaySystemMessage("Usage: /peer nick <id> <nick>")
	ui.displaySystemMessage("Set a nick for a peer")
}

func (ui *ChatUI) handlePeerNickCommand(args []string) {

	// No nick given, hence just show the existing nick
	if len(args) == 3 {
		p, err := peer.Lookup(args[2])
		if err != nil {
			ui.displaySystemMessage("Error fetching alias: " + err.Error())
			return
		}
		log.Debugf("%s: %s", p, p.Nick)
		ui.displaySystemMessage(fmt.Sprintf("Alias for %s is set to %s", p.ID, p.Nick))
		return
	}

	if len(args) == 4 {
		p, err := peer.Lookup(args[2])
		if err != nil {
			ui.displaySystemMessage("Error fetching alias: " + err.Error())
			return
		}
		p.Nick = args[3]
		err = peer.Set(p)
		if err != nil {
			ui.displaySystemMessage("Error setting alias: " + err.Error())
			return
		}
		log.Debugf("Setting alias for %s to %s", p.ID, p.Nick)
		ui.displaySystemMessage(fmt.Sprintf("Alias for %s set to %s", p.ID, p.Nick))
		return
	}

	ui.handleHelpPeerNickCommand(args)

}

func (ui *ChatUI) handlePeerShowCommand(args []string) {

	if len(args) == 3 {
		id := args[2]
		val, err := peer.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		peerInfo := fmt.Sprintf("ID: %s\nNick: %s\nAddrs:", val.ID, val.Nick)
		ui.displaySystemMessage(peerInfo)
		for _, a := range val.AddrInfo.Addrs {
			ui.displaySystemMessage(a.String())
		}
	} else {
		ui.handleHelpPeerShowCommand(args)
	}

}

func (ui *ChatUI) handleHelpPeerShowCommand(args []string) {
	ui.displaySystemMessage("Usage: /peer show <id|nick>")
	ui.displaySystemMessage("Shows the peer info")
}
