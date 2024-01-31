package main

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

type WebHandlerData struct {
	Host   host.Host
	Entity *entity.Entity
}

func (h *WebHandlerData) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, h.Host, h.Entity)
}

func webHandler(w http.ResponseWriter, r *http.Request, n host.Host, e *entity.Entity) {
	allConnected := getLivePeerIDs(n)

	doc := New()
	doc.Title = fmt.Sprintf("Pong peer: %s", e.DID.String())
	doc.H1 = doc.Title
	doc.H2 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (n.ID().String()))
	doc.Addrs = n.Addrs()
	doc.AllConnectedPeers = allConnected

	fmt.Fprint(w, doc.String())
}

func getLivePeerIDs(h host.Host) peer.IDSlice {
	peerSet := make(map[peer.ID]struct{})

	for _, conn := range h.Network().Conns() {
		peerSet[conn.RemotePeer()] = struct{}{}
	}

	peers := make(peer.IDSlice, 0, len(peerSet))
	for p := range peerSet {
		peers = append(peers, p)
	}

	return peers
}

type Document struct {
	Title               string
	H1                  string
	H2                  string
	Addrs               []multiaddr.Multiaddr
	PeersWithSameRendez peer.IDSlice
	AllConnectedPeers   peer.IDSlice
	OtherPeers          peer.IDSlice
}

func New() *Document {
	return &Document{}
}

func (d *Document) String() string {

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, int(config.GetDiscoveryTimeout().Seconds()))
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"
	if d.H2 != "" {
		html += "<h2>" + d.H2 + "</h2>\n"
	}

	// Info leak? Not really important anyways.
	// // Addresses
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n<ul>"
		for _, addr := range d.Addrs {
			html += "<li>" + addr.String() + "</li>"
		}
		html += "</ul>"
	}

	// Peers with Same Rendezvous
	if len(d.PeersWithSameRendez) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n<ul>", len(d.PeersWithSameRendez))
		for _, peer := range d.PeersWithSameRendez {
			html += "<li>" + peer.String() + "</li>"
		}
		html += "</ul>"
	}
	// All Connected Peers
	if len(d.AllConnectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n<ul>", len(d.AllConnectedPeers))
		for _, peer := range d.AllConnectedPeers {
			html += "<li>" + peer.String() + "</li>"
		}
		html += "</ul>"
	}

	html += "</body>\n</html>"
	return html
}
