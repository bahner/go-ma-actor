package ui

import (
	"strings"

	"github.com/bahner/go-ma-actor/alias"
)

func (ui *ChatUI) handleCommands(input string) {
	args := strings.Split(input, " ")

	// Update alias when a command is entered
	ui.a.Nick = alias.GetOrCreateEntityAlias(ui.a.DID.String())
	ui.e.Nick = alias.GetOrCreateEntityAlias(ui.e.DID.String())

	switch args[0] {
	case "/help":
		ui.handleHelpCommands(args)
	case "/status":
		ui.handleStatusCommand(args)
	case "/msg":
		ui.handleMsgCommand(args)
	case "/broadcast":
		ui.handleBroadcastCommand(args)
	case "/discover":
		ui.triggerDiscovery()
	case "/enter":
		ui.handleEnterCommand(args)
	case "/alias":
		ui.handleAliasCommand(args)
	case "/aliases":
		ui.handleAliasListCommand(args)
	case "/whereis":
		ui.handleWhereisCommand(args)
	case "/me":
		ui.handleMeCommands(args)
	case "discover":
		ui.triggerDiscovery()
	case "/refresh":
		ui.refreshPeers()
	default:
		ui.displaySystemMessage("Unknown command: " + args[0])
	}
}
