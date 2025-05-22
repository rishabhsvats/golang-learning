package main

import (
	"context"
	"log"
	"time"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/greet/proto"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Println("doLongGreet was invoked")

	reqs := []*pb.GreetRequest{
		{FirstName: "Rishabh"},
		{FirstName: "Ronaldo"},
		{FirstName: "Messi"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling the LongGreet %v \n", err)
	}

	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v\n ", err)
	}

	log.Printf("LongGreet: %s\n", res.Result)
}
