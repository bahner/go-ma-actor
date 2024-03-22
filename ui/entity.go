package ui

import (
	"fmt"
	"strings"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
)

const (
	entityUsage        = "/entity delete|list|nick|resolve|show"
	entityHelp         = "Manages info on seen entities"
	entityConnectUsage = "/entity connect <id|nick>"
	entityConnectHelp  = "Connects to an entity's libp2p node"
	entityListUsage    = "/entity list"
	entityListHelp     = "Lists nicks for entities"
	entityNickUsage    = "/entity nick <id|nick> <nick>"
	entityNickHelp     = `Sets or shows a nick for an entity
The entity to set nick for *MUST* be quoted if it contains spaces.
The nick after the entity to set nick for doesn't need to be quoted.
`
	entityDeleteUsage  = "/entity delete <id|nick>"
	entityDeleteHelp   = "Deletes an entity from the database"
	entityResolveUsage = "/entity resolve <DID|NICK>"
	entityResolveHelp  = "Tries to resolve the most recent version of the DID Document for the given DID or NICK."
	entityShowUsage    = "/entity show <id|nick>"
	entityShowHelp     = "Show info about the entity"
	aliasSeparator     = "\t => "
)

func (ui *ChatUI) handleEntityCommand(args []string) {

	if len(args) >= 2 {
		command := args[1]
		switch command {
		case "connect":
			ui.handleEntityConnectCommand(args)
			return
		case "delete":
			ui.handleEntityDeleteCommand(args)
			return
		case "list":
			ui.handleEntityListCommand(args)
			return
		case "nick":
			ui.handleEntityNickCommand(args)
			return
		case "resolve":
			go ui.handleEntityResolveCommand(args) // This make take some time. No need to block the UI
			return
		case "show":
			ui.handleEntityShowCommand(args)
			return
		default:
			ui.displaySystemMessage("Unknown alias entity command: " + command)
		}
	}

	ui.handleHelpCommand(entityUsage, entityHelp)

}

func (ui *ChatUI) handleEntityListCommand(args []string) {

	log.Debugf("entity list command: %v", args)
	if len(args) == 2 {

		nicks, err := entity.Nicks()
		if err != nil {
			ui.displaySystemMessage("Error fetching nicks: " + err.Error())
			return
		}
		log.Debugf("nicks: %v", nicks)

		if len(nicks) > 0 {
			for k, v := range nicks {
				ui.displaySystemMessage(k + aliasSeparator + v)
			}
		} else {
			ui.displaySystemMessage("No entities found")
		}
		return
	}
	ui.handleHelpCommand(entityListUsage, entityListHelp)

}

func (ui *ChatUI) handleEntityNickCommand(args []string) {

	if len(args) >= 4 {
		id := entity.DID(args[2])
		nick := strings.Join(args[3:], separator)

		e, err := entity.GetOrCreate(id, true)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		err = e.SetNick(nick)
		if err != nil {
			ui.displaySystemMessage("Error setting entity nick: " + err.Error())
			return
		}
		// Change the window title if the ID matches the current entity
		if id == ui.e.DID.Id {
			ui.msgBox.SetTitle(nick)
		}
		ui.displaySystemMessage(e.DID.Id + aliasSeparator + e.Nick())
		return
	}
	ui.handleHelpCommand(entityNickUsage, entityNickHelp)
}

// SHOW
func (ui *ChatUI) handleEntityShowCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		e, err := entity.GetOrCreate(entity.DID(id), true)
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		entityInfo := fmt.Sprintf(e.DID.Id + aliasSeparator + e.Nick())
		ui.displaySystemMessage(entityInfo)
		return
	}
	ui.handleHelpCommand(entityShowUsage, entityShowHelp)

}

// REMOVE
func (ui *ChatUI) handleEntityDeleteCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		err := entity.DeleteNick(id)
		if err != nil {
			ui.displaySystemMessage("Error deleting nick: " + err.Error())
			return
		}
		ui.displaySystemMessage("Nick deleted for " + id + " if it existed")
		return
	}
	ui.handleHelpCommand(entityDeleteUsage, entityDeleteHelp)

}

func (ui *ChatUI) handleEntityConnectCommand(args []string) {

	if len(args) >= 3 {
		id := strings.Join(args[2:], separator)
		id = entity.DID(id)
		e, err := entity.GetOrCreate(id, false) // Lookup up the entity document properly.
		if err != nil {
			ui.displaySystemMessage("Error: " + err.Error())
			return
		}
		log.Debugf("Connecting to peer for entity: %v", id)
		pai, err := e.ConnectPeer()
		if err != nil {
			msg := fmt.Sprintf("Error connecting to entity peer: %v. %v", pai.ID, err)
			ui.displaySystemMessage(msg)
			return
		}
		ui.displaySystemMessage("Connected to " + id + aliasSeparator + pai.ID.String())
		return
	}
	ui.handleHelpCommand(entityConnectUsage, entityConnectHelp)
}

func (ui *ChatUI) handleEntityResolveCommand(args []string) {

	if len(args) >= 3 {

		// We must absolutely get the entity, so we can get the DID Document.
		// So no caching.
		CACHED := false

		id := strings.Join(args[2:], separator)
		id = entity.DID(id)

		e, err := entity.GetOrCreate(id, CACHED)
		if err != nil {
			ui.displaySystemMessage("Error fetching entity: " + err.Error())
			return
		}

		ui.displaySystemMessage("Resolving DID Document for " + e.DID.Id + "...")
		d, c, err := doc.Fetch(id, false)
		if err != nil {
			ui.displaySystemMessage("Error fetching DID Document: " + err.Error())
			return
		}
		ui.displaySystemMessage("Resolved DID Document for " + e.DID.Id + " (CID: " + c.String() + ")")
		e.Doc = d

		return

	}
	ui.handleHelpCommand(entityResolveUsage, entityResolveHelp)

}
