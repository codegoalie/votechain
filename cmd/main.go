package main

import (
	"fmt"

	"github.com/codegoalie/votechain/vchain"
)

func main() {
	chain := vchain.NewChain((2))

	chain.AddVote(vchain.Vote{
		VoterProof: "Someh:SomeH(c||m)",
		RaceID:     1,
		Selection:  "1,2,3",
	})

	chain.AddVote(vchain.Vote{
		VoterProof: "Someotherh:SomeotherH(c||m)",
		RaceID:     1,
		Selection:  "2,3,1",
	})

	fmt.Printf("%+v", chain)
}
