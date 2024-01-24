package main

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

func webHandler(w http.ResponseWriter, r *http.Request) {

	desiredPeers := config.GetDesiredPeers()

	allConnected := p.GetAllConnectedPeers()
	if allConnected == nil {
		log.Error("Failed to get connected peers.")
		allConnected = peer.IDSlice{}
	}
	peersWithRendezvous := p.GetConnectedProtectedPeers()
	if peersWithRendezvous == nil {
		log.Error("Failed to get connected peers with rendezvous.")
		peersWithRendezvous = peer.IDSlice{}
	}

	doc := New()
	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s. Version %s ", ma.RENDEZVOUS, VERSION)
	doc.H1 = fmt.Sprintf("%s@%s v%s", ma.RENDEZVOUS, (p.Node.ID().String()), VERSION)
	doc.H1 += fmt.Sprintf("<br>Found %d/%d peers with rendezvous %s", len(peersWithRendezvous), desiredPeers, ma.RENDEZVOUS)
	doc.Addrs = p.Node.Addrs()
	if allConnected == nil {
		allConnected = peer.IDSlice{}
	}
	doc.AllConnectedPeers = allConnected
	if peersWithRendezvous == nil {
		peersWithRendezvous = peer.IDSlice{}
	}
	doc.MaPeers = peersWithRendezvous
	doc.AllPeers = allConnected

	fmt.Fprint(w, doc.String())
}

type Document struct {
	Title             string
	H1                string
	Addrs             []multiaddr.Multiaddr
	MaPeers           peer.IDSlice
	AllConnectedPeers peer.IDSlice
	AllPeers          peer.IDSlice
}

func New() *Document {
	return &Document{}
}

func (d *Document) String() string {

	retryInterval := config.GetDiscoveryRetryInterval()

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, int(retryInterval.Seconds()))
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"

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
	if len(d.MaPeers) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n<ul>", len(d.MaPeers))
		for _, peer := range d.MaPeers {
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
