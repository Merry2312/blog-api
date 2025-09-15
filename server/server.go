package server

import (
	"context"
	"errors"
	"log"
	"time"

	pb "blog-api/proto/blog/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlogServer struct {
	collection *mongo.Collection
	pb.UnimplementedBlogServer
}

func NewBlogServer(collection *mongo.Collection) *BlogServer {
	return &BlogServer{collection: collection}
}

// TODO:
// Delete
// Update
// Get All

func (s *BlogServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	article := bson.M{
		"title":     req.Title,
		"content":   req.Content,
		"author":    req.Author,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}

	res, err := s.collection.InsertOne(ctx, article)
	if err != nil {
		log.Fatal("Error creating article:", err)
		return nil, err
	}

	// Convert MongoDB ObjectID to string
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert MongoDB ID")
	}

	// Fetch the inserted document to return exact object from DB
	var insertedDoc ArticleModel
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&insertedDoc)
	if err != nil {
		return nil, err
	}

	// Return gRPC Article object wrapped in CreateArticleResponse
	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:        insertedDoc.ID.Hex(),
			Title:     insertedDoc.Title,
			Content:   insertedDoc.Content,
			CreatedAt: timestamppb.New(insertedDoc.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(insertedDoc.UpdatedAt.Time()),
		},
	}, nil
}

func (s *BlogServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Fetch the inserted document to return exact object from DB
	var doc ArticleModel
	err = s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	// Return gRPC Article object wrapped in GetArticleResponse
	return &pb.GetArticleResponse{
		Article: &pb.Article{
			Id:        doc.ID.Hex(),
			Title:     doc.Title,
			Content:   doc.Content,
			Author:    doc.Author,
			CreatedAt: timestamppb.New(doc.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(doc.UpdatedAt.Time()),
		},
	}, nil
}

func (s *BlogServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*emptypb.Empty, error) {
	// TODO: Implement the logic to delete an article from MongoDB
	return &emptypb.Empty{}, nil
}
