package vchain

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// Block represents a set of votes to be added to the chain
type Block struct {
	Number int
	Parent string
	Nonce  int

	Capacity int
	votes    []Vote
}

// Hash is the SHA256 checksum of the block
func (b *Block) Hash() string {
	hasher := sha256.New()
	hash := ""
	for {
		hasher.Write([]byte(fmt.Sprintf("%+v", b)))
		hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		if hash[0:3] == "000" {
			break
		}
		b.Nonce++
	}
	return hash
}

// AddVote verifies and appends the vote to the block
func (b *Block) AddVote(v Vote) error {
	if len(b.votes) > b.Capacity {
		return errors.New("Block full")
	}

	if ok, err := v.Validate(); !ok {
		return err
	}

	b.votes = append(b.votes, v)
	return nil
}
