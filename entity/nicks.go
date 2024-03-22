package entity

import (
	"errors"
	"strings"

	"github.com/bahner/go-ma-actor/db"
)

const (
	entityNickPrefix = entityPrefix + "nick:"
)

var (
	ErrFailedToCreateNick = errors.New("failed to set entity nick")
	ErrDIDNotFound        = errors.New("DID not found")
	ErrNickNotFound       = errors.New("Nick not found")
)

// This returns the DID for the nick.
// If it doesn't exist it returns the query.
func DID(query string) string {

	queryBytes := []byte(entityNickPrefix + query)

	id, err := db.Get(queryBytes)
	if err != nil {
		return query
	}

	return string(id)
}

// Removes a an entity nick from the database if it exists.
// It does a lookup to you can enter both a nick or a DID.
func DeleteNick(q string) error {

	nick := Nick(q)
	nickBytes := []byte(entityNickPrefix + nick)

	return db.Delete(nickBytes)
}

// Nick is the opposite of LookupDID. It returns the nick for a DID.
// This means iterating over all nicks to find the right one,
// hence it's not as efficient as LookupDID.
func Nick(q string) string {

	qBytes := []byte(q)

	id, err := db.Lookup(qBytes)
	if err != nil {
		return q
	}

	return strings.TrimPrefix(string(id), entityNickPrefix)
}

func Nicks() (map[string]string, error) {
	return db.Keys(entityNickPrefix)
}
