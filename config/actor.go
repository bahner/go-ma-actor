package config

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/actor"
	log "github.com/sirupsen/logrus"
)

var (
	a *actor.Actor
)

// initActor initializes the actor. initKeyset and initP2P must've been run before this.
func initActor() {

	ctx := context.Background()
	ps := GetPubSub()
	k := GetKeyset()

	log.Infof("Intializing actor with identity: %s", k.IPNSKey.DID)
	a, err = actor.NewFromKeyset(ctx, ps, k, *forcePublish)
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}
	log.Infof("Actor initialized: %s", a.Entity.DID.Fragment)

}

func GetActor() *actor.Actor {
	return a
}
