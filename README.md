# Votechain

Votechain is a simple blockchain implementation specifically designed to record
votes in an election. The chain does not hold other data about the election;
only the votes cast. Another system, such as
[LockTheVote](https://github.com/codegoalie/lockthevote) should be used to track
elections, races, candidates, and verify voter identities.

"Transactions" on the chain are votes with fields to identify the race and
selections. A vote is not a record of a full ballot, just a selection for a
single race within an election.

Voter proofs are any value which uniquely differentiate a single voter from
other voters. If an election does not require verifiable votes, voter proofs
can even be blank or a single value. However, unique voter proof values can
be used to verify inclusion in the election. The simplest voter proof value
is a plain text identifer; such as an email or ID. This can easily be verified
visually or programatically. On the opposite end of the spectrum, voter proofs
can be non-interactive zero knowledge proofs. This kind of voter proof value
allows voters (who hold their private key) to verify their vote was correctly
recorded in the election while maintaining their anonymity. 
