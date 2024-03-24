package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	name = "robot"
)

func init() {
	pflag.String("openai-key", "", "The (paid) key to use with the OpenAI API")

	viper.BindPFlag("robot.openai.key", pflag.Lookup("openai-key"))
}

type RobotConfigStruct struct {
	OpenAIConfigStruct `yaml:"openai"`
}

type OpenAIConfigStruct struct {
	Key string `yaml:"key"`
}
type RobotConfig struct {
	actor.ActorConfig
	Robot RobotConfigStruct `yaml:"robot"`
}

func initConfig(name string) RobotConfig {

	actorConfig := actor.Config(name) // Assuming this returns an ActorConfigStruct

	// Create a new RobotConfig with the base config and the new key
	robotConfig := RobotConfig{
		ActorConfig: actorConfig,
		Robot: RobotConfigStruct{
			OpenAIConfigStruct: OpenAIConfigStruct{
				Key: openAIKey(),
			},
		},
	}
	r := robotConfig

	if config.GenerateFlag() {
		config.HandleGenerate(&r)
		os.Exit(0)
	}

	if config.ShowConfigFlag() {
		r.Print()
		os.Exit(0)
	}

	return r
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
