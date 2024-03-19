package actor

import (
	"errors"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma/did/doc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig() {

	// Always parse the flags first
	config.InitCommonFlags()
	config.InitActorFlags()
	pflag.Parse()
	config.SetProfile(config.Profile())
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new actor and node identity")
		actor, node := generateActorIdentitiesOrPanic(config.Profile())
		actorConfig := configTemplate(actor, node)
		config.Generate(actorConfig)
		os.Exit(0)
	}

	// At this point an actor *must* be initialized
	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}

func generateActorIdentitiesOrPanic(name string) (string, string) {
	actor, node, err := config.GenerateActorIdentities(name)
	if err != nil {
		if errors.Is(err, doc.ErrAlreadyPublished) {
			log.Warnf("Actor document already published: %v", err)
		} else {
			log.Fatal(err)
		}
	}
	return actor, node
}

func configTemplate(identity string, node string) map[string]interface{} {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	return map[string]interface{}{
		"actor": map[string]interface{}{
			"identity": identity,
			"location": config.ActorLocation(),
			"nick":     config.ActorNick(),
		},
		"db": map[string]interface{}{
			"file": config.DefaultDbFile,
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
	}
}
