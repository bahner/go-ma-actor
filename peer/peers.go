package peer

var peers map[string]*Peer

func init() {
	peers = make(map[string]*Peer)
}

// Add adds a peer to the map
func Add(p *Peer) {
	peers[p.ID()] = p
}

// Get returns a peer from the map
func Get(id string) *Peer {
	return peers[id]
}

// List returns a list of peers
func List() []*Peer {
	var result []*Peer
	for _, p := range peers {
		result = append(result, p)
	}
	return result
}

// Remove removes a peer from the map
func Delete(id string) {
	delete(peers, id)
}

func ListAliases() []string {
	var result []string
	for _, p := range peers {
		result = append(result, p.Alias())
	}
	return result
}

func ListIDs() []string {
	var result []string
	for _, p := range peers {
		result = append(result, p.ID())
	}
	return result
}
