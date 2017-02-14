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
	chain       vchain.Chain
	broadcaster chan<- *pb.Vote
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
	s.broadcaster <- msg

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

func (s pollingPlaceServer) Mined(ctx context.Context, msg *pb.Block) (*pb.Empty, error) {
	// add block to chain
	return nil, nil
}

func newServer(capacity int, broadcaster chan<- *pb.Vote) *pollingPlaceServer {
	s := new(pollingPlaceServer)
	s.chain = vchain.NewChain(capacity)
	s.broadcaster = broadcaster

	return s
}

func main() {
	port := flag.String("port", "4000", "Port to listen on")
	peerPort := flag.String("peerport", "5000", "Port to coordinate with")
	capacity := flag.Int("capacity", 3, "Number of votes per block")
	flag.Parse()

	fmt.Printf("Listening on :%s\n", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	peers := []string{
		fmt.Sprintf("localhost:%s", *peerPort),
	}

	newVotes := make(chan *pb.Vote, 5)

	go broadcaster(newVotes, peers)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPollingStationServer(grpcServer, newServer(*capacity, newVotes))

	grpcServer.Serve(lis)
}

func broadcaster(newVotes <-chan *pb.Vote, peers []string) {
	fmt.Println("Starting broadcaster")
	clients := []pb.PollingStationClient{}
	for _, addr := range peers {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial: %s", err)
			continue
		}

		client := pb.NewPollingStationClient(conn)
		clients = append(clients, client)
	}

	ctx := context.Background()
	fmt.Printf("broadcaster ready with clients: %+v\n", clients)

	select {
	case vote := <-newVotes:
		for _, client := range clients {
			client.Cast(ctx, vote)
			fmt.Printf("broadcaster: Sent to %v\n", client)
		}
	}
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
