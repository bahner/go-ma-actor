package ui

import (
	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/peer"
)

// handleAliasCommand handles the /alias command
func (ui *ChatUI) handleAliasCommand(args []string) {

	if len(args) > 1 {
		switch args[1] {
		case "node":
			ui.handleNodeAliasCommand(args)
		case "entity":
			ui.handleEntityAliasCommand(args)
		case "list":
			ui.handleAliasesCommand(args)
		default:
			ui.displaySystemMessage("Unknown alias type: " + args[1])
		}
	} else {
		ui.displaySystemMessage("Usage: /alias [node|entity|list]")
	}
}

func (ui *ChatUI) handleAliasesCommand(args []string) {

	ui.displaySystemMessage(alias.PrintEntityAliases())
	ui.displaySystemMessage(alias.PrintNodeAliases())

}

func (ui *ChatUI) handleEntityAliasCommand(args []string) {

	if len(args) > 2 {
		switch args[2] {
		case "add":
			ui.handleEntityAliasAddCommand(args)
		case "remove":
			ui.handleEntityAliasRemoveCommand(args)
		case "list":
			ui.displaySystemMessage(alias.PrintEntityAliases())
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + args[2])
		}
	}

}

func (ui *ChatUI) handleEntityAliasAddCommand(args []string) {

	if len(args) > 3 {
		alias.AddEntityAlias(args[3], args[4])
	}

}

func (ui *ChatUI) handleEntityAliasRemoveCommand(args []string) {

	if len(args) > 3 {
		alias.RemoveEntityAlias(args[3])
	}

}

func (ui *ChatUI) handleNodeAliasCommand(args []string) {

	if len(args) > 2 {
		switch args[2] {
		case "add":
			ui.handleNodeAliasAddCommand(args)
		case "remove":
			ui.handleNodeAliasRemoveCommand(args)
		case "list":
			ui.displaySystemMessage(alias.PrintNodeAliases())
		default:
			ui.displaySystemMessage("Unknown alias node command: " + args[2])
		}
	}

}

func (ui *ChatUI) handleNodeAliasAddCommand(args []string) {

	ap := peer.GetByAlias(args[3])
	if ap == nil {
		ui.displaySystemMessage("Unknown node with alias: " + args[3])
		return
	}
	if len(args) > 3 {
		alias.AddNodeAlias(ap.ID, args[4])
		ap.Alias = args[4]
	}

}

func (ui *ChatUI) handleNodeAliasRemoveCommand(args []string) {

	if len(args) > 3 {
		alias.RemoveNodeAlias(args[3])
	}

}
