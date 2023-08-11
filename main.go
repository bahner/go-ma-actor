package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
	"github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p/core/peer"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "myspace"

var (
	ps  *pubsub.Service
	log *logrus.Logger
)

func main() {

	ctx := context.Background()

	log = logrus.New()

	initConfig()

	h := host.New()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	initPubSubService(ctx, wg, h)
	wg.Wait()

	// join the chat room
	cr, err := JoinChatRoom(ctx, ps.Sub, h.Node.ID(), nick, room)
	if err != nil {
		panic(err)
	}

	// draw the UI
	ui := NewChatUI(cr)
	if err := ui.Run(); err != nil {
		printErr("error running text UI: %s", err)
	}
}

// printErr is like fmt.Printf, but writes to stderr.
func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

// shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}
