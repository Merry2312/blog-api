# blog-api

An implementation of an API that can power a personal blog.

## Operations

RPC | Description
--|--
getAllArticles | Retrieve all articles or a list of filtered articles. Articles can be filtered by creation time, author and text search on the article's title and content.
getArticle | Retrieve an article by it's ID
createArticle | Create an article
deleteArticle | Delete an article
updateArticle | Update an article

# Build protos

Generate the proto files using one of the following commands:

```
./build.sh
```

or

```
protoc --go_out=proto --go-grpc_out=proto proto/blog.proto
```

# Start the database

To start the MongoDB docker container:

```
docker compose up -d
```

# Start the golang gRPC server

To start the gRPC server, run:

```
go run main.go
```

# To test with gRPC UI

To start the gRPC UI instance:

```
grpcui -plaintext localhost:50051
```
