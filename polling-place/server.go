package main

import (
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/codegoalie/votechain/vchain"
	pb "github.com/codegoalie/votechain/votechain"
)

type pollingPlaceServer struct {
	chain vchain.Chain
}

func (p *pollingPlaceServer) Cast(ctx context.Context, vote *pb.Vote) (*pb.Result, error) {
	log.Printf("Received cast: \n%+v", vote)
	err := p.chain.AddVote(vchain.Vote{
		VoterProof: vote.VoterProof,
		RaceID:     int(vote.RaceId),
		Selection:  vote.Selection,
	})

	if err != nil {
		log.Printf("Could not cast vote: \n%+v", p.chain)
		return &pb.Result{
			Success: false,
			Message: "Could not cast vote",
		}, err
	}

	log.Printf("Vote Cast! \n%+v", p.chain)
	return &pb.Result{
		Success: true,
		Message: "Vote cast! Thanks for being a part!",
	}, nil
}

func newServer() *pollingPlaceServer {
	s := new(pollingPlaceServer)
	s.chain = vchain.NewChain(3)
	s.chain.AddVote(vchain.Vote{
		VoterProof: "Dummy",
		RaceID:     1,
		Selection:  "2,3,1",
	})

	return s
}

func main() {
	fmt.Println("Listening on :4000")
	lis, err := net.Listen("tcp", fmt.Sprintf(":4000"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPollingStationServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
