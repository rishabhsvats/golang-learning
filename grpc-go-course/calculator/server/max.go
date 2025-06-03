package main

import (
	"io"
	"log"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/calculator/proto"
)

func (s *Server) Max(stream pb.CalculatorService_MaxServer) error {
	log.Println("Max function was invoked")
	var maximum int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v\n", err)
		}

		if number := req.Number; number > maximum {
			maximum = number

			err := stream.Send(&pb.MaxResponse{
				Result: maximum,
			})
			if err != nil {
				log.Fatalf("error while sending data to client: %v\n", err)
			}

		}
	}
}
