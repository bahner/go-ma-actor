package actor

import (
	"context"
	"time"

	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

const broadcastWait = 3 * time.Second

func HelloWorld(ctx context.Context, a *Actor, b *p2ppubsub.Topic) {

	time.Sleep(broadcastWait) // Wait for the network to be ready. This is why we run in a goroutine.

	if a == nil {
		return
	}

	if b == nil {
		return
	}

	if a.Entity == nil {
		return
	}

	me := a.Entity.DID.Id
	greeting := []byte("Hello, world! " + me + " is here.")

	mesg, _ := msg.NewBroadcast(me, greeting, "text/plain", a.Keyset.SigningKey.PrivKey)
	mesg.Broadcast(ctx, b)
}
