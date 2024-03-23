package main

import (
	"time"

	"github.com/bahner/go-ma-actor/config"
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

type NodeStruct struct {
	Cookie        string        `yaml:"cookie"`
	Name          string        `yaml:"name"`
	Space         string        `yaml:"space"`
	DebugInterval time.Duration `yaml:"debug-interval"`
}

type NodeConfigStruct struct {
	Node NodeStruct `yaml:"node"`
}

type NodeConfig struct {
	Actor config.ActorConfigStruct `yaml:"actor"`
	API   config.APIConfigStruct   `yaml:"api"`
	DB    config.DBConfigStruct    `yaml:"db"`
	HTTP  config.HTTPConfigStruct  `yaml:"http"`
	Node  NodeConfigStruct         `yaml:"node"`
	Log   config.LogConfigStruct   `yaml:"log"`
	P2P   config.P2PConfigStruct   `yaml:"p2p"`
}

func Config(name string) NodeConfig {

	config.ActorFlags()
	pflag.Parse()

	config.SetProfile(name)

	c := NodeConfig{
		Actor: config.ActorConfig(),
		API:   config.APIConfig(),
		DB:    config.DBConfig(),
		HTTP:  config.HTTPConfig(),
		Node: NodeConfigStruct{
			Node: NodeStruct{
				Cookie:        NodeCookie(),
				Name:          NodeName(),
				Space:         NodeSpace(),
				DebugInterval: NodeDebugInterval(),
			},
		},
		Log: config.LogConfig(),
		P2P: config.P2PConfig(),
	}

	if config.GenerateFlag() {
		config.Save(&c)
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
