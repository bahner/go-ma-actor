package config

import (
	"context"
	"fmt"
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
	defaultDesiredPeers  int = 3

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnMgrGrace           time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 1
)

func init() {

	// P2P Settings
	pflag.Int("low_watermark", defaultLowWaterMark, "Low watermark for peer discovery.")
	viper.BindPFlag("libp2p.connmgr.low_watermark", pflag.Lookup("low_watermark"))

	pflag.Int("high_watermark", defaultHighWaterMark, "High watermark for peer discovery.")
	viper.BindPFlag("libp2p.connmgr.high_watermark", pflag.Lookup("high_watermark"))

	pflag.Int("desired_peers", defaultDesiredPeers, "Desired number of peers to connect to.")
	viper.BindPFlag("libp2p.connmgr.desired_peers", pflag.Lookup("desired_peers"))

	pflag.Duration("grace_period", defaultConnMgrGrace, "Grace period for connection manager.")
	viper.BindPFlag("libp2p.connmgr.grace_period", pflag.Lookup("grace_period"))

	pflag.Duration("discovery_retry", defaultDiscoveryRetryInterval, "Retry interval for peer discovery.")
	viper.BindPFlag("libp2p.discovery_retry", pflag.Lookup("discovery_retryl"))

	pflag.Duration("discovery_timeout", defaultDiscoveryTimeout, "Timeout for peer discovery.")
	viper.BindPFlag("libp2p.connmgr.discovery_timeout", pflag.Lookup("discoveryTimeout"))
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

	_, privKeyBytes, err := mb.Decode(viper.GetString("libp2p.identity"))
	if err != nil {
		return nil
	}

	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil
	}

	return privKey

}

func generateNodeIdentity() (string, error) {
	pk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		log.Errorf("failed to generate node identity: %s", err)
		return "", err
	}

	pkBytes, err := pk.Raw()
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
func GetDiscoveryTimeout() time.Duration {
	return time.Duration(viper.GetDuration("discoveryTimeout")) * time.Second
}

func GetDiscoveryTimeoutString() string {
	return GetDiscoveryTimeout().String()
}

func GetLowWaterMark() int {
	return viper.GetInt("lowWaterMark")
}

func GetLowWatermarkString() string {
	return fmt.Sprint(GetLowWaterMark())
}

func GetHighWaterMark() int {
	return viper.GetInt("highWaterMark")
}

func GetHighWatermarkString() string {
	return fmt.Sprint(GetHighWaterMark())
}

func GetConnMgrGracePeriod() time.Duration {
	return viper.GetDuration("connmgrGracePeriod")
}

func GetConnMgrGraceString() string {
	return GetConnMgrGracePeriod().String()
}

func GetDiscoveryContext() (context.Context, func()) {

	ctx := context.Background()

	discoveryCtx, cancel := context.WithTimeout(ctx, GetDiscoveryTimeout())

	return discoveryCtx, cancel
}

func GetDiscoveryRetryInterval() time.Duration {
	return viper.GetDuration("discoveryRetryInterval")
}

func GetDiscoveryRetryIntervalString() string {
	return GetDiscoveryRetryInterval().String()
}

func GetDesiredPeers() int {
	return viper.GetInt("desiredPeers")
}

func GetDesiredPeersString() string {
	return fmt.Sprint(GetDesiredPeers())
}
