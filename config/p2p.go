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
	defaultConnmgrLowWatermark  int = 10
	defaultConnmgrHighWatermark int = 30
	defaultListenPort           int = 0

	defaultDiscoveryTimeout       time.Duration = time.Second * 30
	defaultConnmgrGracePeriod     time.Duration = time.Minute * 1
	defaultDiscoveryRetryInterval time.Duration = time.Second * 10
)

func init() {

	// P2P Settings
	pflag.Int("low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")
	viper.SetDefault("libp2p.connmgr.low-watermark", defaultConnmgrLowWatermark)
	viper.BindPFlag("libp2p.connmgr.low-watermark", pflag.Lookup("low-watermark"))

	pflag.Int("high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
	viper.SetDefault("libp2p.connmgr.high-watermark", defaultConnmgrHighWatermark)
	viper.BindPFlag("libp2p.connmgr.high-watermark", pflag.Lookup("high-watermark"))

	pflag.Duration("grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
	viper.SetDefault("libp2p.connmgr.grace-period", defaultConnmgrGracePeriod)
	viper.BindPFlag("libp2p.connmgr.grace-period", pflag.Lookup("grace-period"))

	pflag.Duration("discovery-retry", defaultDiscoveryRetryInterval, "Retry interval for peer discovery.")
	viper.SetDefault("libp2p.discovery-retry", defaultDiscoveryRetryInterval)
	viper.BindPFlag("libp2p.discovery-retry", pflag.Lookup("discovery-retryl"))

	pflag.Duration("discovery-timeout", defaultDiscoveryTimeout, "Timeout for peer discovery.")
	viper.SetDefault("libp2p.discovery-timeout", defaultDiscoveryTimeout)
	viper.BindPFlag("libp2p.discovery-timeout", pflag.Lookup("discoveryTimeout"))

	pflag.Int("listen-port", defaultListenPort, "Port for libp2p node to listen on.")
	viper.SetDefault("libp2p.port", defaultListenPort)
	viper.BindPFlag("libp2p.port", pflag.Lookup("listen-port"))
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
	return time.Duration(viper.GetDuration("libp2p.discovery-timeout"))
}

func GetDiscoveryTimeoutString() string {
	return GetDiscoveryTimeout().String()
}

func GetLowWatermark() int {
	return viper.GetInt("libp2p.connmgr.low-watermark")
}

func GetLowWatermarkString() string {
	return fmt.Sprint(GetLowWatermark())
}

func GetHighWatermark() int {
	return viper.GetInt("libp2p.connmgr.high-watermark")
}

func GetHighWatermarkString() string {
	return fmt.Sprint(GetHighWatermark())
}

func GetConnMgrGracePeriod() time.Duration {
	return viper.GetDuration("libp2p.connmgr.grace-period")
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
	return viper.GetDuration("libp2p.discovery-retry")
}

func GetDiscoveryRetryIntervalString() string {
	return GetDiscoveryRetryInterval().String()
}
