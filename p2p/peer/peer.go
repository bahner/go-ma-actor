package peer

import (
	"fmt"

	"github.com/bahner/go-ma-actor/config"
	"github.com/fxamacker/cbor/v2"
	p2peer "github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	// ID is the peer's ID
	ID string
	// Name is the peer's name
	Nick string
	// AddrInfo
	AddrInfo *p2peer.AddrInfo
	// Allowed
	Allowed bool
}

// Create a new aliased addrinfo peer
// Always creates an alias of the last 8 characters of the ID initially
func New(addrInfo *p2peer.AddrInfo) Peer {

	addrInfo.MarshalJSON()
	id := addrInfo.ID.String()
	return Peer{
		ID:       id,
		Nick:     createNodeAlias(id),
		AddrInfo: addrInfo,
		Allowed:  config.P2PDiscoveryAllowAll(),
	}
}

func GetOrCreate(addrInfo *p2peer.AddrInfo) (Peer, error) {

	id := addrInfo.ID.String()

	p, err := Get(id)
	if err == nil {
		return p, nil
	}

	p = New(addrInfo)
	err = Set(p)
	if err != nil {
		return Peer{}, err
	}

	return p, nil
}

func createNodeAlias(id string) string {

	if len(id) <= defaultAliasLength {
		return id
	}

	return id[len(id)-defaultAliasLength:]

}

// Marshal returns the CBOR encoding of Peer, converting AddrInfo to JSON first.
func (p *Peer) MarshalToCBOR() ([]byte, error) {
	// Marshal AddrInfo to JSON
	addrInfoJSON, err := p.AddrInfo.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AddrInfo to JSON: %w", err)
	}

	// Create a map to represent the Peer struct including AddrInfo as a JSON string
	data := map[string]interface{}{
		"ID":       p.ID,
		"Nick":     p.Nick,
		"AddrInfo": string(addrInfoJSON), // Store the JSON string
	}

	// Marshal the map to CBOR
	return cbor.Marshal(data)
}

// Unmarshal decodes a CBOR-encoded Peer and assigns the result to the object, converting AddrInfo from JSON.
// NB! We need to pass a pointer to the Peer object to assign the result to it.
func UnmarshalFromCBOR(data []byte, p *Peer) error {
	var intermediateMap map[string]interface{}
	if err := cbor.Unmarshal(data, &intermediateMap); err != nil {
		return fmt.Errorf("failed to unmarshal CBOR to map: %w", err)
	}

	// Extract the AddrInfo JSON string
	addrInfoJSON, ok := intermediateMap["AddrInfo"].(string)
	if !ok {
		return fmt.Errorf("AddrInfo is not a valid JSON string")
	}

	// Unmarshal the JSON string to an AddrInfo object
	var addrInfo p2peer.AddrInfo
	if err := addrInfo.UnmarshalJSON([]byte(addrInfoJSON)); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to AddrInfo: %w", err)
	}

	// Assign the values to the Peer struct
	p.ID = intermediateMap["ID"].(string)
	p.Nick = intermediateMap["Nick"].(string)
	p.AddrInfo = &addrInfo

	return nil
}
