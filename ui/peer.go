package ui

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/p2p/peer"

	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

const (
	peerUsage         = "/peer connect|find|nick|remove|show"
	peerHelp          = "Manages peer info"
	peerShowUsage     = "/peer show <id|nick>"
	peerShowHelp      = "Shows the peer info"
	peerConnectUsage  = "/peer connect <id|nick>"
	peerConnectHelp   = "Connects to a peer"
	peerFindUsage     = "/peer find id"
	peerFindHelp      = "Looks up a host in the distributed hash tables\nThis might take a while."
	peerDeleteUsage   = "/peer remove <id|nick>"
	peerDeleteHelp    = "Deletes a peer from the database, but not from the network"
	peerNickUsage     = "/peer nick list|set|show"
	peerNickHelp      = "Manages peer nicks"
	peerNickListUsage = "/peer nick list"
	peerNickListHelp  = "List peer DID and nicks"
	peerNickSetUsage  = "/peer nick set <id|nick> <nick>"
	peerNickSetHelp   = "Sets a nick for an peer"
	peerNickShowUsage = "/peer nick show <id|nick>"
	peerNickShowHelp  = "Shows the peer info"
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

	if len(args) == 3 {
		id := args[2]
		p, err := peer.Get(id)
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
		ui.handleHelpCommand(peerShowUsage, peerShowHelp)
	}
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

	ui.handleHelpCommand(peerNickUsage, peerNickHelp)
}

// LIST
func (ui *ChatUI) handlePeerNickListCommand(args []string) {

	log.Debugf("peer nick list command: %v", args)
	if len(args) == 3 {

		nicks := peer.Nicks()

		if len(nicks) > 0 {
			for k, v := range nicks {
				ui.displaySystemMessage(k + " : " + v)
			}
		} else {
			ui.displaySystemMessage("No peers found")
		}
	} else {
		ui.handleHelpCommand(peerNickListUsage, peerNickListHelp)
	}
}

// SET
func (ui ChatUI) handlePeerNickSetCommand(args []string) {

	if len(args) == 5 {
		id := args[3]
		nick := args[4]
		p, err := ui.p.GetOrCreatePeerFromIDString(peer.Lookup(id))
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
		ui.handleHelpCommand(peerNickSetUsage, peerNickSetHelp)
		return
	}
}

// SHOW
func (ui *ChatUI) handlePeerNickShowCommand(args []string) {

	if len(args) == 4 {
		id := args[3]
		p, err := peer.Get(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		peerInfo := fmt.Sprintf(p.ID + " is also known as " + p.Nick)
		ui.displaySystemMessage(peerInfo)
	} else {
		ui.handleHelpCommand(peerNickShowUsage, peerNickShowHelp)
	}
}

func (ui *ChatUI) handlePeerConnectCommand(args []string) {

	if len(args) == 3 {
		id := peer.Lookup(args[2])

		addrInfo, err := ui.p.GetPeerAddrInfoFromIDString(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		err = ui.p.DHT.PeerConnectAndUpdateIfSuccessful(context.Background(), addrInfo)
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

	if len(args) == 3 {
		id, err := p2peer.Decode(args[2])
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}

		ai, err := ui.p.DHT.FindPeer(context.Background(), id)
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

	if len(args) == 3 {
		id, err := peer.LookupID(args[2])
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
