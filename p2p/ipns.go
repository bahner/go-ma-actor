package p2p

import (
	"github.com/ipfs/boxo/namesys"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
)

func newIPNSPublisher(d *p2pDHT.IpfsDHT) *namesys.IPNSPublisher {

	// Set up a DHT for the host
	// Create a new synchronized in-memory datastore
	ds := dssync.MutexWrap(datastore.NewMapDatastore())

	return namesys.NewIPNSPublisher(d, ds)
}

func newIPNSResolver(d *p2pDHT.IpfsDHT) *namesys.IPNSResolver {

	return namesys.NewIPNSResolver(d)
}
