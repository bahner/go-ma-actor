package ui

import (
	"time"

	"github.com/spf13/viper"
)

func getUIPeerslistWidth() int {
	return viper.GetInt("ui.peerslist-width")
}
func getUIPeersRefreshInterval() time.Duration {
	return viper.GetDuration("ui.refresh")
}
