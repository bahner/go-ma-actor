package main

import (
	"errors"
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma/did/doc"
	"github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultPongReply   = "Pong!"
	defaultFortuneMode = false
	pong               = "pong"
)

var defaultFortuneArgs = []string{"-s"}

func init() {
	pflag.String("pong-reply", defaultPongReply, "The message to send back to the sender")

	viper.BindPFlag("mode.pong.reply", pflag.Lookup("pong-reply"))
	viper.SetDefault("mode.pong.reply", defaultPongReply)

	viper.BindPFlag("mode.pong.fortune.enable", pflag.Lookup("pong-fortune"))
	viper.SetDefault("mode.pong.fortune.enable", defaultFortuneMode)

	viper.SetDefault("mode.pong.fortune.args", defaultFortuneArgs)
}

func p2pOptions() p2p.Options {
	return p2p.Options{
		DHT: []p2pDHT.Option{
			p2pDHT.Mode(p2pDHT.ModeServer),
		},
		P2P: []libp2p.Option{
			libp2p.DefaultTransports,
			libp2p.DefaultSecurity,
			libp2p.EnableRelay(),
			libp2p.EnableRelayService(),
			libp2p.Ping(true),
		}}
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
		actor, node := generateActorIdentitiesOrPanic(pong)
		pongConfig := configTemplate(actor, node)
		config.Generate(pongConfig)
		os.Exit(0)
	}

	config.InitActor()

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}

func pongFortuneMode() bool {
	return viper.GetBool("mode.pong.fortune.enable")
}

func pongFortuneArgs() []string {
	return viper.GetStringSlice("mode.pong.fortune.args")
}

func pongReply() string {
	return viper.GetString("mode.pong.reply")
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
			"nick":     pong,
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
			"pong": map[string]interface{}{
				"reply": pongReply(),
				"fortune": map[string]interface{}{
					"enable": pongFortuneMode(),
					"args":   pongFortuneArgs(),
				},
			},
		},
	}
}
