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
	// 1. Connect to MongoDB
	collection := server.ConnectMongo("mongodb://localhost:27017", "blogdb", "articles")

	// 2. Create your BlogServer instance
	blogServer := server.NewBlogServer(collection)

	// 3. Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// 4. Register BlogServer with gRPC
	pb.RegisterBlogServer(grpcServer, blogServer)

	// Enable reflection
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
