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
			ui.handleHelpBroadcastCommand()
		case "discover":
			ui.handleHelpDiscoverCommand()
		case "enter":
			ui.handleHelpEnterCommand()
		case "entity":
			ui.handleHelpEntityCommands()
		case "me":
			ui.handleHelpMeCommands()
		case "msg":
			ui.handleHelpMsgCommand()
		case "peer":
			ui.handleHelpPeerCommands()
		case "refresh":
			ui.handleHelpRefreshCommand()
		case "resolve":
			ui.handleHelpResolveCommand()
		case "set":
			ui.handleHelpSetCommands()
		case "status":
			ui.handleHelpStatusCommands()
		case "whereis":
			ui.handleHelpWhereisCommand()
		case "quit":
			ui.handleHelpQuitCommand()
		default:
			ui.handleHelpUnknownCommand(args)
		}
	}
}

func (ui *ChatUI) handleHelpQuitCommand() {
	ui.displaySystemMessage("Usage: /quit")
	ui.displaySystemMessage("Quits the chat client")
}

func (ui *ChatUI) handleHelpUnknownCommand(args []string) {
	ui.displaySystemMessage("Unknown command: " + args[1])
}
