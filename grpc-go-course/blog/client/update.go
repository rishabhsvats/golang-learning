package main

import (
	"context"
	"log"

	pb "github.com/rishabhsvats/golang-learning/grpc-go-course/blog/proto"
)

func updateBlog(c pb.BlogServiceClient, id string) {
	log.Println("updateBlog was invoked")

	newBlog := &pb.Blog{
		Id:       id,
		AuthorId: "Rishabh",
		Title:    "A new title",
		Content:  "Updating the content of the first blog",
	}
	_, err := c.UpdateBlog(context.Background(), newBlog)

	if err != nil {
		log.Fatalf("Error happened while updating: %v\n", err)
	}

	log.Println("Blog was updated")
}
