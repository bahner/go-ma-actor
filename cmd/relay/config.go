package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type RelayConfig struct {
	API  config.APIConfig  `yaml:"api"`
	DB   config.DBConfig   `yaml:"db"`
	HTTP config.HTTPConfig `yaml:"http"`
	Log  config.LogConfig  `yaml:"log"`
	P2P  config.P2PConfig  `yaml:"p2p"`
}

func initConfig(defaultProfileName string) RelayConfig {

	pflag.Parse()
	config.SetDefaultProfileName(defaultProfileName)
	config.Init()

	c := RelayConfig{
		API:  config.API(),
		DB:   config.DB(),
		HTTP: config.HTTP(),
		Log:  config.Log(),
		P2P:  config.P2P(),
	}

	if config.GenerateFlag() {
		config.Generate(&c)
	}

	if config.ShowConfigFlag() {
		c.Print()
	}

	if config.ShowConfigFlag() || config.GenerateFlag() {
		os.Exit(0)
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
