package main

import (
	"errors"
	"flag"
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

func (s *pollingPlaceServer) Cast(ctx context.Context, vote *pb.Vote) (*pb.Result, error) {
	log.Printf("Received cast: \n%+v", vote)
	err := s.chain.AddVote(vchain.Vote{
		VoterProof: vote.VoterProof,
		RaceID:     int(vote.RaceId),
		Selection:  vote.Selection,
	})

	if err != nil {
		log.Printf("Could not cast vote: \n%+v", s.chain)
		return &pb.Result{
			Success: false,
			Message: "Could not cast vote",
		}, err
	}

	log.Printf("Vote Cast! \n%+v", s.chain)
	return &pb.Result{
		Success: true,
		Message: "Vote cast! Thanks for being a part!",
	}, nil
}

func (s pollingPlaceServer) GetLatestBlock(ctx context.Context, _ *pb.Empty) (*pb.Block, error) {
	currentBlock := s.chain.CurrentBlock

	votes := make([]*pb.Vote, len(currentBlock.Votes))
	for _, vote := range currentBlock.Votes {
		votes = append(votes, &pb.Vote{
			VoterProof: vote.VoterProof,
			RaceId:     int32(vote.RaceID),
			Selection:  vote.Selection,
		})
	}

	return &pb.Block{
		Number:   int32(currentBlock.Number),
		Parent:   currentBlock.Parent,
		Nonce:    int32(currentBlock.Nonce),
		Capacity: int32(currentBlock.Capacity),
		Votes:    votes,
	}, nil
}

func (s pollingPlaceServer) GetBlock(ctx context.Context, in *pb.BlockNumber) (*pb.Block, error) {
	if block, ok := s.chain.Blocks[in.Hash]; ok {
		votes := make([]*pb.Vote, len(block.Votes))
		for _, vote := range block.Votes {
			votes = append(votes, &pb.Vote{
				VoterProof: vote.VoterProof,
				RaceId:     int32(vote.RaceID),
				Selection:  vote.Selection,
			})
		}

		return &pb.Block{
			Number:   int32(block.Number),
			Parent:   block.Parent,
			Nonce:    int32(block.Nonce),
			Capacity: int32(block.Capacity),
			Votes:    votes,
		}, nil
	}

	return nil, errors.New("Unknown block")
}

func (s pollingPlaceServer) Coordinate(client pb.PollingStation_CoordinateServer) error {
	return nil
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
	port := flag.String("port", "4000", "Port to listen on")
	flag.Parse()

	fmt.Printf("Listening on :%s\n", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPollingStationServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
