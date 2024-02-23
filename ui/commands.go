package ui

import (
	"strings"
)

func (ui *ChatUI) handleCommands(input string) {
	args := strings.Split(input, " ")

	switch args[0] {
	case "/help":
		ui.handleHelpCommands(args)
	case "/status":
		ui.handleStatusCommand(args)
	case "/msg":
		ui.handleMsgCommand(args)
	case "/broadcast":
		ui.handleBroadcastCommand(args)
	case "/set":
		ui.handleSetCommand(args)
	case "/edit":
		ui.handleEditCommand(args)
	case "/resolve":
		go ui.handleResolveCommand(args) // This make take some time. No need to block the UI
	case "/discover":
		ui.triggerDiscovery()
	case "/enter":
		ui.handleEnterCommand(args)
	case "/nick":
		ui.handleNickCommand(args)
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
