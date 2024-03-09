package relay

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma-actor/p2p/peer"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
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
	ProtectedPeers   map[string]string
	UnprotectedPeers map[string]string
}

func NewWebHandlerDocument() *WebHandlerDocument {
	return &WebHandlerDocument{}
}

func (data *WebHandlerData) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, data.P2P)
}

func webHandler(w http.ResponseWriter, _ *http.Request, p *p2p.P2P) {

	doc := NewWebHandlerDocument()

	doc.Title = fmt.Sprintf("Bootstrap peer for rendezvous %s.", ma.RENDEZVOUS)
	doc.H1 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Host.ID().String()))
	doc.H1 += fmt.Sprintf("<br>Found %d peers with rendezvous %s", len(p.ConnectedProtectedPeers()), ma.RENDEZVOUS)
	doc.Addrs = p.Host.Addrs()
	doc.ProtectedPeers = createSortedMapOfPeersNicks(p.ConnectedProtectedPeers())
	doc.UnprotectedPeers = createSortedMapOfPeersNicks(p.ConnectedUnprotectedPeers())
	// doc.AllConnectedPeers = p.GetAllConnectedPeers()

	fmt.Fprint(w, doc.String())
}

// Take peerIDslice and convert them to a map of map[id] = nick
func createSortedMapOfPeersNicks(peers p2peer.IDSlice) map[string]string {
	unsortedPeersMap := make(map[string]string)
	for _, p := range peers {
		id := p.String()
		unsortedPeersMap[id] = peer.Lookup(id)
	}

	var keys []string
	for _, p := range peers {
		keys = append(keys, p.String())
	}
	sort.Strings(keys)

	peersMap := make(map[string]string)
	for _, k := range keys {
		peersMap[k] = unsortedPeersMap[k]
	}
	return peersMap
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
		html += createListFromMap(d.ProtectedPeers)
		html += "</ul>"
	}
	// All Connected Peers
	if len(d.UnprotectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n<ul>", len(d.UnprotectedPeers))
		html += createListFromMap(d.UnprotectedPeers)
		html += "</ul>"
	}

	html += "</body>\n</html>"
	return html
}

func createListFromMap(m map[string]string) string {
	list := ""
	for k, v := range m {
		if k == v {
			list += "<li>" + k + "</li>"
		} else {
			list += "<li>" + k + "(" + v + ")</li>"
		}
	}
	return list
}
