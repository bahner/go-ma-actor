package peer

import "errors"

var (
	ErrPeerNotFound            = errors.New("peer not found")
	ErrInvalidPeerType         = errors.New("invalid peer type")
	ErrAddrInfoAddrsEmpty      = errors.New("addrInfo.Addrs is empty")
	ErrPeerNotAllowed          = errors.New("peer not allowed")
	ErrPeerNotAllowedByDefault = errors.New("peer not allowed by default")
)
