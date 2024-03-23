package actor

import (
	"github.com/bahner/go-ma-actor/config"
	"gopkg.in/yaml.v2"
)

type ActorConfigStruct struct {
	Actor config.ActorConfigStruct `yaml:"actor"`
	API   config.APIConfigStruct   `yaml:"api"`
	DB    config.DBConfigStruct    `yaml:"db"`
	HTTP  config.HTTPConfigStruct  `yaml:"http"`
	Log   config.LogConfigStruct   `yaml:"log"`
	P2P   config.P2PConfigStruct   `yaml:"p2p"`
}

func Config(name string) ActorConfigStruct {

	config.SetProfile(name)

	return ActorConfigStruct{
		Actor: config.ActorConfig(),
		API:   config.APIConfig(),
		DB:    config.DBConfig(),
		HTTP:  config.HTTPConfig(),
		Log:   config.LogConfig(),
		P2P:   config.P2PConfig(),
	}
}

func (c *ActorConfigStruct) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *ActorConfigStruct) Print() {
	config.Print(c)
}

func (c *ActorConfigStruct) Save() error {
	return config.Save(c)
}
