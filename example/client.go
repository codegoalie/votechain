package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"

	pb "github.com/codegoalie/votechain/votechain"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPollingStationClient(conn)
	ticker := time.NewTicker(time.Second * 3)

	for _ = range ticker.C {
		vote := pb.Vote{
			VoterProof: randString(5),
			RaceId:     int32(rand.Int()),
			Selection:  randString(10),
		}
		log.Printf("Sending vote: %v\n", vote)
		_, err = client.Cast(context.Background(), &vote)

		if err != nil {
			log.Fatalf("Failed to send a vote: %v", err)
		}
	}
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
