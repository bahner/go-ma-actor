package peer

import "errors"

var (
	ErrAddrInfoAddrsEmpty = errors.New("addrInfo.Addrs is empty")
	ErrAlreadyConnected   = errors.New("already connected")
	ErrNickNotFound       = errors.New("Nick not found")
)
