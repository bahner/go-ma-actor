package peer

import "github.com/bahner/go-ma-actor/config/db"

const _SELECT_IDS = "SELECT id FROM peers"

// Returns a slic of all known peer IDs.
func IDS() ([]string, error) {
	db, err := db.Get()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(_SELECT_IDS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
func Peers() []Peer {

	peers := []Peer{}
	ids, err := IDS()
	if err != nil {
		return nil
	}
	for _, id := range ids {
		p, err := Get(id)
		if err != nil {
			return nil
		}
		peers = append(peers, p)
	}

	return peers
}
