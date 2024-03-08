package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func generateActorConfigFile(identity string, node string) {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	config := map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"location": ActorLocation(),
			"nick":     ActorNick(),
		},
		"db": map[string]interface{}{
			"file": DefaultDbFile,
		},
		// Use default log settings, so as not to pick up debug log settings
		"log": map[string]interface{}{
			"level": viper.GetString("log.level"),
			"file":  viper.GetString("log.file"),
		},
		// NB! This is a cross over from go-ma
		"api": map[string]interface{}{
			// This must be set corretly for generation to work
			"maddr": viper.GetString("api.maddr"),
		},
		"http": map[string]interface{}{
			"socket": HttpSocket(),
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     P2PPort(),
			"connmgr": map[string]interface{}{
				"low-watermark":  P2PConnmgrLowWatermark(),
				"high-watermark": P2PConnmgrHighWatermark(),
				"grace-period":   P2PConnMgrGracePeriod(),
			},
			"discovery": map[string]interface{}{
				"advertise-ttl":   P2PDiscoveryAdvertiseTTL(),
				"advertise-limit": P2PDiscoveryAdvertiseLimit(),
				"allow-all":       P2PDiscoveryAllowAll(),
			},
		},
	}

	// Convert the map of defaults to YAML
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	if generateFlag() {
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

	// Determine the file open flags based on the forceFlag
	var flags int
	if forceFlag() {
		// Allow overwrite
		log.Warnf("Force flag set, overwriting existing config file %s", filePath)
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	} else {
		// Prevent overwrite
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	file, err := os.OpenFile(filePath, flags, configFileMode)
	if err != nil {
		if os.IsExist(err) {
			errMsg = fmt.Sprintf("File %s already exists.", filePath)
		} else {
			errMsg = fmt.Sprintf("Failed to open file: %v", err)
		}
		panic(errMsg)
	}
	defer file.Close()

	// Write content to file.
	if _, err := file.Write(content); err != nil {
		panic(fmt.Sprintf("Failed to write to file: %v", err))
	}

	log.Printf("Generated config file %s", filePath)
}

func generatePongConfigFile(identity string, node string) {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	config := map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"nick":     pong,
		},
		"db": map[string]interface{}{
			"file": DefaultDbFile,
		},
		"log": map[string]interface{}{
			"level": LogLevel(),
			"file":  LogFile(),
		},
		// NB! This is a cross over from go-ma
		"api": map[string]interface{}{
			// This must be set corretly for generation to work
			"maddr": viper.GetString("api.maddr"),
		},
		"http": map[string]interface{}{
			"socket": HttpSocket(),
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     P2PPort(),
			"connmgr": map[string]interface{}{
				"low-watermark":  P2PConnmgrLowWatermark(),
				"high-watermark": P2PConnmgrHighWatermark(),
				"grace-period":   P2PConnMgrGracePeriod(),
			},
			"discovery": map[string]interface{}{
				"advertise-ttl":   P2PDiscoveryAdvertiseTTL(),
				"advertise-limit": P2PDiscoveryAdvertiseLimit(),
				"allow-all":       P2PDiscoveryAllowAll(),
			},
		},
		"mode": map[string]interface{}{
			"pong": map[string]interface{}{
				"reply": DEFAULT_PONG_REPLY,
			},
		},
	}

	// Convert the map of defaults to YAML
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	if generateFlag() {
		writeGeneratedConfigFile(configYAML)
	} else {
		fmt.Println(string(configYAML))
	}
}

func generateRelayConfigFile(node string) {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	config := map[string]interface{}{
		"db": map[string]interface{}{
			"file": DefaultDbFile,
		},
		"log": map[string]interface{}{
			"level": LogLevel(),
			"file":  LogFile(),
		},
		"http": map[string]interface{}{
			"socket": HttpSocket(),
		},
		"p2p": map[string]interface{}{
			"identity": node,
			"port":     P2PPort(),
			"connmgr": map[string]interface{}{
				"low-watermark":  P2PConnmgrLowWatermark(),
				"high-watermark": P2PConnmgrHighWatermark(),
				"grace-period":   P2PConnMgrGracePeriod(),
			},
			"discovery": map[string]interface{}{
				"advertise-ttl":   P2PDiscoveryAdvertiseTTL(),
				"advertise-limit": P2PDiscoveryAdvertiseLimit(),
				"allow-all":       P2PDiscoveryAllowAll(),
			},
		},
	}

	// Convert the map of defaults to YAML
	configYAML, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	if generateFlag() {
		writeGeneratedConfigFile(configYAML)
	} else {
		fmt.Println(string(configYAML))
	}
}
