package main

import (
	"context"
	"log"
	"time"

	pb "github.com/rishabhsvats/grpc/proto"
)

func callSayHelloClientStream(client pb.GreetServiceClient, name *pb.NameList) {
	log.Printf("client streaming started")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("could not send names: %v", err)
	}

	for _, name := range name.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("error while sending: %v", err)
		}
		log.Printf("sent the request with name: %s", name)
		time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	log.Printf("client streaming finished")
	if err != nil {
		log.Fatalf("error while receving: %v", err)
	}
	log.Printf("%v", res.Messages)
}
