# blog-api

An implementation of an API that can power a personal blog.

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
