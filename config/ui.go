package config

import (
	"github.com/spf13/viper"
)

func GetUIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
