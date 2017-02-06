package vchain

import "testing"

func TestBlockHash(t *testing.T) {
	block := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		Votes: []Vote{
			{
				VoterProof: "something",
				RaceID:     1,
				Selection:  "Candidate 1",
			},
		},
	}
	expected := "BmW2GzA-ehN6fF7hEGkvqQSyJ1DDuwlnfqyCLIN5ts8="

	actual := block.Hash()

	if actual != expected {
		t.Errorf("block.Hash() = %s, expected %s", actual, expected)
	}
}

func TestBlockHashIsUnique(t *testing.T) {
	block := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		Votes: []Vote{
			{
				VoterProof: "something",
				RaceID:     1,
				Selection:  "Candidate 1",
			},
		},
	}
	block2 := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		Votes: []Vote{
			{
				VoterProof: "else",
				RaceID:     1,
				Selection:  "Candidate 1",
			},
		},
	}

	if block.Hash() == block2.Hash() {
		t.Error("vote.Hash() == vote2.Hash(), expected them to be different")
	}
}
