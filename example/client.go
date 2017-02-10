package main

import (
	"io"
	"log"

	context "golang.org/x/net/context"
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

	stream, err := client.Coordinate(context.Background())
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

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a vote/block : %v", err)
			}
			log.Printf("Got message %+v\n", in.Transaction)
		}
	}()
	for _, vote := range votes {
		log.Printf("Sending vote: %v\n", vote)
		err := stream.Send(&pb.Coordination{
			Transaction: &pb.Coordination_Vote{
				Vote: &vote,
			}})

		if err != nil {
			log.Fatalf("Failed to send a vote: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}
