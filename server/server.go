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

func (s *BlogServer) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	if req.Id == "" {
		return nil, errors.New("article ID is required")
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	var article ArticleModel
	err = s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		return nil, err
	}

	return &pb.GetArticleResponse{
		Article: &pb.Article{
			Id:        article.ID.Hex(),
			Title:     article.Title,
			Content:   article.Content,
			Author:    article.Author,
			CreatedAt: timestamppb.New(article.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(article.UpdatedAt.Time()),
		},
	}, nil
}

func (s *BlogServer) GetAllArticles(ctx context.Context, req *pb.GetAllArticlesRequest) (*pb.GetAllArticlesResponse, error) {
	var articles []*pb.Article

	var andFilters []bson.M

	if len(req.AuthorsFilter) > 0 {
		andFilters = append(andFilters, bson.M{"author": bson.M{"$in": req.AuthorsFilter}})
	}

	if req.TimeFilter != nil {
		andFilters = append(andFilters, bson.M{
			"createdAt": bson.M{
				"$gte": req.TimeFilter.Start.AsTime(),
				"$lte": req.TimeFilter.End.AsTime(),
			},
		})
	}

	if len(req.TextSearch) > 0 {
		searchString := ""
		for i, search := range req.TextSearch {
			if i > 0 {
				searchString += " "
			}
			searchString += search
		}
		andFilters = append(andFilters, bson.M{"$text": bson.M{"$search": searchString}})
	}

	var filters bson.M
	if len(andFilters) > 0 {
		filters = bson.M{"$and": andFilters}
	} else {
		filters = bson.M{}
	}

	log.Printf("GetAllArticles with filters: %s", filters)
	cursor, err := s.collection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var article ArticleModel
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, &pb.Article{
			Id:        article.ID.Hex(),
			Title:     article.Title,
			Content:   article.Content,
			Author:    article.Author,
			CreatedAt: timestamppb.New(article.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(article.UpdatedAt.Time()),
		})
	}

	return &pb.GetAllArticlesResponse{
		Articles: articles,
	}, nil
}

func (s *BlogServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {

	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}
	if req.Author == "" {
		return nil, errors.New("author is required")
	}

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

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert MongoDB ID")
	}

	var createdArticle ArticleModel
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&createdArticle)
	if err != nil {
		return nil, err
	}

	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:        createdArticle.ID.Hex(),
			Title:     createdArticle.Title,
			Content:   createdArticle.Content,
			CreatedAt: timestamppb.New(createdArticle.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(createdArticle.UpdatedAt.Time()),
		},
	}, nil
}

func (s *BlogServer) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*emptypb.Empty, error) {
	if req.Id == "" {
		return nil, errors.New("article ID is required")
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *BlogServer) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	if req.Id == "" {
		return nil, errors.New("article ID is required")
	}
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}
	if req.Author == "" {
		return nil, errors.New("author is required")
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"title":     req.Title,
		"content":   req.Content,
		"author":    req.Author,
		"updatedAt": time.Now(),
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	var updatedArticle ArticleModel
	err = s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedArticle)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateArticleResponse{
		Article: &pb.Article{
			Id:        updatedArticle.ID.Hex(),
			Title:     updatedArticle.Title,
			Content:   updatedArticle.Content,
			Author:    updatedArticle.Author,
			CreatedAt: timestamppb.New(updatedArticle.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(updatedArticle.UpdatedAt.Time()),
		},
	}, nil
}
