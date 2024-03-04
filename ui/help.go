package ui

import (
	"fmt"
	"strings"
)

const helpIndent = "       "

func (ui *ChatUI) displayHelpUsage(msg string) {
	out := withColor("purple", msg)
	fmt.Fprintf(ui.msgW, "Usage: %s\n", out)
}
func (ui *ChatUI) displayHelpText(msg string) {
	lines := strings.Split(msg, "\n")
	fmt.Fprintln(ui.msgW)
	for _, line := range lines {
		out := withColor("purple", line)
		fmt.Fprintf(ui.msgW, helpIndent+"%s\n", out)
	}
}

const (
	helpUsage = "/help [command]"
	helpText  = `Displays help for the given command, or a list of available commands if no command is given
NB! The input parser removes duplicate consecutive spaces and args.
    SO /entity nick set set FOO will be parsed as /entity nick set FOO
`
	quitUsage = "/quit"
	quitText  = "Quits the chat client"
)

func (ui *ChatUI) handleHelpCommands(args []string) {

	if len(args) == 1 {
		ui.displayHelpUsage(helpUsage)
		ui.displayHelpText(helpText)
		ui.displaySystemMessage("")
		ui.displaySystemMessage("Available commands:")
		ui.displaySystemMessage("/help broadcast")
		ui.displaySystemMessage("/help discover")
		ui.displaySystemMessage("/help enter")
		ui.displaySystemMessage("/help entity")
		ui.displaySystemMessage("/help me # Pun intended")
		ui.displaySystemMessage("/help msg")
		ui.displaySystemMessage("/help peer")
		ui.displaySystemMessage("/help quit")
		ui.displaySystemMessage("/help refresh")
		ui.displaySystemMessage("/help resolve")
		ui.displaySystemMessage("/help set")
		ui.displaySystemMessage("/help status")
		ui.displaySystemMessage("/help whereis")
		ui.displaySystemMessage("/help")
	} else {
		switch args[1] {
		case "broadcast":
			ui.handleHelpCommand(broadcastUsage, broadcastHelp)
		case "discover":
			ui.handleHelpCommand(discoverUsage, discoverHelp)
		case "enter":
			ui.handleHelpCommand(enterUsage, enterHelp)
		case "entity":
			ui.handleHelpCommand(entityUsage, entityHelp)
		case "me":
			ui.handleHelpCommand(meUsage, meHelp)
		case "msg":
			ui.handleHelpCommand(msgUsage, msgHelp)
		case "peer":
			ui.handleHelpCommand(peerUsage, peerHelp)
		case "refresh":
			ui.handleHelpCommand(refreshUsage, refreshHelp)
		case "resolve":
			ui.handleHelpCommand(resolveUsage, resolveHelp)
		case "set":
			ui.handleHelpCommand(setUsage, setHelp)
		case "status":
			ui.handleHelpCommand(statusUsage, statusHelp)
		case "whereis":
			ui.handleHelpCommand(whereisUsage, whereisHelp)
		case "quit":
			ui.handleHelpCommand(quitUsage, quitText)
		default:
			ui.handleHelpUnknownCommand(args)
		}
	}
}

func (ui *ChatUI) handleHelpUnknownCommand(args []string) {
	ui.displaySystemMessage("Unknown command: " + args[1])
}

func (ui *ChatUI) handleHelpCommand(usage string, help string) {
	ui.displayHelpUsage(usage)
	ui.displayHelpText(help)
}
