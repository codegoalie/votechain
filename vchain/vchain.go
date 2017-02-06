package vchain

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
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
	Number int
	Parent string
	Nonce  int

	Capacity int
	votes    []Vote
}

// Chain represents a block chain. It exposes functions for ease of use.
type Chain struct {
	currentBlock  Block
	blockCapacity int
	blocks        map[string]Block
}

// NewChain returns a new Chain with the starter block initilized to the Capacity
func NewChain(cap int) Chain {
	return Chain{
		currentBlock:  Block{Capacity: cap, Parent: "0000000000000000000000000000000000000000"},
		blockCapacity: cap,
		blocks:        map[string]Block{},
	}
}

// AddVote adds a vote to the current block on the chain. If he block is full,
// the block will be added to the chain and a new currentBlock will be created.
func (c *Chain) AddVote(v Vote) error {
	err := c.currentBlock.AddVote(v)
	if err != nil {
		return err
	}

	if len(c.currentBlock.votes) >= c.blockCapacity {
		blockHash := c.currentBlock.Hash()
		c.blocks[blockHash] = c.currentBlock
		c.currentBlock = Block{Number: c.currentBlock.Number + 1, Capacity: c.blockCapacity}
		c.currentBlock.Parent = blockHash
	}

	return nil
}

// Validate checks the VoterProof for authenticity (currently uninmplemented)
func (v Vote) Validate() (bool, error) {
	return true, nil
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
