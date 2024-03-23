package config

import (
	"github.com/spf13/viper"
)

type UIStruct struct {
	PeerslistWidth int `yaml:"peerslist-width"`
}

type UIConfigStruct struct {
	UI UIStruct `yaml:"ui"`
}

func UIConfig() UIConfigStruct {
	return UIConfigStruct{
		UI: UIStruct{
			PeerslistWidth: UIPeerslistWidth(),
		},
	}
}

func UIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
