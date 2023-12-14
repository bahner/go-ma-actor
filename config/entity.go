package config

import (
	"flag"

	"go.deanishe.net/env"
)

var (
	entity = flag.String("entity", env.Get(GO_MA_ACTOR_ENTITY_VAR, defaultEntity),
		"DID of the entity to communicate with. You can use environment variable "+GO_MA_ACTOR_ENTITY_VAR+" to set this.")
)

func GetEntity() string {
	return *entity
}
