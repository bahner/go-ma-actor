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
	pflag.String("pong-reply", defaultPongReply, "The message to send back to the sender")

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", defaultPongReply)

	viper.BindPFlag("mode.pong.fortune.enable", pflag.Lookup("pong-fortune"))
	viper.SetDefault("mode.pong.fortune.enable", defaultFortuneMode)

	viper.SetDefault("mode.pong.fortune.args", defaultFortuneArgs)
}

type PongFortuneStruct struct {
	Enable bool     `yaml:"enable"`
	Args   []string `yaml:"args"`
}

type PongConfigStruct struct {
	Reply   string            `yaml:"reply"`
	Fortune PongFortuneStruct `yaml:"fortune"`
}

type PongConfig struct {
	Actor config.ActorConfig `yaml:"actor"`
	API   config.APIConfig   `yaml:"api"`
	DB    config.DBConfig    `yaml:"db"`
	HTTP  config.HTTPConfig  `yaml:"http"`
	Log   config.LogConfig   `yaml:"log"`
	P2P   config.P2PConfig   `yaml:"p2p"`
	Pong  PongConfigStruct   `yaml:"pong"`
}

func initConfig(defaultProfileName string) PongConfig {

	config.SetDefaultProfileName(defaultProfileName)
	actor.Config()

	c := PongConfig{
		Actor: config.Actor(),
		API:   config.API(),
		DB:    config.DB(),
		HTTP:  config.HTTP(),
		Log:   config.Log(),
		P2P:   config.P2P(),
		Pong: PongConfigStruct{
			Reply: pongReply(),
			Fortune: PongFortuneStruct{
				Enable: pongFortuneMode(),
				Args:   pongFortuneArgs()},
		},
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

func (c *PongConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *PongConfig) Print() {
	config.Print(c)
}

func (c *PongConfig) Save() error {
	return config.Save(c)
}

func pongFortuneMode() bool {
	return viper.GetBool("mode.pong.fortune.enable")
}

func pongFortuneArgs() []string {
	return viper.GetStringSlice("mode.pong.fortune.args")
}

func pongReply() string {
	return viper.GetString("mode.pong.reply")
}
