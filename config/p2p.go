package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	mb "github.com/multiformats/go-multibase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultLowWaterMark  int = 3
	defaultHighWaterMark int = 10
	defaultListenPort    int = 0

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnMgrGrace           time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 60
)

func init() {

	// P2P Settings
	pflag.Int("low_watermark", defaultLowWaterMark, "Low watermark for peer discovery.")
	viper.SetDefault("libp2p.connmgr.low_watermark", defaultLowWaterMark)
	viper.BindPFlag("libp2p.connmgr.low_watermark", pflag.Lookup("low_watermark"))

	pflag.Int("high_watermark", defaultHighWaterMark, "High watermark for peer discovery.")
	viper.SetDefault("libp2p.connmgr.high_watermark", defaultHighWaterMark)
	viper.BindPFlag("libp2p.connmgr.high_watermark", pflag.Lookup("high_watermark"))

	// pflag.Int("desired_peers", defaultDesiredPeers, "Desired number of peers to connect to.")
	// viper.SetDefault("libp2p.connmgr.desired_peers", defaultDesiredPeers)
	// viper.BindPFlag("libp2p.connmgr.desired_peers", pflag.Lookup("desired_peers"))

	pflag.Duration("grace_period", defaultConnMgrGrace, "Grace period for connection manager.")
	viper.SetDefault("libp2p.connmgr.grace_period", defaultConnMgrGrace)
	viper.BindPFlag("libp2p.connmgr.grace_period", pflag.Lookup("grace_period"))

	pflag.Duration("discovery_retry", defaultDiscoveryRetryInterval, "Retry interval for peer discovery.")
	viper.SetDefault("libp2p.discovery_retry", defaultDiscoveryRetryInterval)
	viper.BindPFlag("libp2p.discovery_retry", pflag.Lookup("discovery_retryl"))

	pflag.Duration("discovery_timeout", defaultDiscoveryTimeout, "Timeout for peer discovery.")
	viper.SetDefault("libp2p.discovery_timeout", defaultDiscoveryTimeout)
	viper.BindPFlag("libp2p.connmgr.discovery_timeout", pflag.Lookup("discoveryTimeout"))

	pflag.Int("listen_port", defaultListenPort, "Port for libp2p node to listen on.")
	viper.SetDefault("libp2p.port", defaultListenPort)
	viper.BindPFlag("libp2p.port", pflag.Lookup("listen_port"))
}

// P2P Node identity
func InitP2P() {

	var i string

	if viper.GetString("identity") == "" {
		i, _ = generateNodeIdentity()

		viper.Set("identity", i)
	}

	log.Debugf("Node identity: %s", i)
}

func GetNodeMultibasePrivKey() string {

	return viper.GetString("identity")
}

func GetNodeIdentity() crypto.PrivKey {

	log.Debugf("config.GetNodeIdentity: %s", viper.GetString("libp2p.identity"))
	_, privKeyBytes, err := mb.Decode(viper.GetString("libp2p.identity"))
	if err != nil {
		log.Debugf("config.GetNodeIdentity: Failed to decode node identity: %v", err)
		return nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		log.Debugf("config.GetNodeIdentity: Failed to unmarshal node identity: %v", err)
		return nil
	}

	log.Debug("Config.GetNodeIdentity: ", privKey.GetPublic())
	return privKey

}

func generateAndPrintNodeIdentity() error {

	p2pPrivKey, err := generateNodeIdentity()
	if err != nil {
		return fmt.Errorf("config.initIdentity: Failed to generate node identity: %v", err)
	}

	fmt.Println(ENV_PREFIX + "_LIBP2P_IDENTITY=" + p2pPrivKey)

	return nil
}

func generateNodeIdentity() (string, error) {
	pk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	pkBytes, err := crypto.MarshalPrivateKey(pk)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	ni, err := mb.Encode(mb.Base58BTC, pkBytes)
	if err != nil {
		log.Errorf("failed to encode node identity: %s", err)
		return "", err
	}

	return ni, nil

}

// P2P Settings

func GetDiscoveryContext() (context.Context, func()) {

	ctx := context.Background()

	discoveryCtx, cancel := context.WithTimeout(ctx, GetDiscoveryTimeout())

	return discoveryCtx, cancel
}

func GetDiscoveryTimeout() time.Duration {
	return time.Duration(viper.GetDuration("libp2p.discovery_timeout")) * time.Second
}

func GetDiscoveryTimeoutString() string {
	return GetDiscoveryTimeout().String()
}

func GetLowWaterMark() int {
	return viper.GetInt("libp2p.connmgr.low_watermark")
}

func GetLowWatermarkString() string {
	return fmt.Sprint(GetLowWaterMark())
}

func GetHighWaterMark() int {
	return viper.GetInt("libp2p.connmgr.high_watermark")
}

func GetHighWatermarkString() string {
	return fmt.Sprint(GetHighWaterMark())
}

func GetConnMgrGracePeriod() time.Duration {
	return viper.GetDuration("libp2p.connmgr.grace_period")
}

func GetConnMgrGraceString() string {
	return GetConnMgrGracePeriod().String()
}

func GetListenPort() int {
	return viper.GetInt("libp2p.port")
}

func GetListenPortString() string {
	return strconv.Itoa(GetListenPort())
}

func GetDiscoveryRetryInterval() time.Duration {
	return viper.GetDuration("libp2p.discovery_retry")
}

func GetDiscoveryRetryIntervalString() string {
	return GetDiscoveryRetryInterval().String()
}
