package repository

import (
	"context"
	"github.com/Alexx1088/userservice/internal/db"

	"github.com/Alexx1088/userservice/internal/model"
	"github.com/jackc/pgx/v5"
)

type PgUserRepository struct{}

func (r *PgUserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := db.Pool.Exec(ctx,
		"INSERT INTO users (id, name, email, password, score) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Name, user.Email, user.Score)
	return err
}

func (r *PgUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := db.Pool.QueryRow(ctx,
		"SELECT id, name, email, password, score FROM users WHERE id = $1", id)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Score)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *PgUserRepository) UpdateScore(ctx context.Context, id string, delta int32) (*model.User, error) {
	_, err := db.Pool.Exec(ctx,
		"UPDATE users SET score = score + $1 WHERE id = $2", delta, id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}
