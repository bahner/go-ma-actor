package entity

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/mode"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

type WebHandlerData struct {
	P2P    *p2p.P2P
	Entity *entity.Entity
}

type WebHandlerDocument struct {
	Title               string
	H1                  string
	H2                  string
	Addrs               []multiaddr.Multiaddr
	PeersWithSameRendez peer.IDSlice
	AllConnectedPeers   peer.IDSlice
}

func NewWebHandlerDocument() *WebHandlerDocument {
	return &WebHandlerDocument{}
}

func (data *WebHandlerData) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, data.P2P, data.Entity)
}

func webHandler(w http.ResponseWriter, _ *http.Request, p *p2p.P2P, e *entity.Entity) {

	doc := NewWebHandlerDocument()

	titleStr := fmt.Sprintf("Entity: %s", e.DID.Id)
	h1str := titleStr

	if config.PongMode() {
		h1str = fmt.Sprintf("%s (Pong mode)", titleStr)
	}

	if config.PongFortuneMode() {
		h1str = fmt.Sprintf("%s (Pong mode with fortune cookies)", titleStr)
	}

	doc.Title = titleStr
	doc.H1 = h1str
	doc.H2 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Host.ID().String()))
	doc.Addrs = p.Host.Addrs()
	doc.AllConnectedPeers = p.AllConnectedPeers()
	doc.PeersWithSameRendez = p.ConnectedProtectedPeers()

	fmt.Fprint(w, doc.String())
}

func (d *WebHandlerDocument) String() string {

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, config.HttpRefresh())
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
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n", len(d.PeersWithSameRendez))
		html += mode.UnorderedListFromPeerIDSlice(d.PeersWithSameRendez)
	}
	// All Connected Peers
	if len(d.AllConnectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n", len(d.AllConnectedPeers))
		html += mode.UnorderedListFromPeerIDSlice(d.AllConnectedPeers)
	}

	html += "</body>\n</html>"
	return html
}
