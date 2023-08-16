package main

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

func extractPublicKeyFromPeerID(pid string) ([]byte, error) {

	// Decode the Peer ID from its string representation
	peerID, err := peer.Decode(pid)
	if err != nil {
		log.Fatalf("Failed to decode peer ID: %s", err)
	}

	// Extract the public key from the Peer ID
	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		log.Fatalf("Failed to extract public key from peer ID: %s", err)
	}

	// Print the public key
	pubKeyBytes, err := pubKey.Raw()
	if err != nil {
		log.Fatalf("Failed to get raw public key bytes: %s", err)
	}

	return pubKeyBytes, nil
}
