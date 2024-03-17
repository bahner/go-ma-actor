package main

import (
	"github.com/bahner/go-ma-actor/config"
)

func configTemplate(node string) map[string]interface{} {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	return map[string]interface{}{
		"db": map[string]interface{}{
			"file": config.DefaultDbFile,
		},
		"log": map[string]interface{}{
			"level": config.LogLevel(),
			"file":  config.LogFile(),
		},
		"http": map[string]interface{}{
			"socket":  config.HttpSocket(),
			"refresh": config.HttpRefresh(),
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     config.P2PPort(),
			"connmgr": map[string]interface{}{
				"low-watermark":  config.P2PConnmgrLowWatermark(),
				"high-watermark": config.P2PConnmgrHighWatermark(),
				"grace-period":   config.P2PConnMgrGracePeriod(),
			},
			"discovery": map[string]interface{}{
				"advertise-ttl":      config.P2PDiscoveryAdvertiseTTL(),
				"advertise-limit":    config.P2PDiscoveryAdvertiseLimit(),
				"advertise-interval": config.P2PDiscoveryAdvertiseInterval(),
				"dht":                config.P2PDiscoveryDHT(),
				"mdns":               config.P2PDiscoveryMDNS(),
			},
		},
	}

}
