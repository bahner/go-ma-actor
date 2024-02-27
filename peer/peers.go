package peer

import (
	"errors"

	"github.com/bahner/go-ma-actor/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_SELECT = "SELECT node FROM nodes WHERE id =?"
	_UPSERT = "INSERT INTO nodes (id, node) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET node = ?"
	_DELETE = "DELETE FROM nodes WHERE id = ?"
	_PEERS  = "SELECT node FROM nodes"
)

const (
	defaultAliasLength = 8
)

var (
	ErrNotFound = errors.New("Peer not found")
)

// Sets a node in the database
// The key is the node's ID
func Set(p Peer) error {

	d, err := db.Get()
	if err != nil {
		return err
	}

	v, err := p.MarshalToCBOR()
	if err != nil {
		return err
	}
	_, err = d.Exec(_UPSERT, p.ID, v, v)
	if err != nil {
		return err
	}

	return nil
}

// Fetches a nodefrom the database
// Returns a Peer{} if it does not exist
func Get(id string) (Peer, error) {

	db, err := db.Get()
	if err != nil {
		return Peer{}, err
	}

	var (
		p Peer
		v []byte
	)

	err = db.QueryRow(_SELECT, id).Scan(&v)
	if err != nil {
		return Peer{}, err
	}

	err = UnmarshalFromCBOR(v, &p)
	if err == nil {
		return p, nil
	}

	return Peer{}, err
}

// Removes a node from the database if it exists
func Remove(id string) error {

	d, err := db.Get()
	if err != nil {
		return err
	}

	_, err = d.Exec(_DELETE, id)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

// Returns the ID of a node. Returns the input if the node does not exist
func GetID(id string) string {

	p, err := LookupNick(id)
	if err == nil {
		return p.ID
	}

	return id

}

func LookupNick(id string) (Peer, error) {

	peers, err := Peers()
	if err != nil {
		return Peer{}, err
	}

	for _, p := range peers {
		if p.Nick == id {
			return p, nil
		}
	}

	return Peer{}, ErrNotFound
}

func Lookup(id string) (Peer, error) {

	// Lookup Nick first, that's probably the most used case
	p, err := LookupNick(id)
	if err == nil {
		return p, nil
	}

	return Get(id)
}

// Peers returns a slice of all peers in the database
// If an error occurs, an empty slice is returned
func Peers() ([]Peer, error) {
	var (
		peers = []Peer{}
		p     Peer
		b     []byte
	)

	d, err := db.Get()
	if err != nil {
		return nil, err
	}

	rows, err := d.Query(_PEERS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}
		p = Peer{}
		err = UnmarshalFromCBOR(b, &p)
		if err != nil || p.ID == "" || p.Nick == "" || p.AddrInfo == nil {
			continue
		}

		peers = append(peers, p)
	}

	return peers, nil
}
