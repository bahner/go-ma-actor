package config

import (
	"strconv"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConnmgrLowWatermark  int           = 512
	defaultConnmgrHighWatermark int           = 768
	defaultConnmgrGracePeriod   time.Duration = time.Minute * 2

	defaultDiscoveryAdvertiseInterval time.Duration = time.Minute * 5
	defaultDiscoveryAdvertiseTTL      time.Duration = time.Minute * 60
	defaultDiscoveryAdvertiseLimit    int           = 100

	ALLOW_ALL_PEERS bool = true // Allow all peers by default. This is the norm for now. Use connmgr threshold and protection instead.

	defaultListenPort int    = 0
	fakeP2PIdentity   string = "NO_DEFAULT_NODE_IDENITY"
	defaultDHT        bool   = true
	defaultMDNS       bool   = true
)

func init() {

	pflag.Bool("dht", defaultDHT, "Whether to discover using DHT")
	pflag.Bool("mdns", defaultMDNS, "Whether to discover using MDNS")
	pflag.Duration("connmgr-grace-period", defaultConnmgrGracePeriod, "Grace period for connection manager.")
	pflag.Duration("discovery-advertise-interval", defaultDiscoveryAdvertiseInterval, "How often to advertise our presence to libp2p")
	pflag.Duration("discovery-advertise-ttl", defaultDiscoveryAdvertiseTTL, "Hint of TimeToLive for advertising peer discovery.")
	pflag.Int("connmgr-high-watermark", defaultConnmgrHighWatermark, "High watermark for peer discovery.")
	pflag.Int("connmgr-low-watermark", defaultConnmgrLowWatermark, "Low watermark for peer discovery.")
	pflag.Int("discovery-advertise-limit", defaultDiscoveryAdvertiseLimit, "Limit for advertising peer discovery.")
	pflag.Int("port", defaultListenPort, "Port for libp2p node to listen on.")

	// Bind pflags
	viper.BindPFlag("p2p.connmgr.grace-period", pflag.Lookup("connmgr-grace-period"))
	viper.BindPFlag("p2p.connmgr.high-watermark", pflag.Lookup("connmgr-high-watermark"))
	viper.BindPFlag("p2p.connmgr.low-watermark", pflag.Lookup("connmgr-low-watermark"))
	viper.BindPFlag("p2p.discovery.advertise-interval", pflag.Lookup("discovery-advertise-interval"))
	viper.BindPFlag("p2p.discovery.advertise-limit", pflag.Lookup("discovery-advertise-limit"))
	viper.BindPFlag("p2p.discovery.advertise-ttl", pflag.Lookup("discovery-advertise-ttl"))
	viper.BindPFlag("p2p.discovery.dht", pflag.Lookup("dht"))
	viper.BindPFlag("p2p.discovery.mdns", pflag.Lookup("mdns"))
	viper.BindPFlag("p2p.port", pflag.Lookup("port"))
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
	Identity  string          `yaml:"identity"`
	Port      int             `yaml:"port"`
	Connmgr   ConnmgrStruct   `yaml:"connmgr"`
	Discovery DiscoveryStruct `yaml:"discovery"`
}

func P2P() P2PConfig {
	viper.SetDefault("p2p.identity", fakeP2PIdentity)

	p2pIdentity, err := P2PIdentity()
	if err != nil {
		panic(err)
	}

	return P2PConfig{
		Identity: p2pIdentity,
		Port:     P2PPort(),
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

func P2PIdentity() (string, error) {

	if GenerateFlag() {
		return generateNodeIdentity()
	}

	return viper.GetString("p2p.identity"), nil
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

func P2PPort() int {
	return viper.GetInt("p2p.port")
}

// String functions

func P2PPortString() string {
	return strconv.Itoa(P2PPort())
}
