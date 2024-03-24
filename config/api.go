package config

import (
	"github.com/bahner/go-ma"
	"github.com/spf13/viper"
)

type APIConfig struct {
	Maddr string `yaml:"maddr"`
}

func API() APIConfig {
	viper.SetDefault("api.maddr", ma.DEFAULT_IPFS_API_MULTIADDR)

	return APIConfig{
		Maddr: APIAddr(),
	}
}
func APIAddr() string {
	return viper.GetString("api.maddr")
}
