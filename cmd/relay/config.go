package main

import (
	"os"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	p2pDHT "github.com/libp2p/go-libp2p-kad-dht"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

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
		}}
}

func initConfig(name string) {

	// Always parse the flags first
	config.InitCommonFlags()
	pflag.Parse()
	config.SetProfile(name)
	config.Init()

	if config.GenerateFlag() {
		// Reinit logging to STDOUT
		log.SetOutput(os.Stdout)
		log.Info("Generating new node identity")
		node, err := config.GenerateNodeIdentity()
		if err != nil {
			log.Fatal(err)
		}
		relayConfig := configTemplate(node)
		config.Generate(relayConfig)
		os.Exit(0)
	}

	// This flag is dependent on the actor to be initialized to make sense.
	if config.ShowConfigFlag() {
		config.Print()
		os.Exit(0)
	}

}

func configTemplate(node string) map[string]interface{} {

	// Get the default settings as a map
	// Note: Viper does not have a built-in way to directly extract only the config
	// so we manually recreate the structure based on the config we have set.
	return map[string]interface{}{
		"db": map[string]interface{}{
			"dir": config.DefaultDbPath,
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
