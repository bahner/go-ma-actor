package config

import (
	"time"

	"github.com/spf13/viper"
)

func UIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
func UIPeersRefreshInterval() time.Duration {
	return viper.GetDuration("ui.refresh")
}
