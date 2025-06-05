package main

import (
	"context"
	"log"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/calculator/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func doSqrt(c pb.CalculatorServiceClient, n int32) {
	log.Println("doSqrt was invoked")
	res, err := c.Sqrt(context.Background(), &pb.SqrtRequest{Number: n})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error message from server: %s\n", e.Message())
			log.Printf("Error message from server: %s\n", e.Code())

			if e.Code() == codes.InvalidArgument {
				log.Println("we probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("A non grpc error: %v\n", err)
		}
	}

	log.Printf("Sqrt: %f\n", res.Result)

}
