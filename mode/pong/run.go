package pong

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("pong-reply", config.DefaultPongReply, "The message to send back to the sender")
}

// Run the pong actor. Cancel it from outside to stop it.
func Run(a *actor.Actor, p *p2p.P2P) {

	ctx := context.Background()

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", config.DefaultPongReply)

	fmt.Printf("Starting pong mode as %s\n", a.Entity.DID.Id)
	go p.DiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")
	go a.Subscribe(ctx, a.Entity)
	fmt.Println("Subscribed to self.")

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)
	fmt.Println("Started event handlers.")

	actor.HelloWorld(ctx, a)
	fmt.Println("Sent hello world.")

	fmt.Printf("Running in pong mode as %s@%s\n", a.Entity.DID.Id, p.DHT.Host().ID())
	fmt.Println("Press Ctrl-C to stop.")

	for {
		// A background loop that does nothing.
		// The ctx will never be cancelled, so this will run forever.
		<-ctx.Done()
		fmt.Println("Pong run loop cancelled, exiting...")
		return
	}
}
