package vchain

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Vote represents a single vote cast into the system. It is analagous to a
// transaction in a ledger-based chain.
type Vote struct {
	VoterProof string
	RaceID     int
	Selection  string
}

// Block represents a set of votes to be added to the chain
type Block struct {
	Parent string
	Votes  []Vote
}

// Hash is the SHA256 checksum of the block
func (b Block) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%+v", b)))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
