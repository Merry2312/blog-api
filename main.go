package main

import (
	"log"
	"net"

	pb "blog-api/proto/blog/proto"
	"blog-api/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Printf("Connecting to MongoDB...")
	collection := server.ConnectMongo("mongodb://localhost:27017", "blogdb", "articles")
	log.Printf("Successfully connected to MongoDB")

	blogServer := server.NewBlogServer(collection)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server...")
	grpcServer := grpc.NewServer()

	pb.RegisterBlogServer(grpcServer, blogServer)

	reflection.Register(grpcServer)

	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
