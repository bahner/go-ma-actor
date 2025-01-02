package main

import (
	"os"
	"time"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	defautSpaceNodeName      = "space@localhost"
	defaultNodeCookie        = "spacecookie"
	defaultNodeName          = "pubsub@localhost"
	defaultNodeDebugInterval = time.Second * 60
	name                     = "node"
)

func init() {

	// Erlang node config
	pflag.String("spacenode", defautSpaceNodeName, "Name of the node running the actual SPACE")
	viper.BindPFlag("node.space", pflag.Lookup("spacenode"))
	viper.SetDefault("node.space", defautSpaceNodeName)

	pflag.String("nodecookie", defaultNodeCookie, "Secret shared between erlang nodes in the cluster")
	viper.BindPFlag("node.cookie", pflag.Lookup("nodecookie"))
	viper.SetDefault("node.cookie", defaultNodeCookie)

	pflag.String("nodename", defaultNodeName, "Name of the erlang node")
	viper.BindPFlag("node.name", pflag.Lookup("nodename"))
	viper.SetDefault("node.name", defaultNodeName)

	pflag.Duration("node_debug_interval", defaultNodeDebugInterval, "Interval for debug output")
	viper.BindPFlag("node.debug_interval", pflag.Lookup("_node_debug_interval"))
	viper.SetDefault("node.debug_interval", defaultNodeDebugInterval)

}

type NodeConfigStruct struct {
	Cookie        string        `yaml:"cookie"`
	Name          string        `yaml:"name"`
	Space         string        `yaml:"space"`
	DebugInterval time.Duration `yaml:"debug-interval"`
}

type NodeConfig struct {
	Actor config.ActorConfig `yaml:"actor"`
	DB    config.DBConfig    `yaml:"db"`
	HTTP  config.HTTPConfig  `yaml:"http"`
	Log   config.LogConfig   `yaml:"log"`
	Node  NodeConfigStruct   `yaml:"node"`
	P2P   config.P2PConfig   `yaml:"p2p"`
}

func Config(defaultProfileName string) NodeConfig {

	config.SetDefaultProfileName(defaultProfileName)
	actor.Config()

	c := NodeConfig{
		Actor: config.Actor(),
		DB:    config.DB(),
		HTTP:  config.HTTP(),
		Node: NodeConfigStruct{
			Cookie:        NodeCookie(),
			Name:          NodeName(),
			Space:         NodeSpace(),
			DebugInterval: NodeDebugInterval(),
		},
		Log: config.Log(),
		P2P: config.P2P(),
	}

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

func (c *NodeConfig) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *NodeConfig) Print() {
	config.Print(c)
}

func (c *NodeConfig) Save() error {
	return config.Save(c)
}

func NodeSpace() string {
	return viper.GetString("node.space")
}

func NodeCookie() string {
	return viper.GetString("node.cookie")
}

func NodeName() string {
	return viper.GetString("node.name")
}

func NodeDebugInterval() time.Duration {
	return viper.GetDuration("node.debug_interval")
}
