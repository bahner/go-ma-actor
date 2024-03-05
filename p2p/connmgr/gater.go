package connmgr

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/control"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ConnectionGater is a struct that implements the network.ConnectionGater interface.
// It uses a sync.Map to store valid peer IDs that have been discovered using the correct rendezvous string.
type ConnectionGater struct {
	validPeers sync.Map // Stores peer.ID as key
}

// New creates a new CustomConnectionGater instance.
func NewConnectionGater() *ConnectionGater {
	return &ConnectionGater{}
}

// InterceptPeerDial checks if we should allow dialing the specified peer.
func (cg *ConnectionGater) InterceptPeerDial(p peer.ID) (allow bool) {
	_, found := cg.validPeers.Load(p)
	return found
}

// InterceptAccept checks if an incoming connection from the specified network address should be allowed.
func (cg *ConnectionGater) InterceptAccept(conn network.ConnMultiaddrs) (allow bool) {
	return true
}

// InterceptSecured, InterceptUpgraded, and other methods can be implemented as needed.
// For simplicity, they are set to allow all connections in this example.
func (cg *ConnectionGater) InterceptSecured(_ network.Direction, p peer.ID, _ network.ConnMultiaddrs) (allow bool) {
	_, found := cg.validPeers.Load(p)

	return found

}

func (cg *ConnectionGater) InterceptAddrDial(_ peer.ID, _ multiaddr.Multiaddr) (allow bool) {
	return true
}

func (cg *ConnectionGater) InterceptUpgraded(_ network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}

func (cg *ConnectionGater) AddPeer(p peer.ID) {
	cg.validPeers.Store(p, struct{}{})
}

func (cg *ConnectionGater) RemovePeer(p peer.ID) {
	cg.validPeers.Delete(p)
}

func (cg *ConnectionGater) ListPeers() []peer.ID {
	var peers []peer.ID
	cg.validPeers.Range(func(k, v interface{}) bool {
		peers = append(peers, k.(peer.ID))
		return true
	})
	return peers
}

func (cg *ConnectionGater) Clear() {
	cg.validPeers = sync.Map{}
}

func (cg *ConnectionGater) Count() int {
	var count int
	cg.validPeers.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

func (cg *ConnectionGater) HasPeer(p peer.ID) bool {
	_, found := cg.validPeers.Load(p)
	return found
}

func (cg *ConnectionGater) Close() {
	cg.Clear()
}

func (cg *ConnectionGater) IsEmpty() bool {
	return cg.Count() == 0
}
