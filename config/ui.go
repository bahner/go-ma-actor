package config

import (
	"github.com/spf13/viper"
)

func UIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
