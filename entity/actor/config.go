package actor

import (
	"github.com/bahner/go-ma-actor/config"
	"gopkg.in/yaml.v2"
)

type ActorConfig struct {
	Actor config.ActorConfig `yaml:"actor"`
	API   config.APIConfig   `yaml:"api"`
	DB    config.DBConfig    `yaml:"db"`
	HTTP  config.HTTPConfig  `yaml:"http"`
	Log   config.LogConfig   `yaml:"log"`
	P2P   config.P2PConfig   `yaml:"p2p"`
}

// This is an all-inclusive configuration function that sets up the configuration for the actor.
// flags and everything. It is used in the main function of simple actors programmes.
// It also parses the common flags.
func Config() ActorConfig {

	config.ActorFlagsetParse(true)
	config.Init()

	return ActorConfig{
		Actor: config.Actor(),
		API:   config.API(),
		DB:    config.DB(),
		HTTP:  config.HTTP(),
		Log:   config.Log(),
		P2P:   config.P2P(),
	}
}

func (c *ActorConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *ActorConfig) Print() {
	config.Print(c)
}

func (c *ActorConfig) Save() error {
	return config.Save(c)
}
