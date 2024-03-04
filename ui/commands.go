package ui

import (
	"encoding/csv"
	"strings"
)

const commandSeparator = ' ' // This is a rune

func (ui *ChatUI) handleCommands(input string) {
	args, err := parseCommandsInput(input)
	if err != nil {
		ui.displaySystemMessage("Error parsing input: " + err.Error())
		return
	}

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
		ui.handleEditCommand()
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

// Takes the input and returns a slice of strings. This is used to split the input
// into a command and its arguments. Where "The Bar™" is considered a single argument,
func parseCommandsInput(input string) ([]string, error) {

	reader := csv.NewReader(strings.NewReader(input))
	// Set the delimiter to space
	reader.Comma = commandSeparator
	// Consider quotes as optional for fields
	reader.LazyQuotes = true
	// Read one line of input
	commands, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return deleteEmpties(commands), nil
}

func deleteEmpties(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
