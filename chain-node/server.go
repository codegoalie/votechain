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

type chainNodeServer struct {
	chain vchain.Chain
}

func (s chainNodeServer) GetLatestBlock(ctx context.Context, _ *pb.Empty) (*pb.Block, error) {
	currentBlock := s.chain.CurrentBlock
	return &pb.Block{
		Number:   int32(currentBlock.Number),
		Parent:   currentBlock.Parent,
		Nonce:    currentBlock.Nonce,
		Capacity: int32(currentBlock.Capacity),
		Votes:    currentBlock.votes,
	}, nil
}

func (s chainNodeServer) GetBlock(ctx context.Context, in *pb.BlockNumber) (*pb.Block, error) {
	return &pb.Block{}, nil
}

func (s chainNodeServer) Coordinate(client pb.ChainNode_CoordinateServer) error {
	return nil
}

func newServer() *chainNodeServer {
	s := new(chainNodeServer)
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":4000"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChainNodeServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
