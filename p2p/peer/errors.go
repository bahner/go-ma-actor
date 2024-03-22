package peer

import "errors"

var (
	ErrAddrInfoAddrsEmpty      = errors.New("addrInfo.Addrs is empty")
	ErrPeerDenied              = errors.New("peer not allowed")
	ErrPeerNotAllowedByDefault = errors.New("peer not allowed by default")
	ErrAlreadyConnected        = errors.New("already connected")
	ErrFailedToCreateNick      = errors.New("failed to set entity nick")
	ErrDIDNotFound             = errors.New("DID not found")
	ErrNickNotFound            = errors.New("Nick not found")
)
