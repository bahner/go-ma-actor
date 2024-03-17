package main

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

	list := "<table>\n"
	for _, v := range keys {
		if peersMap[v] == v {
			list += "<tr><td span=2>" + v + "</td></tr>\n"
		} else {
			list += "<tr><td>" + v + "</td><td>" + peersMap[v] + "</td></tr>\n"
		}
	}
	list += "</table>\n"
	return list
}
