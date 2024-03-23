package config

import (
	"github.com/bahner/go-ma"
	"github.com/spf13/viper"
)

type APIStruct struct {
	Maddr string `yaml:"maddr"`
}

type APIConfigStruct struct {
	Api APIStruct `yaml:"api"`
}

func APIConfig() APIConfigStruct {
	viper.SetDefault("api.maddr", ma.DEFAULT_IPFS_API_MULTIADDR)

	return APIConfigStruct{
		Api: APIStruct{
			Maddr: APIAddr(),
		},
	}
}
func APIAddr() string {
	return viper.GetString("api.maddr")
}
