package entity

import (
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/db"
	log "github.com/sirupsen/logrus"
)

var nicks *sync.Map

func init() {
	nicks = new(sync.Map)
}

// Returns the Entity's nick. Uses and sets the DID as fallback.
// This does not actually save the nick, as it it will always return
// the same value for the same DID.
func (e Entity) Nick() string {
	if nick, ok := nicks.Load(e.DID.Id); ok {
		return nick.(string)
	}
	return e.DID.Id
}

// Sets a node in the database
// takes new did and nick. If an old  did  for the alias exists it is removed.
// This makes this the only alias for the DID and the only complex function in this file.
func (e Entity) SetNick(nick string) error {
	nicks.Store(e.DID.Id, nick)
	return db.Save(nicks, config.DBPeers())
}

func DeleteNick(id string) {
	nicks.Delete(id)
}

// Lookup takes an input string that can be either a nick or an ID,
// and returns the corresponding ID if found; otherwise, it returns the input.
func Lookup(q string) string {
	found := ""
	nicks.Range(func(key, value interface{}) bool {
		id := key.(string)
		nick := value.(string)

		// Check if input matches either the key (ID) or the value (nick)
		if q == id || q == nick {
			found = id
			return false // Stop iterating
		}
		return true // Continue iterating
	})

	if found != "" {
		return found
	}
	return q // Return the input if no match is found
}

func Nicks() map[string]string {
	nmap := make(map[string]string) // Pun intended

	nicks.Range(func(key, value interface{}) bool {
		strKey, okKey := key.(string)
		strValue, okValue := value.(string)
		if okKey && okValue {
			nmap[strKey] = strValue
		}
		return true
	})

	return nmap
}

func WatchCSV() error {
	filename := config.DBEntities()
	log.Infof("entity.WatchCSV: watching %s", filename)
	return db.Watch(filename, nicks)
}
