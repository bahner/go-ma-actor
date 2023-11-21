package config

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bahner/go-ma/p2p"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

const nodeListenPort = 4001

var (
	err error

	ctxInitDiscoveryTimeout context.Context
	cancel                  context.CancelFunc

	libp2pOpts []libp2p.Option
	n          host.Host
	ps         *pubsub.PubSub
)

// Initialize the p2p node and pubsub service.
// Requires the keyset to be initialized.
func initP2P(timeout int) {

	discoveryTimeout := time.Duration(timeout) * time.Second

	wgDiscovery := &sync.WaitGroup{}
	wgDiscovery.Add(1)
	doDiscovery(wgDiscovery, discoveryTimeout)
	wgDiscovery.Wait()
}

func doDiscovery(wg *sync.WaitGroup, timeout time.Duration) {

	defer wg.Done()

	ctx := context.Background()
	k := GetKeyset()

	ctxInitDiscoveryTimeout, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	// Configure libp2p from here only
	libp2pOpts = []libp2p.Option{
		libp2p.ListenAddrStrings(getListenAddrStrings()...),
		libp2p.Identity(k.IPNSKey.PrivKey),
	}
	n, ps, err = p2p.Init(ctxInitDiscoveryTimeout, libp2pOpts...)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize p2p: %v", err))
	}

}

func getListenAddrStrings() []string {

	port := strconv.Itoa(nodeListenPort)

	return []string{
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip4/0.0.0.0/udp/" + port + "/quic",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",
		"/ip6/::/tcp/" + port,
		"/ip6/::/udp/" + port + "/quic",
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",
	}
}

func GetNode() host.Host {
	return n
}

func GetPubSub() *pubsub.PubSub {
	return ps
}
