package peer

type Peer struct {
	// ID is the peer's ID
	ID string
	// Name is the peer's name
	Alias string
}

func New(id string, alias string) *Peer {
	return &Peer{
		ID:    id,
		Alias: alias,
	}
}

func NewFromID(id string) *Peer {
	return New(id, id[len(id)-8:])
}

func GetOrCreate(id string) *Peer {
	p := Get(id)
	if p == nil {
		p = NewFromID(id)
		Add(p)
	}
	return p
}
