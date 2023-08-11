package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/pubsub"
	"github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p/core/peer"
)

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
	go initPubSubService(ctx, wg, h)
	wg.Wait()

	// create and join the chat room, ps is now initialized.
	cr, err := newChatRoom(ctx, ps, nick, room)
	if err != nil {
		panic(err)
	}

	// draw the UI
	ui := NewChatUI(ctx, cr)
	if err := ui.Run(); err != nil {
		printErr("error running text UI: %s", err)
	}
}

func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}
