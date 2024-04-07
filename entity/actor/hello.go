package actor

import (
	"context"
	"time"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	"github.com/bahner/go-ma/msg"
)

const broadcastWait = 3 * time.Second

func (a *Actor) HelloWorld(ctx context.Context) {

	topic, err := pubsub.GetOrCreateTopic(ma.BROADCAST_TOPIC)
	if err != nil {
		return
	}

	time.Sleep(broadcastWait) // Wait for the network to be ready. This is why we run in a goroutine.

	if a == nil {
		return
	}

	if a.Entity == nil {
		return
	}

	me := a.Entity.DID.Id
	greeting := []byte("Hello, world! " + me + " is here.")

	mesg, _ := msg.NewBroadcast(me, greeting, "text/plain", a.Keyset.SigningKey.PrivKey)
	mesg.Broadcast(ctx, topic)
}
