package mode

import (
	"sort"

	"github.com/bahner/go-ma-actor/p2p/peer"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

func UnorderedListFromPeerIDSlice(peers p2peer.IDSlice) string {
	peersMap := make(map[string]string)
	for _, p := range peers {
		id := p.String()
		nick, err := peer.LookupNick(id)
		if err != nil {
			peersMap[id] = id
		} else {
			peersMap[id] = nick
		}
	}

	var keys []string
	for _, p := range peers {
		keys = append(keys, p.String())
	}
	sort.Strings(keys)

	list := "<ul>\n"
	for _, v := range keys {
		if peersMap[v] == v {
			list += "<li>" + v + "</li>\n"
		} else {
			list += "<li>" + v + "(" + peersMap[v] + ")</li>\n"
		}
	}
	list += "</ul>\n"
	return list
}
