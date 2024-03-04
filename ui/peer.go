package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/peer"
	log "github.com/sirupsen/logrus"
)

// PEER
func (ui *ChatUI) handleHelpPeerCommands() {
	ui.displaySystemMessage("Usage: /peer show|nick")
	ui.displaySystemMessage("Manages peer info")
}

func (ui *ChatUI) handlePeerCommand(args []string) {

	if len(args) >= 2 {
		command := args[1]
		switch command {
		case "nick":
			ui.handlePeerNickCommand(args)
			return
		case "show":
			ui.handlePeerShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown peer command: " + command)
		}
	}

	ui.handleHelpPeerCommands()
}

// SHOW
func (ui *ChatUI) handlePeerShowCommand(args []string) {

	if len(args) == 3 {
		id := args[2]
		p, err := peer.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		ui.displaySystemMessage("ID: " + p.ID)
		ui.displaySystemMessage("Nick: " + p.Nick)
		ui.displaySystemMessage("Maddrs:")
		for _, maddr := range p.AddrInfo.Addrs {
			ui.displaySystemMessage(maddr.String())
		}
	} else {
		ui.handleHelpPeerShowCommand()
	}
}

func (ui *ChatUI) handleHelpPeerShowCommand() {
	ui.displaySystemMessage("Usage: /peer show <id|nick>")
	ui.displaySystemMessage("       Shows the peer info")
}

// NICK
func (ui *ChatUI) handleHelpPeerNickCommand() {
	ui.displaySystemMessage("Usage: /peer nick list|set|show")
	ui.displaySystemMessage("       Manages peer nicks")
}

func (ui *ChatUI) handlePeerNickCommand(args []string) {

	if len(args) >= 3 {
		command := args[2]
		switch command {
		case "list":
			ui.handlePeerNickListCommand(args)
			return
		case "set":
			ui.handlePeerNickSetCommand(args)
			return
		case "show":
			ui.handlePeerNickShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias peer command: " + command)
		}
	}

	ui.handleHelpPeerNickCommand()
}

// LIST
func (ui *ChatUI) handlePeerNickListCommand(args []string) {

	log.Debugf("peer list command: %v", args)
	if len(args) == 3 {

		peers, err := peer.Peers()
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		log.Debugf("peers: %v", peers)

		if len(peers) > 0 {
			for _, v := range peers {
				ui.displaySystemMessage(v.ID + " : " + v.Nick)
			}
		} else {
			ui.displaySystemMessage("No peers found")
		}
	} else {
		ui.handleHelpPeerNickListCommand()
	}
}

func (ui *ChatUI) handleHelpPeerNickListCommand() {
	ui.displaySystemMessage("Usage: /peer nick list")
	ui.displaySystemMessage("List peer DID and nicks")
}

// SET
func (ui ChatUI) handlePeerNickSetCommand(args []string) {

	if len(args) == 5 {
		id := args[3]
		nick := args[4]
		p, err := peer.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		p.Nick = nick
		err = peer.Set(p)
		if err != nil {
			ui.displaySystemMessage("Error setting peer nick: " + err.Error())
			return
		}
		ui.displaySystemMessage(p.ID + " is now known as " + p.Nick)
	} else {
		ui.handleHelpPeerNickSetCommand()
		return
	}
}

func (ui *ChatUI) handleHelpPeerNickSetCommand() {
	ui.displaySystemMessage("Usage: /peer nick set <id|nick> <nick>")
	ui.displaySystemMessage("       Sets a nick for an peer")
}

// SHOW
func (ui *ChatUI) handlePeerNickShowCommand(args []string) {

	if len(args) == 4 {
		id := args[3]
		p, err := peer.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		peerInfo := fmt.Sprintf(p.ID + " is also known as " + p.Nick)
		ui.displaySystemMessage(peerInfo)
	} else {
		ui.handleHelpPeerNickShowCommand()
	}
}

func (ui *ChatUI) handleHelpPeerNickShowCommand() {
	ui.displaySystemMessage("Usage: /peer nick show <id|nick>")
	ui.displaySystemMessage("       Shows the peer info")
}
