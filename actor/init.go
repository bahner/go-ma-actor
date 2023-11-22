package actor

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p/pubsub"
	log "github.com/sirupsen/logrus"
)

var (
	err error
	a   *Actor
)

// initActor initializes the actor. initKeyset and initP2P must've been run before this.
func init() {

	ctx := context.Background()
	k := config.GetKeyset()
	ps := pubsub.Get()

	log.Infof("Intializing actor with identity: %s", k.IPNSKey.DID)
	a, err = NewFromKeyset(ctx, ps, k, config.GetForcePublish())
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.Entity.DID.Fragment)

}

func GetActor() *Actor {
	return a
}
