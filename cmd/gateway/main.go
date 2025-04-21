package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"

	pb "userservice/pb/user"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()} // for dev only
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
