package repository

import (
	"context"
	"example.com/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, id, authToken, name string) error
	GetUserByUserId(ctx context.Context, id string) (*domain.User, error)
}
