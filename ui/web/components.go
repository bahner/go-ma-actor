package web

import (
	"sort"

	"github.com/bahner/go-ma-actor/p2p/peer"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

func unorderedListFromPeerIDSlice(peers p2peer.IDSlice) string {
	peersMap := make(map[string]string)
	for _, p := range peers {
		nick := peer.GetOrCreateNick(p)
		peersMap[p.String()] = nick
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

func unorderedListFromStringSlice(s []string) string {

	list := "<table>\n"
	for _, v := range s {
		list += "<tr><td>" + v + "</td></tr>\n"
	}
	list += "</table>\n"
	return list
}
