package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

const (
	entityUsage = "/entity nick|show"
	entityHelp  = `Manages info on seen entities
At this point only nicks are handled`
	entityConnectUsage    = "/entity connect <id|nick>"
	entityConnectHelp     = "Connects to an entity's libp2p node"
	entityNickUsage       = "/entity nick list|set|remove|show"
	entityNickHelp        = "Manages nicks for entities"
	entityNickListUsage   = "/entity nick list"
	entityNickListHelp    = "Lists nicks for entities"
	entityNickSetUsage    = "/entity nick set <id|nick> <nick>"
	entityNickSetHelp     = "Sets a nick for an entity"
	entityNickRemoveUsage = "/entity nick remove <id|nick>"
	entityNickRemoveHelp  = "Removes a nick for an entity"
	entityNickShowUsage   = "/entity nick show <id|nick>"
	entityNickShowHelp    = "Shows the entity info"
)

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
		case "connect":
			ui.handleEntityConnectCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpCommand(entityUsage, entityHelp)

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
		ui.handleHelpCommand(entityNickListUsage, entityNickListHelp)
	}

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

	ui.handleHelpCommand(entityNickUsage, entityNickHelp)
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
		ui.handleHelpCommand(entityNickSetUsage, entityNickSetHelp)
	}

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
		ui.handleHelpCommand(entityNickShowUsage, entityNickShowHelp)
	}

}

// REMOVE
func (ui *ChatUI) handleEntityNickRemoveCommand(args []string) {

	if len(args) == 4 {
		id := entity.GetDID(args[3])
		entity.RemoveNick(id)
		ui.displaySystemMessage("Nick removed for " + id + " if it existed")
	} else {
		ui.handleHelpCommand(entityNickRemoveUsage, entityNickRemoveHelp)
	}

}

func (ui *ChatUI) handleEntityConnectCommand(args []string) {

	if len(args) == 3 {
		e, err := entity.GetOrCreate(args[2])
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		pai, err := e.ConnectPeer()
		if err != nil {
			ui.displaySystemMessage("Error connecting to enityty peer: " + err.Error())
			return
		}
		ui.displaySystemMessage("Connected to " + e.DID.Id + ": " + pai.ID.String())
	} else {
		ui.handleHelpCommand(entityConnectUsage, entityConnectHelp)
	}
}
