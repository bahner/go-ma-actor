package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type RelayConfig struct {
	API  config.APIConfigStruct  `yaml:"api"`
	DB   config.DBConfigStruct   `yaml:"db"`
	HTTP config.HTTPConfigStruct `yaml:"http"`
	Log  config.LogConfigStruct  `yaml:"log"`
	P2P  config.P2PConfigStruct  `yaml:"p2p"`
}

func Config(name string) RelayConfig {

	pflag.Parse()
	config.SetProfile(name)
	config.Init()

	c := RelayConfig{
		API:  config.APIConfig(),
		DB:   config.DBConfig(),
		HTTP: config.HTTPConfig(),
		Log:  config.LogConfig(),
		P2P:  config.P2PConfig(),
	}

	if config.GenerateFlag() {
		config.Save(&c)
	}

	return c
}

func (c *RelayConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *RelayConfig) Print() {
	config.Print(c)
}

func (c *RelayConfig) Save() error {
	return config.Save(c)
}
