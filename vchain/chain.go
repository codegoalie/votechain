package vchain

// Chain represents a block chain. It exposes functions for ease of use.
type Chain struct {
	CurrentBlock  Block
	blockCapacity int
	Blocks        map[string]Block
}

// NewChain returns a new Chain with the starter block initilized to the Capacity
func NewChain(cap int) Chain {
	return Chain{
		CurrentBlock:  Block{Capacity: cap, Parent: "0000000000000000000000000000000000000000"},
		blockCapacity: cap,
		Blocks:        map[string]Block{},
	}
}

// AddVote adds a vote to the current block on the chain. If he block is full,
// the block will be added to the chain and a new currentBlock will be created.
func (c *Chain) AddVote(v Vote) error {
	err := c.CurrentBlock.AddVote(v)
	if err != nil {
		return err
	}

	if len(c.CurrentBlock.votes) >= c.blockCapacity {
		blockHash := c.CurrentBlock.Hash()
		c.Blocks[blockHash] = c.CurrentBlock
		c.CurrentBlock = Block{Number: c.CurrentBlock.Number + 1, Capacity: c.blockCapacity}
		c.CurrentBlock.Parent = blockHash
	}

	return nil
}
