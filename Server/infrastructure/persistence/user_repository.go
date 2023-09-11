package persistence

import (
	"context"
	"example.com/domain"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) AddUser(ctx context.Context, id, authToken, name string) error {
	return nil
}

func (u *UserRepository) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}
