package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/codegoalie/votechain/votechain"
)

func main() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPollingStationClient(conn)

	if err != nil {
		log.Fatalf("Failed to get stream: %v", err)
	}

	votes := []pb.Vote{
		{
			VoterProof: "123",
			RaceId:     2,
			Selection:  "1,2",
		},
		{
			VoterProof: "321",
			RaceId:     2,
			Selection:  "2,1",
		},
	}

	for _, vote := range votes {
		log.Printf("Sending vote: %v\n", vote)
		_, err := client.Cast(context.Background(), &vote)

		if err != nil {
			log.Fatalf("Failed to send a vote: %v", err)
		}
	}
}
