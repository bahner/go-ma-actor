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

	defaultDiscoveryAdvertiseTTL   time.Duration = time.Minute * 60
	defaultDiscoveryAdvertiseLimit int           = 100
	DEFAULT_ALLOW_ALL              bool          = true // Allow all peers by default. This is the norm for now. Use connmgr threshold and protection instead.

	defaultListenPort int    = 0
	fakeP2PIdentity   string = "NO_DEFAULT_NODE_IDENITY"
)

func init() {

	// P2P Settings

	// CONNMGR
	pflag.Duration("p2p-connmgr-grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
	pflag.Int("p2p-connmgr-high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
	pflag.Int("p2p-connmgr-low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")

	viper.BindPFlag("p2p.connmgr.grace-period", pflag.Lookup("p2p-connmgr-grace-period"))
	viper.BindPFlag("p2p.connmgr.high-watermark", pflag.Lookup("p2p-connmgr-high-watermark"))
	viper.BindPFlag("p2p.connmgr.low-watermark", pflag.Lookup("p2p-connmgr-low-watermark"))

	viper.SetDefault("p2p.connmgr.grace-period", defaultConnmgrGracePeriod)
	viper.SetDefault("p2p.connmgr.high-watermark", defaultConnmgrHighWatermark)
	viper.SetDefault("p2p.connmgr.low-watermark", defaultConnmgrLowWatermark)

	// DISCOVERY
	pflag.Int("p2p-discovery-advertise-limit", defaultDiscoveryAdvertiseLimit, "Limit for advertising peer discovery.")
	pflag.Duration("p2p-discovery-advertise-ttl", defaultDiscoveryAdvertiseTTL, "Hint o TimeToLive for advertising peer discovery.")
	pflag.Bool("p2p-discovery-allow-all", DEFAULT_ALLOW_ALL, "Number of concurrent peer discovery routines.")

	viper.BindPFlag("p2p.discovery.advertise-limit", pflag.Lookup("p2p-discovery-advertise-limit"))
	viper.BindPFlag("p2p.discovery.advertise-ttl", pflag.Lookup("p2p-discovery-advertise-ttl"))
	viper.BindPFlag("p2p.discovery.allow-all", pflag.Lookup("p2p-discovery-allow-all"))

	viper.SetDefault("p2p.discovery.advertise-limit", defaultDiscoveryAdvertiseLimit)
	viper.SetDefault("p2p.discovery.advertise-ttl", defaultDiscoveryAdvertiseTTL)
	viper.SetDefault("p2p.discovery.allow-all", DEFAULT_ALLOW_ALL)

	// Port
	pflag.Int("p2p-port", defaultListenPort, "Port for libp2p node to listen on.")
	viper.BindPFlag("p2p.port", pflag.Lookup("p2p-port"))
	viper.SetDefault("p2p.port", defaultListenPort)

	// Identity
	viper.SetDefault("p2p.identity", fakeP2PIdentity)

}

// P2P Settings

func P2PIdentity() string {

	return viper.GetString("p2p.identity")
}

func P2PDiscoveryAdvertiseTTL() time.Duration {
	return viper.GetDuration("p2p.discovery.advertise-ttl")
}

func P2PDiscoveryAdvertiseLimit() int {
	return viper.GetInt("p2p.discovery.advertise-limit")
}

func P2PDiscoveryAllowAll() bool {
	return viper.GetBool("p2p.discovery.allow-all")
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

func P2PDiscoveryRetryInterval() time.Duration {
	return viper.GetDuration("p2p.discovery.retry")
}

// String functions

func P2PPortString() string {
	return strconv.Itoa(P2PPort())
}

func P2PDiscoveryRetryIntervalString() string {
	return strconv.Itoa(int(P2PDiscoveryRetryInterval()))
}
