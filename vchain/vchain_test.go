package vchain

import "testing"

func TestBlockHash(t *testing.T) {
	block := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		votes: []Vote{
			{
				VoterProof: "something",
				RaceID:     1,
				Selection:  "Candidate 1",
			},
		},
	}
	expected := "g_NcS4zEKtPoCy_ZIl3rkQfwFP0XcPmYYXJsdVYsVQU="

	actual := block.Hash()

	if actual != expected {
		t.Errorf("block.Hash() = %s, expected %s", actual, expected)
	}
}

func TestBlockHashIsUnique(t *testing.T) {
	block := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		votes: []Vote{
			{
				VoterProof: "something",
				RaceID:     1,
				Selection:  "Candidate 1",
			},
		},
	}
	block2 := Block{
		Parent: "00000000000000000000000000000000000000000000000000",
		votes: []Vote{
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
