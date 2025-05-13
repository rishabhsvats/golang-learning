package main

import (
	"context"
	"log"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/calculator/proto"
)

func doSum(c pb.CalculatorServiceClient) {
	log.Println("doSum was invoked")
	res, err := c.Sum(context.Background(), &pb.SumRequest{
		FirstNumber:  57,
		SecondNumber: 54,
	})

	if err != nil {
		log.Fatalf("could not sum %v\n", err)
	}

	log.Printf("Sum: %d\n", res.Result)
}
