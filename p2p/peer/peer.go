package peer

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/config/db"
	"github.com/bahner/go-ma-actor/internal"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	_DELETE_PEER = "DELETE FROM peers WHERE id = ?"
	_UPSERT_PEER = "INSERT INTO peers (id, nick, allowed) VALUES (?, ?, ?) ON CONFLICT(id) DO UPDATE SET nick = excluded.nick, allowed = excluded.allowed;"
	_SELECT_PEER = "SELECT id, nick, allowed FROM peers WHERE id = ?"
)

type Peer struct {
	Id         string
	Nick       string
	AllowedInt int
}

func New(id string, nick string, allowed bool) Peer {
	return Peer{
		Id:         id,
		Nick:       nick,
		AllowedInt: internal.Bool2int(allowed),
	}
}

func (p Peer) Allowed() bool {
	return internal.Int2bool(p.AllowedInt)
}

func (p *Peer) SetAllowed(allowed bool) error {
	p.AllowedInt = internal.Bool2int(allowed)
	return p.Commit()
}

func (p *Peer) SetNick(nick string) error {
	p.Nick = nick
	return p.Commit()
}

// Upsert modifies an existing peer's information in the map and the database.
func (p Peer) Commit() error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(_UPSERT_PEER, p.Id, p.Nick, p.AllowedInt)
	if err != nil {
		log.Debugf("Failed to set peer %s: %s", p.Id, err.Error())
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Delete removes a peer from the map and the database by ID.
func Delete(id string) error {
	d, err := db.Get()
	if err != nil {
		return err
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(_DELETE_PEER, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Get a peer by ID.
func Get(id string) (p Peer, err error) {
	d, err := db.Get()
	if err != nil {
		return Peer{}, err
	}

	err = d.QueryRow(_SELECT_PEER, id).Scan(&p.Id, &p.Nick, &p.AllowedInt)
	if err != nil {
		return Peer{}, err
	}

	return p, nil
}

// Takes a quiery of a nick or an ID and returns the ID.
func GetOrCreate(id string) (p Peer, err error) {

	// Id the ID exists, return it
	p, err = Get(id)
	if err == nil {
		return p, nil
	}

	return New(id, GetOrCreateNick(id), config.ALLOW_ALL_PEERS), nil
}
