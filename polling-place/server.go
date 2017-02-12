package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/codegoalie/votechain/vchain"
	pb "github.com/codegoalie/votechain/votechain"
)

type pollingPlaceServer struct {
	chain    vchain.Chain
	newVotes chan *pb.Vote
}

func (s *pollingPlaceServer) Cast(ctx context.Context, msg *pb.Vote) (*pb.Result, error) {
	log.Printf("Received cast: \n%+v", msg)
	vote := vchain.Vote{
		VoterProof: msg.VoterProof,
		RaceID:     int(msg.RaceId),
		Selection:  msg.Selection,
	}

	err := s.chain.AddVote(vote)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not cast vote: %s", err.Error())
		log.Println(errorMessage)
		return &pb.Result{
			Success: false,
			Message: errorMessage,
		}, err
	}

	log.Printf("Vote Cast! \n%+v\nSending to peer(s)...\n", s.chain)
	s.newVotes <- msg

	return &pb.Result{
		Success: true,
		Message: "Vote cast! Thanks for being a part!",
	}, nil
}

func (s pollingPlaceServer) GetLatestBlock(ctx context.Context, _ *pb.Empty) (*pb.Block, error) {
	return cast(s.chain.CurrentBlock), nil
}

func (s pollingPlaceServer) GetBlock(ctx context.Context, in *pb.BlockNumber) (*pb.Block, error) {
	if block, ok := s.chain.Blocks[in.Hash]; ok {
		return cast(block), nil
	}

	return nil, errors.New("Unknown block")
}

func (s pollingPlaceServer) Coordinate(client pb.PollingStation_CoordinateServer) error {
	for {
		fmt.Println("Receiving")
		in, err := client.Recv()
		if err == io.EOF {
			return nil

		}
		if err != nil {
			return err

		}

		switch msg := in.Transaction.(type) {
		case *pb.Coordination_Block:
			block := msg.Block
			if _, ok := s.chain.Blocks[block.Hash]; ok {
				log.Printf("Coordinated block found: %s\n", block.Hash)
			} else {
				log.Printf("Coordinated block not found: %s\nNeed to figure out how to add this\n", block.Hash)
			}
		case *pb.Coordination_Vote:
			vote := msg.Vote
			log.Printf("Coordinated vote received: %+v\n", vote)
			s.chain.AddVote(vchain.Vote{
				VoterProof: vote.VoterProof,
				RaceID:     int(vote.RaceId),
				Selection:  vote.Selection,
			})
			log.Printf("Coordinated vote Added! \n%+v\n", s.chain)
		}

		log.Println("At newVotes select")
		select {
		case vote := <-s.newVotes:
			log.Printf("Sending newVote: %+v\n", vote)
			client.Send(&pb.Coordination{
				Transaction: &pb.Coordination_Vote{
					Vote: vote,
				},
			})
		default:
			log.Println("Skip newVotes")
		}
	}
}

func newServer() *pollingPlaceServer {
	s := new(pollingPlaceServer)
	s.chain = vchain.NewChain(3)

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

func cast(orig vchain.Block) *pb.Block {
	votes := make([]*pb.Vote, len(orig.Votes))
	for _, vote := range orig.Votes {
		votes = append(votes, &pb.Vote{
			VoterProof: vote.VoterProof,
			RaceId:     int32(vote.RaceID),
			Selection:  vote.Selection,
		})
	}

	return &pb.Block{
		Number:   int32(orig.Number),
		Parent:   orig.Parent,
		Nonce:    int32(orig.Nonce),
		Capacity: int32(orig.Capacity),
		Votes:    votes,
	}
}
