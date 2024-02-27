package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleHelpEntityCommands(args []string) {
	ui.displaySystemMessage("Usage: /entity remove|show|nick")
	ui.displaySystemMessage("Manages entity info")
}

func (ui *ChatUI) handleEntityCommand(args []string) {

	if len(args) >= 2 {
		command := args[1]
		switch command {
		case "list":
			ui.handleEntityListCommand(args)
			return
		case "nick":
			ui.handleEntityNickCommand(args)
			return
		case "remove":
			ui.handleEntityRemoveCommand(args)
			return
		case "show":
			ui.handleEntityShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpEntityCommands(args)

}

func (ui *ChatUI) handleEntityListCommand(args []string) {

	log.Debugf("entity list command: %v", args)
	if len(args) == 2 {

		nicks := entity.ListNicks()
		log.Debugf("entities: %v", nicks)

		if len(nicks) > 0 {
			for k, v := range nicks {
				ui.displaySystemMessage(fmt.Sprintf("%s: %s", k, v))
			}
		} else {
			ui.displaySystemMessage("No entities found")
		}
	} else {
		ui.handleHelpEntityListCommand(args)
	}

}

func (ui *ChatUI) handleHelpEntityListCommand(args []string) {
	ui.displaySystemMessage("Usage: /entity list")
	ui.displaySystemMessage("List entity DID and nicks")
}

func (ui *ChatUI) handleHelpEntityRemoveCommand() {
	ui.displaySystemMessage("/entity remove <id|nick>")
}

func (ui *ChatUI) handleEntityRemoveCommand(args []string) {

	if len(args) == 3 {
		err := entity.RemoveNick(args[3])
		if err != nil {
			ui.displaySystemMessage("Error removing entity: " + err.Error())
			return
		}
	} else {
		ui.handleHelpEntityRemoveCommand()
	}

}

func (ui *ChatUI) handleHelpEntityNickCommand(args []string) {
	ui.displaySystemMessage("Usage: /entity nick <id> <nick>")
	ui.displaySystemMessage("Set a nick for a entity")
}

func (ui *ChatUI) handleEntityNickCommand(args []string) {

	// No nick given, hence just show the existing nick
	if len(args) == 3 {
		p, err := entity.Lookup(args[2])
		if err != nil {
			ui.displaySystemMessage("Error fetching alias: " + err.Error())
			return
		}
		log.Debugf("%s: %s", p.DID, p.Nick)
		ui.displaySystemMessage(fmt.Sprintf("Alias for %s is set to %s", p.DID, p.Nick))
		return
	}

	if len(args) == 4 {
		e, err := entity.Lookup(args[2])
		if err != nil {
			ui.displaySystemMessage("Error fetching alias: " + err.Error())
			return
		}
		err = e.SetNick(args[3])
		if err != nil {
			ui.displaySystemMessage("Error setting alias: " + err.Error())
			return
		}
		log.Debugf("Setting alias for %s to %s", e.DID, e.Nick)
		ui.displaySystemMessage(fmt.Sprintf("Alias for %s set to %s", e.DID, e.Nick))
		return
	}

	ui.handleHelpEntityNickCommand(args)

}

func (ui *ChatUI) handleEntityShowCommand(args []string) {

	if len(args) == 3 {
		id := args[2]
		e, err := entity.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		entityInfo := fmt.Sprintf("DID: %s\nNick: %s\n", e.DID, e.Nick)
		ui.displaySystemMessage(entityInfo)
	} else {
		ui.handleHelpEntityShowCommand(args)
	}

}

func (ui *ChatUI) handleHelpEntityShowCommand(args []string) {
	ui.displaySystemMessage("Usage: /entity show <id|nick>")
	ui.displaySystemMessage("Shows the entity info")
}
