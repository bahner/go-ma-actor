package web

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Relay struct {
	P2P *p2p.P2P
}
type RelayDocument struct {
	Title            string
	H1               string
	Addrs            []multiaddr.Multiaddr
	ProtectedPeers   p2peer.IDSlice
	UnprotectedPeers p2peer.IDSlice
}

// NewWebHandler creates a new Relay instance.
func NewRelayHandler(p *p2p.P2P) *Relay {
	return &Relay{
		P2P: p,
	}
}

// ServeHTTP implements the http.Handler interface for Relay.
// This allows Relay to be directly used as an HTTP handler.
func (data *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Delegate the actual handling to the WebHandler method.
	data.WebHandler(w, r)
}

func (d *RelayDocument) String() string {

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	html += "<style>table, th, td {border: 1px solid black;}</style>"

	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, config.HttpRefresh())
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"

	// Info leak? Not really important anyways.
	// // Addresses
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n"
		html += "<table>\n"
		for _, addr := range d.Addrs {
			html += "<tr><td>" + addr.String() + "</td></tr>"
		}
		html += "</table>"
	}

	// Peers with Same Rendezvous
	if len(d.ProtectedPeers) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n", len(d.ProtectedPeers))
		html += unorderedListFromPeerIDSlice(d.ProtectedPeers)
	}
	// All Connected Peers
	if len(d.UnprotectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n", len(d.UnprotectedPeers))

		html += unorderedListFromPeerIDSlice(d.UnprotectedPeers)
	}

	html += "</body>\n</html>"
	return html
}

func newRelayDocument() *RelayDocument {
	return &RelayDocument{}
}

func (data *Relay) WebHandler(w http.ResponseWriter, r *http.Request) {
	relayHandler(w, r, data.P2P)
}

func relayHandler(w http.ResponseWriter, _ *http.Request, p *p2p.P2P) {

	doc := newRelayDocument()

	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s.", ma.RENDEZVOUS)
	doc.H1 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Host.ID().String()))
	doc.H1 += fmt.Sprintf("<br>Found %d peers with rendezvous %s", len(p.ConnectedProtectedPeers()), ma.RENDEZVOUS)
	doc.Addrs = p.Host.Addrs()
	doc.ProtectedPeers = p.ConnectedProtectedPeers()
	doc.UnprotectedPeers = p.ConnectedUnprotectedPeers()
	// doc.AllConnectedPeers = p.GetAllConnectedPeers()

	fmt.Fprint(w, doc.String())
}
