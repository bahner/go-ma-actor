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
	pflag.String("pong-reply", config.DEFAULT_PONG_REPLY, "The message to send back to the sender")

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", config.DEFAULT_PONG_REPLY)

	viper.BindPFlag("mode.pong.fortune.enable", pflag.Lookup("pong-fortune"))
	viper.SetDefault("mode.pong.fortune.enable", config.DEFAULT_PONG_FORTUNE_MODE)

	viper.SetDefault("mode.pong.fortune.args", config.DEFAULT_PONG_FORTUNE_ARGS)
}

// Run the pong actor. Cancel it from outside to stop it.
func Run(a *actor.Actor, p *p2p.P2P) {

	ctx := context.Background()

	fmt.Printf("Starting pong mode as %s\n", a.Entity.DID.Id)
	go p.StartDiscoveryLoop(ctx)
	fmt.Println("Discovery loop started.")
	go a.Subscribe(ctx, a.Entity)
	fmt.Println("Subscribed to self.")

	go handleEnvelopeEvents(ctx, a)
	go handleMessageEvents(ctx, a)
	fmt.Println("Started event handlers.")

	actor.HelloWorld(ctx, a)
	fmt.Println("Sent hello world.")

	fmt.Printf("Running in pong mode as %s@%s\n", a.Entity.DID.Id, p.Host.ID())
	fmt.Println("Press Ctrl-C to stop.")

	for {
		// A background loop that does nothing.
		// The ctx will never be cancelled, so this will run forever.
		<-ctx.Done()
		fmt.Println("Pong run loop cancelled, exiting...")
		return
	}
}

// ui.currentActorCtx, ui.currentActorCancel = context.WithCancel(context.Background())

// // Now that the UI is created, we can start the actor and subscribe to its events.
// go ui.a.Subscribe(ui.currentActorCtx, ui.a.Entity)

// // We want to handle envelopes for the actor, then deliver the messages
// // to the UI from the incoming envelopes.
// go ui.handleIncomingEnvelopes(ui.currentActorCtx, ui.a.Entity, ui.a)
// go ui.handleIncomingMessages(ui.currentActorCtx, ui.a.Entity)

// go actor.HelloWorld(ui.currentActorCtx, ui.a) // This wait a bit before sending the message.

// func (ui *ChatUI) Run() error {

// 	defer ui.end()

// 	// Now we can start continuous discovery in the background.
// 	fmt.Println("Starting discovery loop in the background....")
// 	go ui.p.StartDiscoveryLoop(context.Background())

// 	// The actor should just run in the background for ever.
// 	// It will handle incoming messages and envelopes.
// 	// It shouldn't change - ever.
// 	fmt.Println("Starting actor...")
// 	ui.startActor()

// 	// We must wait for this to finish.
// 	fmt.Printf("Entering %s ...\n", config.ActorLocation())
// 	err := ui.enterEntity(config.ActorLocation(), true)
// 	if err != nil {
// 		ui.displayStatusMessage(err.Error())
// 	}
// 	fmt.Printf("Entered %s\n", config.ActorLocation())

// 	fmt.Println("Starting event loop...")
// 	go ui.handleEvents()

// 	return ui.app.Run()

// }
