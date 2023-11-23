package peer

type Peer struct {
	// ID is the peer's ID
	id string
	// Name is the peer's name
	alias string
}

func New(id string, alias string) *Peer {
	return &Peer{
		id:    id,
		alias: alias,
	}
}

func NewFromID(id string) *Peer {
	return &Peer{
		id:    id,
		alias: id[len(id)-8:],
	}
}

func (p *Peer) Alias() string {
	return p.alias
}

func (p *Peer) ID() string {
	return p.id
}
