package pong

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("pong-reply", config.DefaultPongReply, "The message to send back to the sender")
}

// Run the pong actor. Cancel it from outside to stop it.
func Run(ctx context.Context, a *actor.Actor, b *p2ppubsub.Topic, n *p2p.P2P) {

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", config.DefaultPongReply)

	log.Infof("Starting pong mode as %s", a.Entity.DID.Id)
	go n.DiscoveryLoop(ctx)
	go a.Subscribe(ctx, a.Entity)

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)

	actor.HelloWorld(ctx, a, b)

	fmt.Printf("Running in pong mode as %s@%s\n", a.Entity.DID.Id, n.DHT.Host().ID())
	fmt.Println("Press Ctrl-C to stop.")

	for {
		<-ctx.Done()
		log.Info("Pong run loop cancelled, exiting...")
		return
	}
}
