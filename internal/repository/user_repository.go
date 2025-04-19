package repository

import (
	"context"
	"userservice/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	UpdateScore(ctx context.Context, id string, delta int32) (*model.User, error)
}
