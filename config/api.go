package config

import (
	"github.com/bahner/go-ma"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("api.maddr", ma.DEFAULT_IPFS_API_MULTIADDR)
}
