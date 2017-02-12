package vchain

import "testing"

func TestChainFindVoteinCurrent(t *testing.T) {
	chain := NewChain(2)
	vote := Vote{
		VoterProof: "124",
		RaceID:     1,
		Selection:  "2",
	}

	err := chain.AddVote(vote)
	if err != nil {
		t.Error("Add new vote to chain failed")
	}
	err = chain.AddVote(vote)
	if err == nil {
		t.Error("Add duplicate vote should fail, but didn't")
	}
}

func TestChainFindVoteinParent(t *testing.T) {
	chain := NewChain(1)
	vote := Vote{
		VoterProof: "124",
		RaceID:     1,
		Selection:  "2",
	}

	err := chain.AddVote(vote)
	if err != nil {
		t.Error("Add new vote to chain failed")
	}
	chain.AddVote(Vote{
		VoterProof: "432",
		RaceID:     1,
		Selection:  "1",
	})
	chain.AddVote(Vote{
		VoterProof: "6043",
		RaceID:     1,
		Selection:  "1",
	})

	err = chain.AddVote(vote)
	if err == nil {
		t.Error("Add duplicate vote should fail, but didn't")
	}
}

func TestChainFindVoteNoDup(t *testing.T) {
	chain := NewChain(1)
	vote := Vote{
		VoterProof: "124",
		RaceID:     1,
		Selection:  "2",
	}

	chain.AddVote(Vote{
		VoterProof: "432",
		RaceID:     1,
		Selection:  "1",
	})
	chain.AddVote(Vote{
		VoterProof: "6043",
		RaceID:     1,
		Selection:  "1",
	})

	err := chain.AddVote(vote)
	if err != nil {
		t.Errorf("Add new vote should not fail, but did: %s", err)
	}
}
