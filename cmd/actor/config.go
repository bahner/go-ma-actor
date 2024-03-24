package main

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
)

func initConfig(profileName string) actor.ActorConfig {

	c := actor.Config(profileName)
	config.HandleGenerate(&c)

	return c
}
