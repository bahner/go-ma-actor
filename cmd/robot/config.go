package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	name = "robot"
)

func init() {
	pflag.String("openai-key", "", "The (paid) key to use with the OpenAI API")

	viper.BindPFlag("mode.openai.key", pflag.Lookup("openai-key"))
}

func initConfig(profile string) {

	// Always parse the flags first
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.SetProfile(profile)
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node, err := config.GenerateActorIdentities(name)
		if err != nil {
			log.Fatalf("Failed to generate identities: %v", err)
		}
		openaiConfig := configTemplate(actor, node)
		config.Generate(openaiConfig)
		os.Exit(0)
	}

	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}
func openAIKey() string {
	return viper.GetString("mode.openai.key")
}

func configTemplate(identity string, node string) map[string]interface{} {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	return map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"nick":     name,
		},
		"db": map[string]interface{}{
			"file": config.DefaultDbFile,
		},
		"log": map[string]interface{}{
			"level": config.LogLevel(),
			"file":  config.LogFile(),
		},
		// NB! This is a cross over from go-ma
		"api": map[string]interface{}{
			// This must be set corretly for generation to work
			"maddr": viper.GetString("api.maddr"),
		},
		"http": map[string]interface{}{
			"socket": config.HttpSocket(),
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
		"mode": map[string]interface{}{
			"openai": map[string]interface{}{
				"key": openAIKey(),
			},
		},
	}
}
