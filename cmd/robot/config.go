package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func init() {

	viper.BindPFlag("robot.openai.key", pflag.Lookup("openai-key"))
}

type RobotConfigStruct struct {
	OpenAIConfigStruct `yaml:"openai"`
}

type OpenAIConfigStruct struct {
	Key string `yaml:"key"`
}
type RobotConfig struct {
	Actor config.ActorConfig `yaml:"actor"`
	DB    config.DBConfig    `yaml:"db"`
	HTTP  config.HTTPConfig  `yaml:"http"`
	Log   config.LogConfig   `yaml:"log"`
	P2P   config.P2PConfig   `yaml:"p2p"`
	Robot RobotConfigStruct  `yaml:"robot"`
}

func initConfig(defaultProfileName string) RobotConfig {

	config.SetDefaultProfileName(defaultProfileName)
	actor.Config()

	// Create a new RobotConfig with the base config and the new key
	robotConfig := RobotConfig{
		Actor: config.Actor(),
		DB:    config.DB(),
		HTTP:  config.HTTP(),
		Log:   config.Log(),
		P2P:   config.P2P(),
		Robot: RobotConfigStruct{
			OpenAIConfigStruct: OpenAIConfigStruct{
				Key: openAIKey(),
			},
		},
	}
	c := robotConfig

	if config.GenerateFlag() {
		config.GenerateConfig(&c)
	}

	if config.ShowConfigFlag() {
		c.Print()
	}

	if config.ShowConfigFlag() || config.GenerateFlag() {
		os.Exit(0)
	}

	return c
}

func (c RobotConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c RobotConfig) Print() {
	config.Print(c)
}

func (c RobotConfig) Save() error {
	return config.Save(c)
}

func openAIKey() string {
	return viper.GetString("robot.openai.key")
}
