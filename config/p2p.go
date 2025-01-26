package config

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConnmgrLowWatermark  int           = 10
	defaultConnmgrHighWatermark int           = 30
	defaultConnmgrGracePeriod   time.Duration = time.Minute * 2

	defaultDiscoveryAdvertiseInterval time.Duration = time.Minute * 5
	defaultDiscoveryAdvertiseTTL      time.Duration = time.Minute * 60
	defaultDiscoveryAdvertiseLimit    int           = 100

	ALLOW_ALL_PEERS bool = true // Allow all peers by default. This is the norm for now. Use connmgr threshold and protection instead.

	defaultListenPort   int    = 0
	defaultWSListenPort int    = 0
	fakeP2PIdentity     string = "NO_DEFAULT_NODE_IDENITY"
	defaultDHT          bool   = true
	defaultMDNS         bool   = true
)

var (
	p2pFlagset   = pflag.NewFlagSet("p2p", pflag.ExitOnError)
	p2pFlagsOnce sync.Once
)

func initP2PFlagset() {

	p2pFlagsOnce.Do(func() {

		p2pFlagset.Bool("dht", defaultDHT, "Whether to discover using DHT")
		p2pFlagset.Bool("mdns", defaultMDNS, "Whether to discover using MDNS")
		p2pFlagset.Duration("connmgr-grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
		p2pFlagset.Duration("discovery-advertise-interval", defaultDiscoveryAdvertiseInterval, "How often to advertise our presence to libp2p")
		p2pFlagset.Duration("discovery-advertise-ttl", defaultDiscoveryAdvertiseTTL, "Hint of TimeToLive for advertising peer discovery.")
		p2pFlagset.Int("connmgr-high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
		p2pFlagset.Int("connmgr-low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")
		p2pFlagset.Int("discovery-advertise-limit", defaultDiscoveryAdvertiseLimit, "Limit for advertising peer discovery.")
		p2pFlagset.Int("port", defaultListenPort, "Port for tcp to listen on. Only used for generating listen addresses, not runtime")
		p2pFlagset.Int("port-ws", defaultWSListenPort, "Port for ws to listen on. Only used for generating listen addresses, not runtime")

		// Bind p2pFlagss
		viper.BindPFlag("p2p.connmgr.grace-period", p2pFlagset.Lookup("connmgr-grace-period"))
		viper.BindPFlag("p2p.connmgr.high-watermark", p2pFlagset.Lookup("connmgr-high-watermark"))
		viper.BindPFlag("p2p.connmgr.low-watermark", p2pFlagset.Lookup("connmgr-low-watermark"))
		viper.BindPFlag("p2p.discovery.advertise-interval", p2pFlagset.Lookup("discovery-advertise-interval"))
		viper.BindPFlag("p2p.discovery.advertise-limit", p2pFlagset.Lookup("discovery-advertise-limit"))
		viper.BindPFlag("p2p.discovery.advertise-ttl", p2pFlagset.Lookup("discovery-advertise-ttl"))
		viper.BindPFlag("p2p.discovery.dht", p2pFlagset.Lookup("dht"))
		viper.BindPFlag("p2p.discovery.mdns", p2pFlagset.Lookup("mdns"))
		viper.BindPFlag("p2p.port", p2pFlagset.Lookup("port"))
		viper.BindPFlag("p2p.port-ws", p2pFlagset.Lookup("port-ws"))

		viper.SetDefault("p2p.connmgr.grace-period", defaultConnmgrGracePeriod)
		viper.SetDefault("p2p.connmgr.high-watermark", defaultConnmgrHighWatermark)
		viper.SetDefault("p2p.connmgr.low-watermark", defaultConnmgrLowWatermark)
		viper.SetDefault("p2p.discovery.advertise-interval", defaultDiscoveryAdvertiseInterval)
		viper.SetDefault("p2p.discovery.advertise-limit", defaultDiscoveryAdvertiseLimit)
		viper.SetDefault("p2p.discovery.advertise-ttl", defaultDiscoveryAdvertiseTTL)
		viper.SetDefault("p2p.discovery.dht", defaultDHT)
		viper.SetDefault("p2p.discovery.mdns", defaultMDNS)
		viper.SetDefault("p2p.port", defaultListenPort)
		viper.SetDefault("p2p.port-ws", defaultWSListenPort)

		if HelpNeeded() {
			fmt.Println("P2P Flags:")
			p2pFlagset.PrintDefaults()
		}
	})

}

type ConnmgrStruct struct {
	LowWatermark  int           `yaml:"low-watermark"`
	HighWatermark int           `yaml:"high-watermark"`
	GracePeriod   time.Duration `yaml:"grace-period"`
}

type DiscoveryStruct struct {
	AdvertiseInterval time.Duration `yaml:"advertise-interval"`
	AdvertiseTTL      time.Duration `yaml:"advertise-ttl"`
	AdvertiseLimit    int           `yaml:"advertise-limit"`
	DHT               bool          `yaml:"dht"`
	MDNS              bool          `yaml:"mdns"`
}

type P2PConfig struct {
	Maddrs    []string        `yaml:"maddrs"`
	Connmgr   ConnmgrStruct   `yaml:"connmgr"`
	Discovery DiscoveryStruct `yaml:"discovery"`
}

func P2P() P2PConfig {
	viper.SetDefault("p2p.identity", fakeP2PIdentity)

	return P2PConfig{
		Maddrs: generateListenAddrs(),
		Connmgr: ConnmgrStruct{
			LowWatermark:  P2PConnmgrLowWatermark(),
			HighWatermark: P2PConnmgrHighWatermark(),
			GracePeriod:   P2PConnMgrGracePeriod()},
		Discovery: DiscoveryStruct{
			AdvertiseInterval: P2PDiscoveryAdvertiseInterval(),
			AdvertiseTTL:      P2PDiscoveryAdvertiseTTL(),
			AdvertiseLimit:    P2PDiscoveryAdvertiseLimit(),
			DHT:               P2PDiscoveryDHT(),
			MDNS:              P2PDiscoveryMDNS()},
	}
}

func P2PDiscoveryAdvertiseInterval() time.Duration {
	return viper.GetDuration("p2p.discovery.advertise-interval")
}

func P2PDiscoveryAdvertiseTTL() time.Duration {
	return viper.GetDuration("p2p.discovery.advertise-ttl")
}

func P2PDiscoveryAdvertiseLimit() int {
	return viper.GetInt("p2p.discovery.advertise-limit")
}

func P2PConnmgrLowWatermark() int {
	return viper.GetInt("p2p.connmgr.low-watermark")
}

func P2PConnmgrHighWatermark() int {
	return viper.GetInt("p2p.connmgr.high-watermark")
}

func P2PConnMgrGracePeriod() time.Duration {
	return viper.GetDuration("p2p.connmgr.grace-period")
}

func P2PDiscoveryDHT() bool {
	return viper.GetBool("p2p.discovery.dht")
}

func P2PDiscoveryMDNS() bool {
	return viper.GetBool("p2p.discovery.mdns")
}

func P2PMaddrs() []string {
	return viper.GetStringSlice("p2p.maddrs")
}

func port() int {
	return viper.GetInt("p2p.port")
}

func portws() int {
	return viper.GetInt("p2p.port-ws")
}

func generateListenAddrs() []string {
	port := strconv.Itoa(port())
	portws := strconv.Itoa(portws())

	return []string{
		// Default addresses
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip6/::/tcp/" + port,
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",

		// WebSocket addresses
		"/ip4/0.0.0.0/tcp/" + portws + "/ws",
		"/ip6/::/tcp/" + portws + "/ws",
	}
}
