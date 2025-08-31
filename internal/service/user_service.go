package service

import (
	"context"
	"github.com/Alexx1088/userservice/internal/kafka"
	"github.com/Alexx1088/userservice/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//	userpb "github.com/Alexx1088/userservice/proto"

	"github.com/Alexx1088/userservice/internal/repository"
	pb "github.com/Alexx1088/userservice/pb/user"
	//	"github.com/google/uuid"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	Repo     repository.UserRepository
	Producer *kafka.Producer
}

func NewUserService(repo repository.UserRepository, producer *kafka.Producer) *UserService {
	return &UserService{
		Repo:     repo,
		Producer: producer,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	//	id := uuid.New().String()
	user := &model.User{
		ID:      req.GetUserId(),
		Name:    req.GetName(),
		Surname: req.GetSurname(),
		Email:   req.GetEmail(),
		Score:   0,
	}

	if err := s.Repo.Create(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.CreateUserResponse{
		UserId:  user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Score:   user.Score,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.Repo.GetByID(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserResponse{
		UserId:  user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Score:   user.Score,
	}, nil
}
