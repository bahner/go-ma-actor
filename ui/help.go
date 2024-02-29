package ui

import "fmt"

func (ui *ChatUI) displayHelpUsage(msg string) {
	out := withColor("purple", msg)
	fmt.Fprintf(ui.msgW, "Usage: %s\n", out)
}
func (ui *ChatUI) displayHelpText(msg string) {
	out := withColor("purple", msg)
	fmt.Fprintf(ui.msgW, indent+"%s\n", out)
}

func (ui *ChatUI) handleHelpCommands(args []string) {

	if len(args) == 1 {
		ui.displaySystemMessage("Usage: /help [command]")
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
			ui.handleHelpBroadcastCommand(args)
		case "discover":
			ui.handleHelpDiscoverCommand(args)
		case "enter":
			ui.handleHelpEnterCommand(args)
		case "entity":
			ui.handleHelpEntityCommands(args)
		case "me":
			ui.handleHelpMeCommands(args)
		case "msg":
			ui.handleHelpMsgCommand(args)
		case "peer":
			ui.handleHelpPeerCommands(args)
		case "refresh":
			ui.handleHelpRefreshCommand(args)
		case "resolve":
			ui.handleHelpResolveCommand(args)
		case "set":
			ui.handleHelpSetCommands(args)
		case "status":
			ui.handleHelpStatusCommands(args)
		case "whereis":
			ui.handleHelpWhereisCommand(args)
		case "quit":
			ui.handleHelpQuitCommand(args)
		default:
			ui.handleHelpUnknownCommand(args)
		}
	}
}

func (ui *ChatUI) handleHelpQuitCommand(args []string) {
	ui.displaySystemMessage("Usage: /quit")
	ui.displaySystemMessage("Quits the chat client")
}

func (ui *ChatUI) handleHelpUnknownCommand(args []string) {
	ui.displaySystemMessage("Unknown command: " + args[1])
}
