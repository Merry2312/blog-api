package server

import "go.mongodb.org/mongo-driver/bson/primitive"

// MongoDB structure for an article
type ArticleModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Author    string             `bson:"author"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt"`
}
