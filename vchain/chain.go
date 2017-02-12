package vchain

import "errors"

// Chain represents a block chain. It exposes functions for ease of use.
type Chain struct {
	CurrentBlock  Block
	blockCapacity int
	Blocks        map[string]Block
}

const seedParent = "0000000000000000000000000000000000000000"

// NewChain returns a new Chain with the starter block initilized to the Capacity
func NewChain(cap int) Chain {
	return Chain{
		CurrentBlock:  Block{Capacity: cap, Parent: seedParent},
		blockCapacity: cap,
		Blocks:        map[string]Block{},
	}
}

// AddVote adds a vote to the current block on the chain. If he block is full,
// the block will be added to the chain and a new currentBlock will be created.
func (c *Chain) AddVote(v Vote) error {
	if c.findVote(v) != "" {
		return errors.New("Vote already exists in chain")
	}

	err := c.CurrentBlock.AddVote(v)
	if err != nil {
		return err
	}

	if len(c.CurrentBlock.Votes) >= c.blockCapacity {
		blockHash := c.CurrentBlock.Hash()
		c.Blocks[blockHash] = c.CurrentBlock
		c.CurrentBlock = Block{Number: c.CurrentBlock.Number + 1, Capacity: c.blockCapacity}
		c.CurrentBlock.Parent = blockHash
	}

	return nil
}

// findVote will return the hash of the block containting the vote or
// and empty string if the vote is not found
func (c Chain) findVote(v Vote) string {
	blockToCheck := c.CurrentBlock

	for {
		for _, voteToCheck := range blockToCheck.Votes {
			if voteToCheck == v {
				return blockToCheck.Hash()
			}
		}
		if blockToCheck.Parent == seedParent {
			return ""
		}
		blockToCheck = c.Blocks[blockToCheck.Parent]
	}
}
