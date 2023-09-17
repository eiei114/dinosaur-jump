package infrastructure

import (
	"context"
	"example.com/domain"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	Conn *bun.DB
}

func NewUserRepository(Conn *bun.DB) *UserRepository {
	return &UserRepository{Conn: Conn}
}

func (u *UserRepository) AddUser(ctx context.Context, id, name string) (*domain.User, error) {
	return nil, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (u *UserRepository) GetUserByUserId(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}
