package main

import (
	"github.com/Alexx1088/userservice/internal/db"
	"github.com/Alexx1088/userservice/internal/repository"
	"github.com/Alexx1088/userservice/internal/service"
	pb "github.com/Alexx1088/userservice/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	if err := db.Connect(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Pool.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set up repository and service
	userRepo := &repository.PgUserRepository{}
	userService := &service.UserService{Repo: userRepo}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)
	reflection.Register(grpcServer)

	log.Println("UserService gRPC server listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
