package service

import (
	"context"
	"example.com/domain"
	"example.com/domain/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(UserRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository}
}

func (u *UserService) Add(ctx context.Context, id, name string) (string, error) {
	var user *domain.User
	user, _ = u.UserRepository.AddUser(ctx, id, name)
	return user.Name, nil
}

func (u *UserService) Delete(ctx context.Context, id string) (string, error) {
	_ = u.UserRepository.DeleteUser(ctx, id)
	return "", nil
}

func (u *UserService) Get(ctx context.Context, id string) (*domain.User, error) {
	var user *domain.User
	user, _ = u.UserRepository.GetUserByUserId(ctx, id)
	return user, nil
}
