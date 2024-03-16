package ui

import (
	"context"
	"strings"

	"github.com/bahner/go-ma-actor/p2p/peer"

	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

const (
	peerUsage        = "/peer connect|find|nick|delete|show"
	peerHelp         = "Manages peer info"
	peerShowUsage    = "/peer show <id|nick>"
	peerShowHelp     = "Shows the peer info"
	peerConnectUsage = "/peer connect <id|nick>"
	peerConnectHelp  = "Connects to a peer"
	peerFindUsage    = "/peer find id"
	peerFindHelp     = "Looks up a host in the distributed hash tables\nThis might take a while."
	peerDeleteUsage  = "/peer delete <id|nick>"
	peerDeleteHelp   = "Deletes a peer from the database, but not from the network"
	peerListUsage    = "/peer nick list"
	peerListHelp     = "List peer DID and nicks"
	peerNickUsage    = "/peer nick set <id|nick> <nick>"
	peerNickHelp     = "Sets a nick for an peer"
)

func (ui *ChatUI) handlePeerCommand(args []string) {

	if len(args) >= 2 {
		command := args[1]
		switch command {
		case "connect":
			ui.handlePeerConnectCommand(args)
			return
		case "delete":
			ui.handlePeerDeleteCommand(args)
			return
		case "find":
			ui.handlePeerFindCommand(args)
			return
		case "list":
			ui.handlePeerNickListCommand(args)
			return
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

	ui.handleHelpCommand(peerUsage, peerHelp)
}

// SHOW
func (ui *ChatUI) handlePeerShowCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		id, err := peer.LookupID(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		nick, err := peer.LookupNick(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		pid, err := p2peer.Decode(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		ui.displaySystemMessage("ID: " + id)
		ui.displaySystemMessage("Nick: " + nick)
		ui.displaySystemMessage("Maddrs:")
		for _, maddr := range ui.p.Host.Peerstore().Addrs(pid) {
			ui.displaySystemMessage(maddr.String())
		}
	} else {
		ui.handleHelpCommand(peerShowUsage, peerShowHelp)
	}
}

// LIST
func (ui *ChatUI) handlePeerNickListCommand(args []string) {

	log.Debugf("peer list command: %v", args)
	if len(args) == 2 {

		nicks := peer.Nicks()

		if len(nicks) > 0 {
			for k, v := range nicks {
				ui.displaySystemMessage(k + aliasSeparator + v)
			}
		} else {
			ui.displaySystemMessage("No peers found")
		}
	} else {
		ui.handleHelpCommand(peerListUsage, peerListHelp)
	}
}

// SET
func (ui ChatUI) handlePeerNickCommand(args []string) {

	if len(args) >= 4 {

		id := peer.Lookup(args[2])
		nick := strings.Join(args[3:], separator)

		p, err := peer.GetOrCreate(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}

		err = p.SetNick(nick)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		ui.displaySystemMessage(id + " is now known as " + nick)
	} else {
		ui.handleHelpCommand(peerNickUsage, peerNickHelp)
		return
	}
}

func (ui *ChatUI) handlePeerConnectCommand(args []string) {

	if len(args) == 3 {
		id := strings.Join(args[2:], separator)
		id, err := peer.LookupID(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}

		addrInfo, err := peer.PeerAddrInfoFromPeerIDString(ui.p.Host, id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		err = peer.ConnectAndProtect(context.Background(), ui.p.Host, addrInfo)
		if err != nil {
			ui.displaySystemMessage("Error connecting to peer: " + err.Error())
			return
		}
		ui.displaySystemMessage("Connected to " + id)
	} else {
		ui.handleHelpCommand(peerConnectUsage, peerConnectHelp)
	}
}

func (ui *ChatUI) handlePeerFindCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		id, err := peer.LookupID(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		pid, err := p2peer.Decode(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}

		ai, err := ui.p.DHT.FindPeer(context.Background(), pid)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		ui.displaySystemMessage("ID: " + ai.ID.String())
		for _, maddr := range ai.Addrs {
			ui.displaySystemMessage(maddr.String())
		}
	} else {
		ui.handleHelpCommand(peerFindUsage, peerFindHelp)
	}
}

func (ui *ChatUI) handlePeerDeleteCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		id, err := peer.LookupID(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		err = peer.Delete(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}

		ui.displaySystemMessage("Peer " + id + " deleted")
	} else {
		ui.handleHelpCommand(peerDeleteUsage, peerDeleteHelp)
	}
}
