package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"

	log "github.com/sirupsen/logrus"
)

const defaultMsg = "yo"

func init() {
	pflag.String("msg", defaultMsg, "Message to send as a pong. For fun and identification.")
	viper.BindPFlag("pong.msg", pflag.Lookup("msg"))
	viper.SetDefault("pong.msg", defaultMsg)

}

func main() {

	pflag.Parse()

	config.Init("pong")
	config.InitLogging()
	config.InitP2P()
	config.InitActor()

	ctx := context.Background()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	// We need to discover peers before we can do anything else.
	p.DiscoverPeers()

	k := config.GetKeyset()
	e, err := entity.NewFromKeyset(k, k.DID.Fragment)
	if err != nil {
		log.Errorf("Error initializing actor: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	fmt.Printf("I am : %s\n", e.DID.String())
	fmt.Printf("My public key is: %s\n", p.Node.ID().String())

	// Now we can start continuous discovery in the background.
	go p.DiscoveryLoop(ctx)
	go handleEvents(ctx, e)

	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &entity.WebHandlerData{
		P2P:    p,
		Entity: e,
	}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s", config.GetHttpSocket())
	err = http.ListenAndServe(config.GetHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
