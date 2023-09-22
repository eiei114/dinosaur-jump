package repository

import (
	"context"
	"example.com/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, id, authToken, name string) error
	DeleteUser(ctx context.Context, id string) error
	GetUserByUserId(ctx context.Context, id string) (*domain.User, error)
	GetUserByAuthToken(ctx context.Context, authToken string) (*domain.User, error)
	GetUserRanking(ctx context.Context) ([]*domain.UserRanking, error)
}
