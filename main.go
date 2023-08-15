package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/bahner/go-myspace/p2p/host"
	"github.com/bahner/go-myspace/p2p/key"
	"github.com/bahner/go-myspace/p2p/pubsub"
	"github.com/sirupsen/logrus"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	ps  *pubsub.Service
	log *logrus.Logger
)

func main() {
	ctx := context.Background()

	log = logrus.New()
	f, err := os.OpenFile("go-myspace-client.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer f.Close()

	// Set the log output to the file
	log.SetOutput(f)

	initConfig()

	identity := key.CreateIdentity(secret)

	h := host.New()
	h.AddOption(libp2p.Identity(identity))
	h.AddOption(libp2p.ListenAddrStrings())
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
