package peer

import (
	"strings"

	"github.com/bahner/go-ma-actor/db"
)

// Sets a node in the database
// takes new did and nick. If an old  did  for the alias exists it is removed.
// This makes this the only alias for the DID and the only complex function in this file.
func SetNick(nick string, id string) error {

	prefixBytes := []byte(peerNickPrefix)
	keyBytes := []byte(nick)
	valBytes := []byte(id)

	return db.Upsert(prefixBytes, keyBytes, valBytes)

}

// Returns the Entity's nick. If it doesn't exist it returns the query.
func Nick(id string) string {

	idBytes := []byte(id)

	key, err := db.Lookup(idBytes)
	if err != nil {
		return id
	}

	return strings.TrimPrefix(string(key), peerNickPrefix)
}

func Nicks() (map[string]string, error) {
	return db.Keys(peerNickPrefix)
}

// Removes a nnick from the database if it exists.
func DeleteNick(nick string) error {

	nickBytes := []byte(peerNickPrefix + nick)

	return db.Delete(nickBytes)
}

// ID is the opposite of Nick. It returns the PeerID for a nick.
func ID(q string) string {

	prefixBytes := []byte(peerNickPrefix)
	qBytes := []byte(q)

	id, err := db.Get(append(prefixBytes, qBytes...))
	if err != nil {
		return q
	}

	return string(id)
}

// Always returns a nick for the did, but also an error if it fails to set the nick.
func getOrCreateNick(id string) (nick string) {
	nick = Nick(id)
	if nick == id {
		nick = createDefaultNick(id)
	}
	return nick
}

// Create a default nick from the DID
func createDefaultNick(s string) string {
	runes := []rune(s)

	// This should be made from did's so length should be more than 8
	if len(runes) > 8 {
		return string(runes[len(runes)-defailtNickLength:])
	}
	// This is strange, but just return the string
	return s
}
