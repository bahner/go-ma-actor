package peer

import "errors"

var (
	ErrPeerNotFound    = errors.New("peer not found")
	ErrInvalidPeerType = errors.New("invalid peer type")
)
