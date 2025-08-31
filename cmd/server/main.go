package main

import (
	"context"
	"github.com/Alexx1088/userservice/internal/db"
	"github.com/Alexx1088/userservice/internal/kafka"
	"github.com/Alexx1088/userservice/internal/repository"
	"github.com/Alexx1088/userservice/internal/service"
	pb "github.com/Alexx1088/userservice/pb/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 0) Load .env first (optional in Docker, helpful locally)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 1) Config
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "kafka:9092" // safe default for containers on shared network
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}

	log.Printf("Connecting to DB: %s", dsn)

	// 2) Init DB (assuming db.Init uses envs or the DSN you validated above)
	db.Init()
	defer db.Pool.Close()

	// 3) Init Kafka producer
	producer := kafka.NewProducer(brokers)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("producer.Close error: %v", err)
		}
	}()

	// 4) Build your deps: repo + service (inject the producer)
	userRepo := &repository.PgUserRepository{ /* add fields if needed */ }
	userSvc := service.NewUserService(userRepo, producer) // make a ctor; or &service.UserService{Repo: userRepo, Producer: producer}

	// 5) gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userSvc)
	reflection.Register(grpcServer)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("UserService gRPC listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Bye")
}
