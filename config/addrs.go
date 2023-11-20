package config

func GetListenAddrStrings(port string) []string {

	return []string{
		"/ip4/0.0.0.0/tcp/" + port,
		"/ip4/0.0.0.0/udp/" + port + "/quic",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1",
		"/ip4/0.0.0.0/udp/" + port + "/quic-v1/webtransport",
		"/ip6/::/tcp/" + port,
		"/ip6/::/udp/" + port + "/quic",
		"/ip6/::/udp/" + port + "/quic-v1",
		"/ip6/::/udp/" + port + "/quic-v1/webtransport",
	}
}
