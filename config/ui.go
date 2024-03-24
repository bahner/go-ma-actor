package config

import (
	"github.com/spf13/viper"
)

type UIStruct struct {
	PeerslistWidth int `yaml:"peerslist-width"`
}

type UIConfig struct {
	UI UIStruct `yaml:"ui"`
}

func InitUIConfig() UIConfig {
	return UIConfig{
		UI: UIStruct{
			PeerslistWidth: UIPeerslistWidth(),
		},
	}
}

func UIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
