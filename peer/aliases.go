package peer

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var aliases = viper.Viper{}

type EntityAlias struct {
	Nick string
	Did  string
}

type NodeAlias struct {
	Nick string
	Id   string
}

func init() {

	// Look in the current directory, the home directory and /etc for the config file.
	// In that order.
	aliases = *viper.New()
	aliases.SetConfigName("aliases")
	aliases.SetConfigType("yaml")
	aliases.AddConfigPath("$HOME/.ma")
	aliases.AddConfigPath("/etc/ma")

	// Allow to set config file via command line flag.
	aliases.SetDefault("entities", []EntityAlias{})
	aliases.SetDefault("nodes", []NodeAlias{})

	aliases.WatchConfig()

	err := aliases.ReadInConfig()
	if err != nil {
		log.Errorf("Error reading aliases config file: %s", err)
	}
}

func GetEntityAliases() []EntityAlias {

	var entityAliases = []EntityAlias{}

	err := aliases.UnmarshalKey("entities", &entityAliases)
	if err != nil {
		log.Errorf("Error unmarshalling entity aliases: %s", err)
	}

	return entityAliases
}

func GetNodeAliases() []NodeAlias {
	var nodeAliases = []NodeAlias{}

	err := aliases.UnmarshalKey("nodes", &nodeAliases)
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

func GetNodeAlias(id string) string {

	aliases := GetNodeAliases()

	for _, alias := range aliases {
		if alias.Id == id {
			return alias.Nick
		}
	}

	return ""
}

func GetEntityDid(nick string) string {

	aliases := GetEntityAliases()

	for _, alias := range aliases {
		if alias.Nick == nick {
			return alias.Did
		}
	}

	return ""
}

func GetNodeId(nick string) string {

	aliases := GetNodeAliases()

	for _, alias := range aliases {
		if alias.Nick == nick {
			return alias.Id
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
			aliases.Set("nodes", nodeAliases)

			// Write the changes to the config file
			return aliases.WriteConfig()
		}
	}

	// If the nick does not exist, add it
	nodeAliases = append(nodeAliases, NodeAlias{Nick: nick, Id: id})
	aliases.Set("nodes", nodeAliases)
	return aliases.WriteConfig()

}

func RemoveNodeAlias(nick string) error {

	nodeAliases := GetNodeAliases()

	for i, alias := range nodeAliases {

		// Check if nick already exists
		if alias.Nick == nick {

			// Remove the nick
			log.Debugf("Removing alias %s for entity %s", nick, alias.Id)
			nodeAliases = append(nodeAliases[:i], nodeAliases[i+1:]...)
			aliases.Set("nodes", nodeAliases)

			// Write the changes to the config file
			return aliases.WriteConfig()
		}
	}

	return nil
}

func AddEntityAlias(did string, nick string) error {

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
			aliases.Set("entities", entityAliases)

			// Write the changes to the config file
			return aliases.WriteConfig()
		}
	}

	// If the nick does not exist, add it
	entityAliases = append(entityAliases, EntityAlias{Nick: nick, Did: did})
	aliases.Set("entities", entityAliases)
	return aliases.WriteConfig()
}

func RemoveEntityAlias(nick string) error {

	entityAliases := GetEntityAliases()

	for i, alias := range entityAliases {

		// Check if nick already exists
		if alias.Nick == nick {

			// Remove the nick
			log.Debugf("Removing alias %s for entity %s", nick, alias.Did)
			entityAliases = append(entityAliases[:i], entityAliases[i+1:]...)
			aliases.Set("entities", entityAliases)

			// Write the changes to the config file
			return aliases.WriteConfig()
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

func RefreshNodeAliases() error {

	na := GetNodeAliases()
	if na == nil {
		return fmt.Errorf("no node aliases found")
	}

	for _, a := range na {
		if a.Id != "" {
			p := get(a.Id)
			if p != nil {
				p.Alias = a.Nick
			}
		}
	}

	return nil
}
