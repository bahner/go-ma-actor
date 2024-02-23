package pong

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultReply = "Pong!"

func init() {
	pflag.String("pong-reply", defaultReply, "The message to send back to the sender")
	viper.BindPFlag("pong.reply", pflag.Lookup("pong-reply"))
}

// Run the pong actor. Cancel it from outside to stop it.
func Run(ctx context.Context, a *actor.Actor, n *p2p.P2P) {

	log.Infof("Starting pong mode as %s", a.Entity.DID.Id)
	go a.Subscribe(ctx, a.Entity)

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)

	fmt.Printf("Running in pong mode as %s@%s\n", a.Entity.DID.Id, n.Node.ID())
	fmt.Println("Press Ctrl-C to stop.")

	for {
		<-ctx.Done()
		log.Info("Pong run loop cancelled, exiting...")
		return
	}
}
