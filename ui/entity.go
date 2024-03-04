package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleHelpEntityCommands() {
	ui.displayHelpUsage("/entity nick list|remove|show|nick")
	ui.displayHelpText("Manages entity info")
	ui.displayHelpText("At this point only nick are handled")
}

func (ui *ChatUI) handleEntityCommand(args []string) {

	if len(args) >= 2 {
		command := args[1]
		switch command {
		case "nick":
			ui.handleEntityNickCommand(args)
			return
		case "show":
			ui.handleEntityNickShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpEntityCommands()

}

func (ui *ChatUI) handleEntityNickListCommand(args []string) {

	log.Debugf("entity list command: %v", args)
	if len(args) == 3 {

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
		ui.handleHelpEntityListCommand()
	}

}

func (ui *ChatUI) handleHelpEntityListCommand() {
	ui.displayHelpUsage("/entity list")
	ui.displayHelpText("List entity DID and nicks")
}

func (ui *ChatUI) handleHelpEntityNickCommand(args []string) {
	ui.displayHelpUsage("/entity nick set|remove|show")
	ui.displayHelpText("Set a nick for a entity")
	ui.handleEntityNickSetCommand(args)
	ui.handleEntityNickRemoveCommand(args)
	ui.handleEntityNickShowCommand(args)
}

// case "remove":
// 	ui.handleEntityRemoveCommand(args)
// 	return

func (ui *ChatUI) handleEntityNickCommand(args []string) {

	if len(args) >= 3 {
		command := args[2]
		switch command {
		case "list":
			ui.handleEntityNickListCommand(args)
			return
		case "set":
			ui.handleEntityNickSetCommand(args)
			return
		case "remove":
			ui.handleEntityNickRemoveCommand(args)
			return
		case "show":
			ui.handleEntityNickShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpEntityNickCommand(args)
}

// SET
func (ui ChatUI) handleEntityNickSetCommand(args []string) {

	if len(args) == 5 {
		id := args[3]
		nick := args[4]
		e, err := entity.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		err = e.SetNick(nick)
		if err != nil {
			ui.displaySystemMessage("Error setting entity nick: " + err.Error())
			return
		}
		ui.displaySystemMessage(e.DID.Id + " is now known as " + e.Nick)
	} else {
		ui.handleHelpEntityNickSetCommand()
	}

}

func (ui *ChatUI) handleHelpEntityNickSetCommand() {
	ui.displayHelpUsage("/entity nick set <id|nick> <nick>")
	ui.displayHelpText("Sets a nick for an entity")
}

// SHOW
func (ui *ChatUI) handleEntityNickShowCommand(args []string) {

	if len(args) == 4 {
		id := args[3]
		e, err := entity.Lookup(id)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		entityInfo := fmt.Sprintf(e.DID.Id + " is also known as " + e.Nick)
		ui.displaySystemMessage(entityInfo)
	} else {
		ui.handleHelpEntityShowCommand()
	}

}

func (ui *ChatUI) handleHelpEntityShowCommand() {
	ui.displayHelpUsage("/entity nick show <id|nick>")
	ui.displayHelpText("Shows the entity info")
}

// REMOVE
func (ui *ChatUI) handleEntityNickRemoveCommand(args []string) {

	if len(args) == 4 {
		id := entity.GetDID(args[3])
		entity.RemoveNick(id)
		ui.displaySystemMessage("Nick removed for " + id + " if it existed")
	} else {
		ui.handleHelpEntityNickRemoveCommand()
	}

}

func (ui *ChatUI) handleHelpEntityNickRemoveCommand() {
	ui.displayHelpUsage("/entity nick remove <id|nick>")
	ui.displayHelpText("Removes a nick for an entity")
}
