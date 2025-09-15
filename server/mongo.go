package server

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(uri string, dbName string, collectionName string) *mongo.Collection {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatal("Mongo connection error:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Mongo ping error:", err)
    }

    return client.Database(dbName).Collection(collectionName)
}
