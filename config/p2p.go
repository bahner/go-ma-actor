package config

import (
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConnmgrLowWatermark  int           = 50
	defaultConnmgrHighWatermark int           = 100
	defaultConnmgrGracePeriod   time.Duration = time.Minute * 1

	defaultDiscoveryAdvertiseInterval time.Duration = time.Minute * 5
	defaultDiscoveryAdvertiseTTL      time.Duration = time.Minute * 60
	defaultDiscoveryAdvertiseLimit    int           = 100

	ALLOW_ALL_PEERS bool = true // Allow all peers by default. This is the norm for now. Use connmgr threshold and protection instead.

	defaultListenPort int    = 0
	fakeP2PIdentity   string = "NO_DEFAULT_NODE_IDENITY"
)

func init() {

	// P2P Settings

	// CONNMGR
	pflag.Duration("connmgr-grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
	pflag.Int("connmgr-high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
	pflag.Int("connmgr-low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")

	viper.BindPFlag("p2p.connmgr.grace-period", pflag.Lookup("connmgr-grace-period"))
	viper.BindPFlag("p2p.connmgr.high-watermark", pflag.Lookup("connmgr-high-watermark"))
	viper.BindPFlag("p2p.connmgr.low-watermark", pflag.Lookup("connmgr-low-watermark"))

	viper.SetDefault("p2p.connmgr.grace-period", defaultConnmgrGracePeriod)
	viper.SetDefault("p2p.connmgr.high-watermark", defaultConnmgrHighWatermark)
	viper.SetDefault("p2p.connmgr.low-watermark", defaultConnmgrLowWatermark)

	// DISCOVERY
	pflag.Int("discovery-advertise-limit", defaultDiscoveryAdvertiseLimit, "Limit for advertising peer discovery.")
	pflag.Duration("discovery-advertise-ttl", defaultDiscoveryAdvertiseTTL, "Hint of TimeToLive for advertising peer discovery.")
	pflag.Duration("discovery-advertise-interval", defaultDiscoveryAdvertiseInterval, "How often to advertise our presence to libp2p")

	viper.BindPFlag("p2p.discovery.advertise-interval", pflag.Lookup("discovery-advertise-interval"))
	viper.BindPFlag("p2p.discovery.advertise-limit", pflag.Lookup("discovery-advertise-limit"))
	viper.BindPFlag("p2p.discovery.advertise-ttl", pflag.Lookup("discovery-advertise-ttl"))

	viper.SetDefault("p2p.discovery.advertise-interval", defaultDiscoveryAdvertiseInterval)
	viper.SetDefault("p2p.discovery.advertise-limit", defaultDiscoveryAdvertiseLimit)
	viper.SetDefault("p2p.discovery.advertise-ttl", defaultDiscoveryAdvertiseTTL)

	// Port
	pflag.Int("port", defaultListenPort, "Port for libp2p node to listen on.")
	viper.BindPFlag("p2p.port", pflag.Lookup("port"))
	viper.SetDefault("p2p.port", defaultListenPort)

	// Identity
	viper.SetDefault("p2p.identity", fakeP2PIdentity)

}

// P2P Settings

func P2PIdentity() string {

	return viper.GetString("p2p.identity")
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

func P2PPort() int {
	return viper.GetInt("p2p.port")
}

// String functions

func P2PPortString() string {
	return strconv.Itoa(P2PPort())
}
