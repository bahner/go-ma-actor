package ui

import "strings"

func (ui *ChatUI) handleCommands(input string) {
	args := strings.Split(input, " ")

	switch args[0] {
	case "/status":
		ui.handleStatusCommand(args)
	case "/discover":
		ui.triggerDiscovery()
	case "/enter":
		ui.handleEnterCommand(args)
	case "/alias":
		ui.handleAliasCommand(args)
	case "/refresh":
		ui.app.Draw()
	default:
		ui.displaySystemMessage("Unknown command: " + args[0])
	}
}
