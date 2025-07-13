package main

import (
	"context"
	"log"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/blog/proto"
)

func deleteBlog(c pb.BlogServiceClient, id string) {
	log.Println("deleteBlog was invoked")

	_, err := c.DeleteBlog(context.Background(), &pb.BlogId{Id: id})
	if err != nil {
		log.Printf("Error while deleting: %v\n", err)
	}

	log.Printf("Blog was deleted")

}
