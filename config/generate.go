package config

import (
	"fmt"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma/key/set"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Well, actually find a nice one for this!
const defaultHome = "did:ma:k2k4r8kzkhamrqz9m5yy0tihj1fso3t6znnuidu00dbtnh3plazatrfw#pong"

func generateConfigFile(actor string, node string) {

	ks, err := set.Unpack(actor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	logFile, err := homedir.Expand("~/.ma/" + viper.GetString("actor.nick") + ".log")
	if err != nil {
		log.Fatalf("generateConfigFile: %v", err)
	}

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the defaults
	// so we manually recreate the structure based on the defaults we have set.
	defaults := map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": actor,
			"home":     defaultHome,
			"nick":     viper.GetString("actor.nick"),
		},
		"aliases": defaultAliasesDB,
		"log": map[string]interface{}{
			"level": defaultLogLevel,
			"file":  logFile,
		},
		// NB! This is a cross over from go-ma
		"api": map[string]interface{}{
			"maddr": ma.DEFAULT_IPFS_API_MULTIADDR,
		},
		"http": map[string]interface{}{
			"socket": defaultHttpSocket,
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     defaultListenPort,
			"connmgr": map[string]interface{}{
				"low-watermark":  defaultConnmgrLowWatermark,
				"high-watermark": defaultConnmgrHighWatermark,
				"grace-period":   defaultConnmgrGracePeriod,
			},
			"discovery-retry":   defaultDiscoveryRetryInterval,
			"discovery-timeout": defaultDiscoveryTimeout,
		},
		"mode": map[string]interface{}{
			"debug": defaultDebugMode,
			"relay": defaultRelayMode,
			"pong": map[string]interface{}{
				"reply":   DefaultPongReply,
				"enabled": defaultPongMode,
			},
		},
	}

	// Convert the map of defaults to YAML
	defaultsYAML, err := yaml.Marshal(defaults)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the YAML defaults
	fmt.Println("# " + ks.DID.Id)
	fmt.Println(string(defaultsYAML))
}
