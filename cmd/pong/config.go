package main

import (
	"errors"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	defaultPongReply   = "Pong!"
	defaultFortuneMode = false
	pong               = "pong"
	profile            = pong
)

var defaultFortuneArgs = []string{"-s"}

func init() {
	pflag.String("pong-reply", defaultPongReply, "The message to send back to the sender")

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", defaultPongReply)

	viper.BindPFlag("mode.pong.fortune.enable", pflag.Lookup("pong-fortune"))
	viper.SetDefault("mode.pong.fortune.enable", defaultFortuneMode)

	viper.SetDefault("mode.pong.fortune.args", defaultFortuneArgs)
}

type FortuneStruct struct {
	Enable bool     `yaml:"enable"`
	Args   []string `yaml:"args"`
}

type PongStruct struct {
	Reply   string        `yaml:"reply"`
	Fortune FortuneStruct `yaml:"fortune"`
}

type PongConfigStruct struct {
	Pong PongStruct `yaml:"pong"`
}

type PongConfig struct {
	API   config.APIConfigStruct   `yaml:"api"`
	Actor config.ActorConfigStruct `yaml:"actor"`
	DB    config.DBConfigStruct    `yaml:"db"`
	Log   config.LogConfigStruct   `yaml:"log"`
	Pong  PongConfigStruct         `yaml:"pong"`
}

func Config() PongConfig {

	config.ActorFlags()
	pflag.Parse()

	// Always parse the flags first
	config.SetProfile(profile)
	config.Init()

	return PongConfig{
		API:   config.APIConfig(),
		Actor: config.ActorConfig(),
		DB:    config.DBConfig(),
		Log:   config.LogConfig(),
		Pong: PongConfigStruct{
			Pong: PongStruct{
				Reply: pongReply(),
				Fortune: FortuneStruct{
					Enable: pongFortuneMode(),
					Args:   pongFortuneArgs()},
			},
		},
	}
}

func (c *PongConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *PongConfig) Print() {
	y, err := c.MarshalToYAML()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(y))
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

func generateActorIdentitiesOrPanic(name string) (string, string) {
	actor, node, err := config.GenerateActorIdentities(name)
	if err != nil {
		if errors.Is(err, doc.ErrAlreadyPublished) {
			log.Warnf("Actor document already published: %v", err)
		} else {
			log.Fatal(err)
		}
	}
	return actor, node
}
