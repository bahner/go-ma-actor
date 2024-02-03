package alias

import (
	"fmt"

	"github.com/bahner/go-ma/did"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EntityAlias struct {
	Nick string
	Did  string
}

type NodeAlias struct {
	Nick string
	Id   string
}

func init() {

	// Allow to set config file via command line flag.
	viper.SetDefault("entities", []EntityAlias{})
	viper.SetDefault("nodes", []NodeAlias{})
}

func GetEntityAliases() []EntityAlias {

	var entityAliases = []EntityAlias{}

	err := viper.UnmarshalKey("entities", &entityAliases)
	if err != nil {
		log.Errorf("Error unmarshalling entity aliases: %s", err)
	}

	return entityAliases
}

func GetNodeAliases() []NodeAlias {
	var nodeAliases = []NodeAlias{}

	err := viper.UnmarshalKey("nodes", &nodeAliases)
	if err != nil {
		log.Errorf("Error unmarshalling node aliases: %s", err)
	}

	return nodeAliases
}

func GetEntityAlias(did string) string {

	aliases := GetEntityAliases()

	for _, alias := range aliases {
		if alias.Did == did {
			return alias.Nick
		}
	}

	return ""
}
func GetEntityDID(nick string) string {

	aliases := GetEntityAliases()

	for _, alias := range aliases {
		if alias.Nick == nick {
			return alias.Did
		}
	}

	return ""
}
func GetNodeDID(nick string) string {

	aliases := GetNodeAliases()

	for _, alias := range aliases {
		if alias.Nick == nick {
			return alias.Id
		}
	}

	return ""
}

func GetNodeAlias(id string) string {

	aliases := GetNodeAliases()

	for _, alias := range aliases {
		if alias.Id == id {
			return alias.Nick
		}
	}

	return ""
}

func AddNodeAlias(id string, nick string) error {

	nodeAliases := GetNodeAliases()

	for i, alias := range nodeAliases {

		// Check if nick already exists
		if alias.Id == id {

			// If the nick is already set, do nothing
			if alias.Nick == nick {
				log.Debugf("Alias %s already set for entity %s", nick, id)
				return nil
			}

			// Otherwise, change the nick
			log.Debugf("Changing alias %s to %s for entity %s", alias.Nick, nick, id)
			nodeAliases[i].Nick = nick
			viper.Set("nodes", nodeAliases)

			// Write the changes to the config file
			return viper.WriteConfig()
		}
	}

	// If the nick does not exist, add it
	nodeAliases = append(nodeAliases, NodeAlias{Nick: nick, Id: id})
	viper.Set("nodes", nodeAliases)
	return viper.WriteConfig()

}

func RemoveNodeAlias(nick string) error {

	nodeAliases := GetNodeAliases()

	for i, alias := range nodeAliases {

		// Check if nick already exists
		if alias.Nick == nick {

			// Remove the nick
			log.Debugf("Removing alias %s for entity %s", nick, alias.Id)
			nodeAliases = append(nodeAliases[:i], nodeAliases[i+1:]...)
			viper.Set("nodes", nodeAliases)

			// Write the changes to the config file
			return viper.WriteConfig()
		}
	}

	return nil
}

func AddEntityAlias(did string, nick string) error {

	// Lookup up possible existing alias
	GetEntityDID(nick)

	entityAliases := GetEntityAliases()

	for i, alias := range entityAliases {

		// Check if nick already exists
		if alias.Did == did {

			// If the nick is already set, do nothing
			if alias.Nick == nick {
				log.Debugf("Alias %s already set for entity %s", nick, did)
				return nil
			}

			// Otherwise, change the nick
			log.Debugf("Changing alias %s to %s for entity %s", alias.Nick, nick, did)
			entityAliases[i].Nick = nick
			viper.Set("entities", entityAliases)

			// Write the changes to the config file
			return viper.WriteConfig()
		}
	}

	// If the nick does not exist, add it
	entityAliases = append(entityAliases, EntityAlias{Nick: nick, Did: did})
	viper.Set("entities", entityAliases)
	return viper.WriteConfig()
}

func RemoveEntityAlias(nick string) error {

	entityAliases := GetEntityAliases()

	for i, alias := range entityAliases {

		// Check if nick already exists
		if alias.Nick == nick {

			// Remove the nick
			log.Debugf("Removing alias %s for entity %s", nick, alias.Did)
			entityAliases = append(entityAliases[:i], entityAliases[i+1:]...)
			viper.Set("entities", entityAliases)

			// Write the changes to the config file
			return viper.WriteConfig()
		}
	}

	return nil
}

func PrintEntityAliases() string {

	aliases := GetEntityAliases()

	aliases_string := "Entities:\n"

	for _, alias := range aliases {
		aliases_string += fmt.Sprintf("%s: %s\n", alias.Nick, alias.Did)
	}

	return aliases_string
}

func PrintNodeAliases() string {

	aliases := GetNodeAliases()

	aliases_string := "Nodes:\n"

	for _, alias := range aliases {
		aliases_string += fmt.Sprintf("%s: %s\n", alias.Nick, alias.Id)
	}

	return aliases_string
}

func Nick(id string) string {

	var nick string

	// First check if this is a valid did
	// Lookup it up or return the fragment
	_did, err := did.New(id)

	// _did is a *did.DID and id is a valid DID
	if err == nil {
		nick = GetEntityAlias(id)
		if nick != "" {
			return nick
		}
		nick = GetNodeAlias(id)
		if nick != "" {
			return nick
		}
		return _did.Fragment
	}

	// This means that what we got is not a valid DID
	// but just a nick itself, so it was already the best
	return nick
}
