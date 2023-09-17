package repository

import (
	"context"
	"example.com/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, id, name string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByUserId(ctx context.Context, id string) (*domain.User, error)
}
