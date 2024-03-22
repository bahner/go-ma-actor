package peer

import (
	"github.com/bahner/go-ma-actor/db"
)

const (
	peerPrefix        = "peer:"
	peerNickPrefix    = peerPrefix + "nick:"
	defailtNickLength = 8
)

// Use the the fact that IsAllowed is set as an indicator of known peers
func IsKnown(id string) bool {

	allowedIdBytes := []byte(peerAllowedPrefix + id)

	_, err := db.Get(allowedIdBytes)
	return err == nil
}

func GetOrCreate(id string) (string, error) {

	nick := getOrCreateNick(id)
	if nick == id {
		return nick, ErrFailedToCreateNick
	}

	return nick, nil
}
