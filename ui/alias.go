package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	log "github.com/sirupsen/logrus"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleAliasCommand(args []string) {

	if len(args) > 1 {
		switch args[1] {
		case "node":
			ui.handleAliasNodeCommand(args)
		case "entity":
			ui.handleAliasEntityCommand(args)
		case "list":
			ui.handleAliasListCommand(args)
		default:
			ui.displaySystemMessage("Unknown alias type: " + args[1])
		}
	} else {
		ui.displaySystemMessage("Usage: /alias [node|entity|list]")
	}
}

func (ui *ChatUI) handleAliasListCommand(args []string) {

	ui.displaySystemMessage(alias.EntityAliases())
	ui.displaySystemMessage(alias.NodeAliases())

}

func (ui *ChatUI) handleAliasEntityCommand(args []string) {

	if len(args) > 2 {
		switch args[2] {
		case "set":
			ui.handleAliasEntitySetCommand(args)
		case "show":
			ui.handleAliasEntityShowCommand(args)
		case "remove":
			ui.handleAliasEntityRemoveCommand(args)
		case "list":
			ui.displaySystemMessage(alias.EntityAliases())
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + args[2])
		}
	}

}

func (ui *ChatUI) handleAliasEntitySetCommand(args []string) {

	// Attempt to see if the first param is an existing nick

	if len(args) == 5 {
		id := alias.LookupEntityNick(args[3])
		nick := args[4]
		alias.SetEntityAlias(id, nick)
		log.Debugf("Setting alias for %s to %s", id, nick)
	} else {
		ui.displaySystemMessage("Usage: /alias entity set <nick> <alias>")
	}

}

func (ui *ChatUI) handleAliasEntityRemoveCommand(args []string) {

	if len(args) == 4 {
		alias.RemoveEntityAlias(args[3])
	} else {
		ui.handleHelpAliasRemoveCommand()
	}

}

func (ui *ChatUI) handleHelpAliasRemoveCommand() {
	ui.displaySystemMessage("/alias entity|node remove <alias>")
}

func (ui *ChatUI) handleAliasNodeCommand(args []string) {

	if len(args) > 2 {
		switch args[2] {
		case "set":
			ui.handleAliasNodeSetCommand(args)
		case "show":
			ui.handleAliasNodeShowCommand(args)
		case "remove":
			ui.handleAliasNodeRemoveCommand(args)
		case "list":
			ui.displaySystemMessage(alias.NodeAliases())
		default:
			ui.displaySystemMessage("Unknown alias node command: " + args[2])
		}
	}

}

func (ui *ChatUI) handleAliasNodeSetCommand(args []string) {

	if len(args) == 5 {

		// Fetch the did if's referenced as an alias
		id := alias.LookupNodeAlias(args[3])

		alias.SetNodeAlias(id, args[4])

	} else {

		ui.displaySystemMessage("Usage: /alias node set <node_ID|node_ShortID> <alias>")
		return
	}

}

func (ui *ChatUI) handleAliasNodeRemoveCommand(args []string) {

	if len(args) > 3 {
		alias.RemoveNodeAlias(args[3])
	}

}

func (ui *ChatUI) handleAliasEntityShowCommand(args []string) {

	if len(args) == 4 {
		ui.displaySystemMessage(alias.GetOrCreateEntityAlias(args[3]))
	} else {
		ui.displaySystemMessage("Usage: /alias entity show <alias>")
	}

}

func (ui *ChatUI) handleAliasNodeShowCommand(args []string) {

	if len(args) == 4 {
		ui.displaySystemMessage(alias.GetOrCreateNodeAlias(args[3]))
	} else {
		ui.displaySystemMessage("Usage: /alias node show <alias>")
	}

}

func (ui *ChatUI) handleHelpAliasCommands(args []string) {
	ui.displaySystemMessage("Usage: /alias entity|node list|remove|show|set <DID|NICK> <alias>")
	ui.displaySystemMessage("Manages aliases for entities and nodes")
}

func (ui *ChatUI) handleHelpAliasesCommand(args []string) {
	ui.displaySystemMessage("Usage: /aliases")
	ui.displaySystemMessage("Lists all aliases")
}
