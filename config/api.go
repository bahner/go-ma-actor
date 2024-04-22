package config

import (
	"github.com/spf13/viper"
)

type APIConfig struct {
	Maddr string `yaml:"maddr"`
}

func API() APIConfig {
	return APIConfig{
		Maddr: APIAddr(),
	}
}
func APIAddr() string {
	return viper.GetString("api.maddr")
}
