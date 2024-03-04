package config

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConnmgrLowWatermark  int = 100
	defaultConnmgrHighWatermark int = 300
	defaultListenPort           int = 0

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnmgrGracePeriod     time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 60

	fakeP2PIdentity string = "NO_DEFAULT_NODE_IDENITY"
)

func init() {

	// P2P Settings
	pflag.Int("p2p-connmgr-low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")
	pflag.Int("p2p-connmgr-high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
	pflag.Duration("p2p-connmgr-grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
	pflag.Duration("p2p-discovery-retry", defaultDiscoveryRetryInterval, "Retry interval for peer discovery.")
	pflag.Duration("p2p-discovery-timeout", defaultDiscoveryTimeout, "Timeout for peer discovery.")
	pflag.Int("p2p-port", defaultListenPort, "Port for libp2p node to listen on.")
}

// P2P Node identity
func InitP2P() {
	viper.BindPFlag("p2p.connmgr.low-watermark", pflag.Lookup("p2p-connmgr-low-watermark"))
	viper.BindPFlag("p2p.connmgr.high-watermark", pflag.Lookup("p2p-connmgr-high-watermark"))
	viper.BindPFlag("p2p.connmgr.grace-period", pflag.Lookup("p2p-connmgr-grace-period"))
	viper.BindPFlag("p2p.discovery-retry", pflag.Lookup("p2p-discovery-retryl"))
	viper.BindPFlag("p2p.discovery-timeout", pflag.Lookup("p2p-discoveryTimeout"))
	viper.BindPFlag("p2p.port", pflag.Lookup("p2p-port"))

	viper.SetDefault("p2p.connmgr.low-watermark", defaultConnmgrLowWatermark)
	viper.SetDefault("p2p.connmgr.high-watermark", defaultConnmgrHighWatermark)
	viper.SetDefault("p2p.connmgr.grace-period", defaultConnmgrGracePeriod)
	viper.SetDefault("p2p.discovery-retry", defaultDiscoveryRetryInterval)
	viper.SetDefault("p2p.discovery-timeout", defaultDiscoveryTimeout)
	viper.SetDefault("p2p.port", defaultListenPort)
	viper.SetDefault("p2p.identity", fakeP2PIdentity)

	var i string

	if P2PIdentity() == fakeP2PIdentity {
		i, _ = generateNodeIdentity()

		viper.Set("p2p.identity", i)
	}

	log.Debugf("P2P identity: %s", i)
}

// P2P Settings

func P2PIdentity() string {

	return viper.GetString("p2p.identity")
}

func P2PDiscoveryContext() (context.Context, func()) {

	ctx := context.Background()

	discoveryCtx, cancel := context.WithTimeout(ctx, P2PDiscoveryTimeout())

	return discoveryCtx, cancel
}

func P2PDiscoveryTimeout() time.Duration {
	return time.Duration(viper.GetDuration("p2p.discovery-timeout"))
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
	return viper.GetDuration("p2p.discovery-retry")
}

// String functions

func P2PPortString() string {
	return strconv.Itoa(P2PPort())
}

func P2PDiscoveryRetryIntervalString() string {
	return strconv.Itoa(int(P2PDiscoveryRetryInterval()))
}
