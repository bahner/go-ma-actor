package peer

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/db"
)

const peerAllowedPrefix = peerPrefix + "allowed:"

func IsAllowed(id string) bool {

	allowedIdBytes := []byte(peerAllowedPrefix + id)

	allowed, err := db.Get(allowedIdBytes)
	if err != nil {
		return config.ALLOW_ALL_PEERS
	}
	return byteToBool(allowed)
}

func SetAllowed(id string, allowed bool) error {

	allowedIdBytes := []byte(peerAllowedPrefix + id)

	return db.Set(allowedIdBytes, boolToByte(allowed))
}

func boolToByte(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

func byteToBool(b []byte) bool {
	return b[0] == 1
}
