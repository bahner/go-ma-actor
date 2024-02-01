package ui

import (
	"strings"

	"github.com/bahner/go-ma-actor/alias"
)

func (ui *ChatUI) handleCommands(input string) {
	args := strings.Split(input, " ")

	// Update alias when a command is entered
	ui.a.Nick = alias.Nick(ui.a.DID.String())
	ui.e.Nick = alias.Nick(ui.e.DID.String())

	switch args[0] {
	case "/status":
		ui.handleStatusCommand(args)
	case "/msg":
		ui.handleMsgCommand(args)
	case "/discover":
		ui.triggerDiscovery()
	case "/enter":
		ui.handleEnterCommand(args)
	case "/alias":
		ui.handleAliasCommand(args)
	case "/aliases":
		ui.handleAliasesCommand(args)
	case "/whereis":
		ui.handleWhereisCommand(args)
	case "discover":
		ui.triggerDiscovery()
	case "/refresh":
		ui.app.Draw()
	default:
		ui.displaySystemMessage("Unknown command: " + args[0])
	}
}
