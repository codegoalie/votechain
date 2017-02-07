package vchain

// Vote represents a single vote cast into the system. It is analagous to a
// transaction in a ledger-based chain.
type Vote struct {
	VoterProof string
	RaceID     int
	Selection  string
}

// Validate checks the VoterProof for authenticity (currently uninmplemented)
func (v Vote) Validate() (bool, error) {
	return true, nil
}
