package ui

import (
	"strings"
)

func (ui *ChatUI) handleCommands(input string) {
	input = strings.TrimSpace(input) // Clear the cruft, if any
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
	case "/peer":
		ui.handlePeerCommand(args)
	case "/entity":
		ui.handleEntityCommand(args)
	case "/whereis":
		ui.handleWhereisCommand(args)
	case "/me":
		ui.handleMeCommands(args)
	case "/refresh":
		ui.refreshPeers()
		ui.msgBox.Clear()
		ui.app.Draw()
	default:
		ui.displaySystemMessage("Unknown command: " + args[0])
	}
}
