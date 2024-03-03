package config

import (
	"fmt"
	"os"

	"github.com/bahner/go-ma"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Well, actually find a nice one for this!
const defaultHome = "did:ma:k2k4r8kzkhamrqz9m5yy0tihj1fso3t6znnuidu00dbtnh3plazatrfw#pong"

func generateConfigFile(identity string, node string) {

	var nick string

	if identity == fakeActorIdentity {
		nick = defaultActor
	} else {
		nick = viper.GetString("actor.nick")
	}

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	config := map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"home":     defaultHome,
			"nick":     nick,
		},
		"db": map[string]interface{}{
			"file": defaultDbFile,
		},
		"log": map[string]interface{}{
			"level": defaultLogLevel,
			"file":  defaultLogfile,
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
			"relay": defaultRelayMode,
			"pong": map[string]interface{}{
				"reply":   DefaultPongReply,
				"enabled": defaultPongMode,
			},
		},
	}

	// Convert the map of defaults to YAML
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Begin the output
	generateFlag, err := pflag.CommandLine.GetBool("generate")
	if err != nil {
		log.Warnf("config.init: %v", err)
	}
	if generateFlag {
		writeGeneratedConfigFile(configYAML)
	} else {
		fmt.Println(string(configYAML))
	}
}

// Write the generated config to the correct file
// NB! This fails fatally in case of an error.
func writeGeneratedConfigFile(content []byte) {

	filePath := configFile()
	var errMsg string

	// Try to open the file with flags to ensure it does not overwrite an existing file.
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			errMsg = fmt.Sprintf("File %s already exists.", filePath)
		} else {
			errMsg = fmt.Sprintf("Failed to open file: %v", err)
		}
		log.Fatalf(errMsg)
	}
	defer file.Close()

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Infof("Generated actor config file %s", filePath)
}
