package peer

import (
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/db"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	log "github.com/sirupsen/logrus"
)

const defaultNickLength = 8

var nicks *sync.Map

func init() {
	nicks = new(sync.Map)
}

// If a nick is not found, it creates and sets a new one.
func AssertNick(pid p2peer.ID) error {
	if IsKnown(pid) {
		return nil
	}

	id := pid.String()
	nick := createNick(id)
	return SetNick(id, nick)
}

// LookupNick looks for a peer's nick by its ID.
// If it does not exist it returns the input
func LookupNick(peerID p2peer.ID) string {

	id := peerID.String()

	if nick, ok := nicks.Load(id); ok {
		return nick.(string)
	}

	return id
}

func DeleteNick(id string) {
	nicks.Delete(id)
}

func IsKnown(peerID p2peer.ID) bool {

	_, ok := nicks.Load(peerID.String())

	return ok
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

// Nick returns the nick for a given peer ID.
func Nick(q string) (string, error) {

	nick, ok := nicks.Load(q)
	if !ok {
		return "", ErrNickNotFound
	}

	return nick.(string), nil

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

// SetNick updates or sets a new nick for a given peer ID.
func SetNick(id, nick string) error {

	// Nicks must be unique, so delete any old ones.
	DeleteNick(Lookup(nick))

	nicks.Store(id, nick)
	return db.Save(nicks, config.DBPeers())
}

func WatchCSV() error {
	filename := config.DBPeers()
	log.Infof("peer.WatchCSV: watching %s", filename)
	return db.Watch(filename, nicks)
}

// createNick creates a short nick from a peer ID.
// It simply takes the last 8 characters of the ID.
// If the ID is shorter than 8 characters, it returns the ID as is.
func createNick(id string) string {
	if len(id) > defaultNickLength {
		return id[len(id)-defaultNickLength:]
	}
	return id
}
