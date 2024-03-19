package web

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/peer"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type WebEntity struct {
	P2P    *p2p.P2P
	Entity *entity.Entity
}

type WebEntityDocument struct {
	Title               string
	H1                  string
	H2                  string
	Addrs               []multiaddr.Multiaddr
	PeersWithSameRendez p2peer.IDSlice
	AllConnectedPeers   p2peer.IDSlice
}

// NewWebHandler creates a new WebEntity instance.
func NewWebEntityHandler(p *p2p.P2P, e *entity.Entity) *WebEntity {
	return &WebEntity{
		P2P:    p,
		Entity: e,
	}
}

// ServeHTTP implements the http.Handler interface for WebEntity.
// This allows WebEntity to be directly used as an HTTP handler.
func (data *WebEntity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Delegate the actual handling to the WebHandler method.
	data.WebHandler(w, r)
}

func (d *WebEntityDocument) String() string {

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
	if d.H2 != "" {
		html += "<h2>" + d.H2 + "</h2>\n"
	}

	// Info leak? Not really important anyways.
	// // Addresses
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n"
		html += "<table>\n"
		for _, addr := range d.Addrs {
			html += "<tr><td>" + addr.String() + "</td></tr>\n"
		}
		html += "</table>\n"
	}

	// Peers with Same Rendezvous
	if len(d.PeersWithSameRendez) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n", len(d.PeersWithSameRendez))
		html += d.unorderedListFromPeerIDSlice(d.PeersWithSameRendez)
	}
	// All Connected Peers
	if len(d.AllConnectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n", len(d.AllConnectedPeers))
		html += d.unorderedListFromPeerIDSlice(d.AllConnectedPeers)
	}

	html += "</body>\n</html>"
	return html
}

func newWebEntityDocument() *WebEntityDocument {
	return &WebEntityDocument{}
}

func (data *WebEntity) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, data.P2P, data.Entity)
}

func webHandler(w http.ResponseWriter, _ *http.Request, p *p2p.P2P, e *entity.Entity) {

	doc := newWebEntityDocument()

	titleStr := fmt.Sprintf("Entity: %s", e.DID.Id)
	h1str := titleStr
	doc.Title = titleStr
	doc.H1 = h1str
	doc.H2 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Host.ID().String()))
	doc.Addrs = p.Host.Addrs()
	doc.AllConnectedPeers = p.AllConnectedPeers()
	doc.PeersWithSameRendez = p.ConnectedProtectedPeers()

	fmt.Fprint(w, doc.String())
}

func (d *WebEntityDocument) unorderedListFromPeerIDSlice(peers p2peer.IDSlice) string {
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
