package relay

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

type WebHandlerData struct {
	P2P *p2p.P2P
}

type WebHandlerDocument struct {
	Title            string
	H1               string
	Addrs            []multiaddr.Multiaddr
	ProtectedPeers   peer.IDSlice
	UnprotectedPeers peer.IDSlice
}

func NewWebHandlerDocument() *WebHandlerDocument {
	return &WebHandlerDocument{}
}

func (data *WebHandlerData) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, data.P2P)
}

func webHandler(w http.ResponseWriter, r *http.Request, p *p2p.P2P) {

	doc := NewWebHandlerDocument()

	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s.", ma.RENDEZVOUS)
	doc.H1 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Node.ID().String()))
	doc.H1 += fmt.Sprintf("<br>Found %d peers with rendezvous %s", len(p.GetConnectedProtectedPeers()), ma.RENDEZVOUS)
	doc.Addrs = p.Node.Addrs()
	doc.ProtectedPeers = p.GetConnectedProtectedPeers()
	doc.UnprotectedPeers = p.GetConnectedUnprotectedPeers()
	// doc.AllConnectedPeers = p.GetAllConnectedPeers()

	fmt.Fprint(w, doc.String())
}

func (d *WebHandlerDocument) String() string {

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
	if len(d.ProtectedPeers) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n<ul>", len(d.ProtectedPeers))
		for _, peer := range d.ProtectedPeers {
			html += "<li>" + peer.String() + "(" + peer.ShortString() + ")</li>"
		}
		html += "</ul>"
	}
	// All Connected Peers
	if len(d.UnprotectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n<ul>", len(d.UnprotectedPeers))
		for _, peer := range d.UnprotectedPeers {
			html += "<li>" + peer.String() + "</li>"
		}
		html += "</ul>"
	}

	html += "</body>\n</html>"
	return html
}
